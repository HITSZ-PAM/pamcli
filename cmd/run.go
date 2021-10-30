/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/HITSZ-PAM/pamcli/modules/client"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a program within PAM context",
	Long:  `Resolve environment variables and pass them to the target program`,
	RunE: func(cmd *cobra.Command, args []string) error {

		newEnvList, err := getAccount()
		if err != nil {
			log.Fatalf("getting account failed: %s", err.Error())
		}

		subShell := exec.Command(args[0], args[1:]...)

		subShell.Env = newEnvList
		subShell.Stdin = os.Stdin
		subShell.Stdout = os.Stdout
		subShell.Stderr = os.Stderr

		done := make(chan bool, 1)
		// Pass all signals to child process
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		// Start subShell
		err = subShell.Start()
		if err != nil {
			return err
		}

		// Listen to signal and pass through to sub process
		go func() {
			select {
			case s := <-signals:
				err := subShell.Process.Signal(s)
				if err != nil && !strings.Contains(err.Error(), "process already finished") {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			case <-done:
				signal.Stop(signals)
				return
			}
		}()

		// Wait for finish
		commandErr := subShell.Wait()

		// close go routine
		done <- true

		if commandErr != nil {
			// Check if the program exited with an error
			exitErr, ok := commandErr.(*exec.ExitError)
			if ok {
				waitStatus, ok := exitErr.Sys().(syscall.WaitStatus)
				if ok {
					// Return the status code returned by the process
					os.Exit(waitStatus.ExitStatus())
					return nil
				}
			}
			return commandErr
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// getAccount return a env list
func getAccount() ([]string, error) {

	// New Client
	cfg := client.Config{
		ServerAddr:   os.Getenv("PAM_SERVER_URL"),
		ClientID:     os.Getenv("PAM_CLIENT_ID"),
		ClientSecret: os.Getenv("PAM_CLIENT_SECRET"),
	}

	ctx := context.Background()
	c, err := client.NewClient(ctx, &cfg)
	if err != nil {
		log.Printf("Failed to create a new client: %v", err)
		return nil, err
	}

	// Match ENV
	usernameRegexp := regexp.MustCompile(`pamcli://username/([\d]+)$`)
	passwordRegexp := regexp.MustCompile(`pamcli://password/([\d]+)$`)
	valueRegexp := regexp.MustCompile(`=(.*)$`)
	envList := os.Environ()
	for idx, env := range envList {
		params := usernameRegexp.FindStringSubmatch(env)
		if len(params) == 2 {
			accoundID := params[1]
			oldstring := valueRegexp.FindStringSubmatch(env)[1]
			username, _, err := c.Resolve(accoundID)
			if err != nil {
				return nil, fmt.Errorf("error resolving credentials: %s", err.Error())
			}
			newEnv := strings.Replace(env, oldstring, username, 1) // 1 means first occurance
			envList[idx] = newEnv
		}
		params = passwordRegexp.FindStringSubmatch(env)
		if len(params) == 2 {
			accoundID := params[1]
			oldstring := valueRegexp.FindStringSubmatch(env)[1]
			_, password, err := c.Resolve(accoundID)
			if err != nil {
				return nil, fmt.Errorf("error resolving credentials: %s", err.Error())
			}
			newEnv := strings.Replace(env, oldstring, password, 1) // 1 means first occurance
			envList[idx] = newEnv
		}
	}

	return envList, nil
}
