package subscriptionManager

import (
	"context"
	"log"
	"strings"
	"subscriptionManager/util"

	oidc "github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	OAuth2Config   oauth2.Config
	OAuth2State    string
	OAuth2Verifier *oidc.IDTokenVerifier
)

func InitAuthenticator() (error *error) {
	ctx := context.Background()
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	provider, err := oidc.NewProvider(ctx, config.Security.OAuth2.IssuerUrl)
	if err != nil {
		return &err
	}

	//redirectURL := "http://localhost:8181/demo/callback"
	// Configure an OpenID Connect aware OAuth2 client.
	OAuth2Config = oauth2.Config{
		ClientID:    config.Security.OAuth2.ClientID,
		RedirectURL: config.Security.OAuth2.WebRedirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: config.Security.OAuth2.Scopes,
	}
	OAuth2State = "WSPEngineOAuth2State"

	oidcConfig := &oidc.Config{
		ClientID: config.Security.OAuth2.ClientID,
		// For some reason when validating the bearer token in IsApiAuthenticated, the token has audience [account]
		// and not the client id (test_app) which causes the verifier to spit out the error oidc:
		// expected audience test_app got [account]". For that reason we skip this check
		SkipClientIDCheck: true,
	}
	OAuth2Verifier = provider.Verifier(oidcConfig)
	return nil
}

// get result and userid
func GetAuthentication(c *gin.Context) (userId string, result int) {
	rawAccessToken := c.Request.Header.Get("Authorization")
	if rawAccessToken == "" {
		result = 1
		return
	}
	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		result = 2
		return
	}
	ctx := context.Background()
	idToken, err := OAuth2Verifier.Verify(ctx, parts[1])
	// idiotic go design - to mute "idToken declared but not used" "error"
	if err != nil {
		log.Println("Failed to verify ID Token: " + err.Error())
		result = 3
		return
	}
	// get user id from idToken
	userId = idToken.Subject
	result = 0
	return
}

func IsApiAuthenticated(c *gin.Context) int {
	_, result := GetAuthentication(c)
	return result
}

func GetUserId(c *gin.Context) string {
	userId, _ := GetAuthentication(c)
	return userId
}
