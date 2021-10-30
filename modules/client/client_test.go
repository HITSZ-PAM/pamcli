package client_test

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/HITSZ-PAM/pamcli/modules/client"
)

func Test_client_Resolve(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Valid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := client.Config{
				ServerAddr:   os.Getenv("PAM_SERVER_URL"),
				ClientID:     os.Getenv("PAM_CLIENT_ID"),
				ClientSecret: os.Getenv("PAM_CLIENT_SECRET"),
			}
			ctx := context.Background()
			c, err := client.NewClient(ctx, &cfg)
			if err != nil {
				t.Errorf("NewClient() error: %v", err)
				return
			}
			_, _, err = c.Resolve(os.Getenv("PAM_ACCOUNT_ID"))
			if err != nil {
				t.Errorf("client.Resolve() error: %v", err)
				return
			}
		})
	}
}

type benchmarkMockClient struct {
}

func (b *benchmarkMockClient) Resolve(dummy string) (string, string, error) {
	return "", "", nil
}

func Benchmark_client_Resolve(b *testing.B) {
	tests := []struct {
		name string
	}{
		{name: "Valid"},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			c := benchmarkMockClient{}

			// Match ENV
			usernameRegexp := regexp.MustCompile(`pamcli://username/([\d]+)$`)
			passwordRegexp := regexp.MustCompile(`pamcli://password/([\d]+)$`)
			valueRegexp := regexp.MustCompile(`=(.*)$`)
			var envList []string
			envList = append(envList, "TEST=pamcli://username/1234567890")
			for i := 0; i < b.N; i++ {
				for _, env := range envList {
					params := usernameRegexp.FindStringSubmatch(env)
					if len(params) == 2 {
						accoundID := params[1]
						oldstring := valueRegexp.FindStringSubmatch(env)[1]
						username, _, _ := c.Resolve(accoundID)
						strings.Replace(env, oldstring, username, 1) // 1 means first occurance
					}
					params = passwordRegexp.FindStringSubmatch(env)
					if len(params) == 2 {
						accoundID := params[1]
						oldstring := valueRegexp.FindStringSubmatch(env)[1]
						_, password, _ := c.Resolve(accoundID)
						strings.Replace(env, oldstring, password, 1) // 1 means first occurance
					}
				}
			}
		})
	}
}
