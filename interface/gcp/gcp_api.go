package gcp

import (
  "context"
  "fmt"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "leapp_daemon/infrastructure/http/http_error"
  "sync"
)

var (
  configMutex      sync.Mutex
  configSingleton  *oauth2.Config
  commonParameters = `"client_id":"32555940559.apps.googleusercontent.com",
  "client_secret":"ZmssLNjJy2998hD4CTg2ejr2",
  "token_uri":"https://oauth2.googleapis.com/token",`
  appScopes = []string{"openid",
    "https://www.googleapis.com/auth/userinfo.email",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/appengine.admin",
    "https://www.googleapis.com/auth/compute",
    "https://www.googleapis.com/auth/accounts.reauth"}
)

type GcpApi struct {
}

func (api *GcpApi) GetOauthUrl() (string, error) {
  config, err := getConfig()
  if err != nil {
    return "", err
  }
  return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
}

func (api *GcpApi) GetOauthToken(authCode string) (*oauth2.Token, error) {
  config, err := getConfig()
  if err != nil {
    return nil, err
  }

  ctx := context.Background()
  token, err := config.Exchange(ctx, authCode)
  if err != nil {
    return nil, http_error.NewBadRequestError(err)
  }
  return token, nil
}

func (api *GcpApi) GetCredentials(oauthToken *oauth2.Token) string {
  scopesList := ""
  for i, appScope := range appScopes {
    if i < len(appScopes)-1 {
      scopesList += fmt.Sprintf("\"%v\",", appScope)
    } else {
      scopesList += fmt.Sprintf("\"%v\"", appScope)
    }
  }

  return fmt.Sprintf(`{
  %v
  "refresh_token": "%v",
  "revoke_uri": "https://accounts.google.com/o/oauth2/revoke",
  "appScopes": [%v],
  "type": "authorized_user"
}`, commonParameters, oauthToken.RefreshToken, scopesList)
}

func getConfig() (*oauth2.Config, error) {
  configMutex.Lock()
  defer configMutex.Unlock()

  if configSingleton == nil {
    configData := fmt.Sprintf(`{
    "installed": {
        %v
        "auth_uri":"https://accounts.google.com/o/oauth2/auth",
        "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
        "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
    }
}`, commonParameters)
    config, err := google.ConfigFromJSON([]byte(configData), appScopes...)
    if err != nil {
      return nil, http_error.NewInternalServerError(err)
    }
    configSingleton = config
  }
  return configSingleton, nil
}
