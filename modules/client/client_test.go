package client_test

import (
	"context"
	"os"
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
			got, got1, err := c.Resolve(os.Getenv("PAM_ACCOUNT_ID"))
			if err != nil {
				t.Errorf("client.Resolve() error: %v", err)
				return
			}
			if got != os.Getenv("PAM_ACCOUNT_USERNAME") {
				t.Errorf("client.Resolve() got = %v, want %v", got, os.Getenv("PAM_ACCOUNT_USERNAME"))
			}
			if got1 != os.Getenv("PAM_ACCOUNT_PASSWORD") {
				t.Errorf("client.Resolve() got1 = %v, want %v", got1, os.Getenv("PAM_ACCOUNT_PASSWORD"))
			}
		})
	}
}
