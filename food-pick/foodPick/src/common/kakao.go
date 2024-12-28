package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/common/aws"
	"net/http"
	"strconv"
	"time"
)

type KakaoData struct {
	ID              int64 `json:"id"`
	ExpiresIn       int64 `json:"expires_in"`
	ExpiresInMillis int64 `json:"expiresInMillis"`
	AppID           int64 `json:"app_id"`
	AppId           int64 `json:"appId"`
}

var KakaoAppID int64

func InitKakaoOauth() error {
	appID, err := aws.AwsSsmGetParam("dev_kakao_app_id")
	if err != nil {
		return err
	}
	KakaoAppID, err = strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return err
	}

	return nil
}

func KakaoValidate(ctx context.Context, token string) (OAuthData, error) {
	fmt.Println(token)
	// Define Kakao API URL for token validation
	url := "https://kapi.kakao.com/v1/user/access_token_info"

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to create request: %v", err), ErrFromKakao)
	}

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+token)

	// Create an HTTP client and set timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to send request: %v", err), ErrFromKakao)
	}
	defer resp.Body.Close()

	// Check if the response status is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to validate token: %v", resp.Status), ErrFromKakao)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to read response body: %v", err), ErrFromKakao)
	}
	// Parse the JSON response
	var data KakaoData
	if err := json.Unmarshal(body, &data); err != nil {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to parse response: %v", err), ErrFromKakao)
	}
	if data.AppID != KakaoAppID || data.AppId != KakaoAppID {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("invalid token + %v", data.AppID), ErrFromKakao)
	}
	return OAuthData{}, nil
}
