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
	"os"
	"os/exec"
	"os/signal"
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

		username, password, err := Get_account()
		newEnv := os.Environ()
		newEnv = append(newEnv, "PAM_ACCOUNT_USERNAME="+username)
		newEnv = append(newEnv, "PAM_ACCOUNT_PASSWORD="+password)

		Son_Shell := exec.Command(args[0], args[1:]...)

		Son_Shell.Env = newEnv
		Son_Shell.Stdin = os.Stdin
		Son_Shell.Stdout = os.Stdout
		Son_Shell.Stderr = os.Stderr

		err = Son_Shell.Start()
		if err != nil {
			return err
		}

		done := make(chan bool, 1)
		// Pass all signals to child process
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		go func() {
			select {
			case s := <-signals:
				err := Son_Shell.Process.Signal(s)
				if err != nil && !strings.Contains(err.Error(), "process already finished") {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			case <-done:
				signal.Stop(signals)
				return
			}
		}()
		commandErr := Son_Shell.Wait()
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Get_account() (string, string, error) {
	cfg := client.Config{
		ServerAddr:   os.Getenv("PAM_SERVER_URL"),
		ClientID:     os.Getenv("PAM_CLIENT_ID"),
		ClientSecret: os.Getenv("PAM_CLIENT_SECRET"),
	}
	ctx := context.Background()
	c, err := client.NewClient(ctx, &cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create a new client:", err)
		return "", "", err
	}
	username, password, err := c.Resolve(os.Getenv("PAM_ACCOUNT_ID"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to resolve account:", err)
		return "", "", err
	}
	return username, password, nil
}
