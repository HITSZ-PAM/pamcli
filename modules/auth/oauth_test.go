package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func startMockServer(c models.Client) *httptest.Server {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(c.ID), jwt.SigningMethodHS512))

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &c)
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	}))
}

func TestAuth_Token(t *testing.T) {
	type fields struct {
		Client       *http.Client
		ServerConfig models.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Good Server",
			fields{
				Client: &http.Client{},
				ServerConfig: models.Client{
					ID:     "000000",
					Secret: "999999",
					Domain: "http://localhost",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := startMockServer(tt.fields.ServerConfig)
			config := Config{
				ClientID:     tt.fields.ServerConfig.ID,
				ClientSecret: tt.fields.ServerConfig.Secret,
				TokenURL:     srv.URL,
			}
			ctx := context.Background()
			auth, _ := NewAuth(ctx, config)
			token, err := auth.Token()

			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Parse and verify jwt access token
			jwtToken, err2 := jwt.ParseWithClaims(token, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("parse error")
				}
				return []byte(tt.fields.ServerConfig.ID), nil
			})
			if (err2 != nil) != tt.wantErr {
				t.Errorf("Auth.Token() error = %v, wantErr %v", err2, tt.wantErr)
				return
			}
			_, ok := jwtToken.Claims.(*generates.JWTAccessClaims)
			if !ok || !jwtToken.Valid {
				t.Errorf("Auth.Token() error = %v, wantErr %v", ok, tt.wantErr)
				return
			}

		})
	}
}
