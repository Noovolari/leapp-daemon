package gcp

import (
  "golang.org/x/oauth2"
  "reflect"
  "testing"
  "time"
)

func TestGetConfig(t *testing.T) {
  config, err := getConfig()

  if err != nil {
    t.Fatalf("unexpected error %v", err)
  }

  expectedScope := []string{"openid",
    "https://www.googleapis.com/auth/userinfo.email",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/appengine.admin",
    "https://www.googleapis.com/auth/compute",
    "https://www.googleapis.com/auth/accounts.reauth"}

  expectedClientSecret := "ZmssLNjJy2998hD4CTg2ejr2"
  expectedClientId := "32555940559.apps.googleusercontent.com"

  if !reflect.DeepEqual(config.Scopes, expectedScope) {
    t.Fatalf("unexpected config scope %v", config.Scopes)
  }

  if config.ClientID != expectedClientId {
    t.Fatalf("unexpected client id %v", config.ClientID)
  }

  if config.ClientSecret != expectedClientSecret {
    t.Fatalf("unexpected client secret %v", config.ClientSecret)
  }
}

func TestGetConfig_singleton(t *testing.T) {
  config, _ := getConfig()
  configSingleton, _ := getConfig()

  if config != configSingleton {
    t.Fatalf("singleton is not returning the same instance")
  }
}

func TestGetCredentials(t *testing.T) {
  token := oauth2.Token{
    AccessToken:  "access token",
    TokenType:    "token type",
    RefreshToken: "refresh token",
    Expiry:       time.Time{},
  }

  api := GcpApi{}
  credentials := api.GetCredentials(&token)

  expectedCredentials := "{\n  \"client_id\":\"32555940559.apps.googleusercontent.com\",\n" +
    "  \"client_secret\":\"ZmssLNjJy2998hD4CTg2ejr2\",\n" +
    "  \"token_uri\":\"https://oauth2.googleapis.com/token\",\n" +
    "  \"refresh_token\": \"refresh token\",\n" +
    "  \"revoke_uri\": \"https://accounts.google.com/o/oauth2/revoke\",\n" +
    "  \"appScopes\": [\"openid\",\"https://www.googleapis.com/auth/userinfo.email\"," +
    "\"https://www.googleapis.com/auth/cloud-platform\",\"https://www.googleapis.com/auth/appengine.admin\"," +
    "\"https://www.googleapis.com/auth/compute\",\"https://www.googleapis.com/auth/accounts.reauth\"],\n" +
    "  \"type\": \"authorized_user\"\n}"

  if credentials != expectedCredentials {
    t.Fatalf("unexpected credentials: %v", credentials)
  }
}
