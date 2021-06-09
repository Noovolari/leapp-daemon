package gcp

import (
	"os"
	"testing"
)

func TestWorkingDefaultAuth(t *testing.T) {
	t.SkipNow()
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		t.Fatalf("GOOGLE_APPLICATION_CREDENTIALS should not be set")
	}

	keys, err := listKeys("leapp-test@forlunch-laboratory.iam.gserviceaccount.com")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	print(keys)
}

func TestOAuthGetAuthUrl(t *testing.T) {
	t.SkipNow()
	config := getConfig()
	url := oauthGetAuthorizationUrl(config)
	print(url)
}

func TestOAuthGetToken(t *testing.T) {
	t.SkipNow()
	config := getConfig()
	authCode := ""

	token := oauthGetTokenFromWeb(config, authCode)
	print("refresh token: " + token.RefreshToken + "\n")
	print("token: " + token.AccessToken + "\n")
}
