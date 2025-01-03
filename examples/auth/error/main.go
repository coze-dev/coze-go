package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

// Using web oauth as example, show how to handle the Auth exception
func main() {
	redirectURI := os.Getenv("COZE_WEB_OAUTH_REDIRECT_URI")
	clientSecret := os.Getenv("COZE_WEB_OAUTH_CLIENT_SECRET")
	clientID := os.Getenv("COZE_WEB_OAUTH_CLIENT_ID")
	ctx := context.Background()

	// The default access is api.coze.com, but if you need to access api.coze.cn,
	// please use base_url to configure the api endpoint to access
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.ComBaseURL
	}

	// The sdk offers the WebOAuthClient class to establish an authorization for Web OAuth.
	// Firstly, it is required to initialize the WebOAuthApp with the client ID and client secret.
	oauth, err := coze.NewWebOAuthClient(clientID, clientSecret, coze.WithAuthBaseURL(cozeAPIBase))
	if err != nil {
		fmt.Printf("Failed to create OAuth client: %v\n", err)
		return
	}

	// Generate the authorization link and direct the user to open it.
	oauthURL := oauth.GetOAuthURL(ctx, &coze.GetWebOAuthURLReq{
		RedirectURI: redirectURI,
		State:       "state",
	})
	fmt.Println(oauthURL)

	// The space permissions for which the Access Token is granted can be specified. As following codes:
	// oauthURL := oauth.GetOAuthURL(&coze.GetWebOAuthURLReq{
	// 	RedirectURI: redirectURI,
	// 	State:       "state",
	// 	WorkspaceID: &workspaceID,
	// })
	// fmt.Println(oauthURL)

	// After the user clicks the authorization consent button, the coze web page will redirect
	// to the redirect address configured in the authorization link and carry the authorization
	// code and state parameters in the address via the query string.
	//
	// Get from the query of the redirect interface: query.get('code')
	code := "mock code"

	// After obtaining the code after redirection, the interface to exchange the code for a
	// token can be invoked to generate the coze access_token of the authorized user.
	resp, err := oauth.GetAccessToken(ctx, &coze.GetWebOAuthAccessTokenReq{
		Code:        code,
		RedirectURI: redirectURI,
	})
	if err != nil {
		fmt.Printf("Failed to get access token: %v\n", err)
		// The SDK has enumerated existing authentication error codes
		// You need to handle the exception and guide users to re-authenticate
		// For different oauth type, the error code may be different,
		// you should read document to get more information
		authErr, ok := coze.AsAuthError(err)
		if ok {
			switch authErr.Code {
			}
		}
		return
	}
	fmt.Println(resp)
}
