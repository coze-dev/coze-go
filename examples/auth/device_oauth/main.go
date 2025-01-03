package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

// This examples is about how to use the device oauth process to acquire user authorization.
//
// Firstly, users need to access https://www.coze.com/open/oauth/apps. For the cn environment,
// users need to access https://www.coze.cn/open/oauth/apps to create an OAuth App of the type
// of TVs/Limited Input devices/Command line programs.
// The specific creation process can be referred to in the document:
// https://www.coze.com/docs/developer_guides/oauth_device_code. For the cn environment, it can be
// accessed at https://www.coze.cn/docs/developer_guides/oauth_device_code.
// After the creation is completed, the client ID can be obtained.

func main() {
	clientID := os.Getenv("COZE_DEVICE_OAUTH_CLIENT_ID")

	// The default access is api.coze.com, but if you need to access api.coze.cn,
	// please use base_url to configure the api endpoint to access
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.ComBaseURL
	}
	ctx := context.Background()

	oauth, err := coze.NewDeviceOAuthClient(clientID, coze.WithAuthBaseURL(cozeAPIBase))
	if err != nil {
		fmt.Printf("Failed to create OAuth client: %v\n", err)
		return
	}

	// In the device oauth authorization process, developers need to first call the interface
	// of Coze to generate the device code to obtain the user_code and device_code. Then generate
	// the authorization link through the user_code, guide the user to open the link, fill in the
	// user_code, and consent to the authorization. Developers need to call the interface of Coze
	// to generate the token through the device_code. When the user has not authorized or rejected
	// the authorization, the interface will throw an error and return a specific error code.
	// After the user consents to the authorization, the interface will succeed and return the
	// access_token.
	//
	// First, make a call to obtain 'getDeviceCode'

	codeResp, err := oauth.GetDeviceCode(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to get device code: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", codeResp)
	fmt.Println(codeResp.LogID())

	// The space permissions for which the Access Token is granted can be specified. As following codes:
	// codeResp, err = oauth.GetDeviceCode(ctx, &coze.GetDeviceOAuthCodeReq{
	// 	WorkspaceID: &workspaceID,
	// })

	// The returned device_code contains an authorization link. Developers need to guide users
	// to open up this link.
	// open codeResp.getVerificationUri

	fmt.Printf("Please open url: %s\n", codeResp.VerificationURL)

	// The developers then need to use the device_code to poll Coze's interface to obtain the token.
	// The SDK has encapsulated this part of the code in and handled the different returned error
	// codes. The developers only need to invoke getAccessToken.

	// if the developers set poll as true, the sdk will automatically handle pending and slow down exception
	resp, err := oauth.GetAccessToken(ctx, &coze.GetDeviceOAuthAccessTokenReq{
		DeviceCode: codeResp.DeviceCode,
		Poll:       true,
	})
	if err != nil {
		authErr, ok := coze.AsAuthError(err)
		if !ok {
			fmt.Printf("Failed to get access token: %v\n", err)
			return
		}
		switch authErr.Code {
		case coze.AccessDenied:
			// The user rejected the authorization.
			// Developers need to guide the user to open the authorization link again.
			fmt.Println("access denied")
		case coze.ExpiredToken:
			// The token has expired. Developers need to guide the user to open
			// the authorization link again.
			fmt.Println("expired token")
		default:
			fmt.Printf("Unexpected error: %v\n", err)

			return
		}
	}
	fmt.Printf("%+v\n", resp)

	// // use the access token to init Coze client
	// use the access token to init Coze client
	cozeCli := coze.NewCozeAPI(coze.NewTokenAuth(resp.AccessToken))
	_ = cozeCli

	// When the token expires, you can also refresh and re-obtain the token
	resp, err = oauth.RefreshToken(ctx, resp.RefreshToken)
	if err != nil {
		fmt.Printf("Failed to refresh token: %v\n", err)
		return
	}
	fmt.Println(resp)
	fmt.Println(resp.LogID())
}
