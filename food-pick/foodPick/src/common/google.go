package common

import (
	"common/aws"
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleConfig = oauth2.Config{}

type tAuthMeta struct {
	GoogleIosID string   `json:"googleIosID" validate:"required"`
	GoogleAndID []string `json:"googleAndID" validate:"required"`
	// KakaoID     int      `json:"kakaoId" validate:"required"`
}

type OAuthData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}
type AuthProvider uint8

const (
	AuthProviderGoogle = AuthProvider(0)
	AuthProviderKakao  = AuthProvider(1)
	AuthProviderNaver  = AuthProvider(2)
)

var authProviderName map[AuthProvider]string = map[AuthProvider]string{
	AuthProviderGoogle: "google",
	AuthProviderKakao:  "kakao",
	AuthProviderNaver:  "naver",
}
var authMeta tAuthMeta
var httpClient = http.DefaultClient
var OpenAIKey string

func InitGoogleOauth() error {
	clientID, err := getClientID()
	if err != nil {
		return err
	}
	clientSecret, err := getClientSecret()
	if err != nil {
		return err
	}
	GoogleConfig = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/v0.2/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	if !Env.IsLocal {
		GoogleConfig.RedirectURL = "https://food-recommendation.jokertrickster.com"
	}
	//706198886562-8bv6bg4dnahtub52keb0v6runkfojbpe.apps.googleusercontent.com
	authMeta.GoogleIosID = "706198886562-8bv6bg4dnahtub52keb0v6runkfojbpe.apps.googleusercontent.com"
	authMeta.GoogleAndID = []string{"706198886562-8bv6bg4dnahtub52keb0v6runkfojbpe.apps.googleusercontent.com"}
	//eyJhbGciOiJSUzI1NiIsImtpZCI6ImQ3YjkzOTc3MWE3ODAwYzQxM2Y5MDA1MTAxMmQ5NzU5ODE5MTZkNzEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI3MDYxOTg4ODY1NjItcXVzazB0dWJmZnVna3N2YWhucWY0dmpybjRzNTVkZ3EuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI3MDYxOTg4ODY1NjItOGJ2NmJnNGRuYWh0dWI1MmtlYjB2NnJ1bmtmb2picGUuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTQxNTU5MDMxNDQ3ODkxNDY5NDciLCJlbWFpbCI6ImEwMTAyODE4Mjc5NkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6IuydtOyasOyjvCIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NKeFJTT1VabTE3dHhPOXpwWTdVTEFBNVdOREEybFhpeDRvdkg5TG1SQXh0LWdRaWdLVD1zOTYtYyIsImdpdmVuX25hbWUiOiLsmrDso7wiLCJmYW1pbHlfbmFtZSI6IuydtCIsImlhdCI6MTcyNjEwODQyMiwiZXhwIjoxNzI2MTEyMDIyfQ.FwecqOnBoHNhBCdjqVDKv9ltXqepgMvjPQQAmOOSFCJDy-NrwAU-dZCYKi3xqW928HOEjIwXhSiXqkQU9CDwq80yABzVALTaw7-js8YaFz-mp7mNXCpdBj2PY5OWRC_O9inJJwq6TQn_PgiELBMRwG7-B79No7Pkosfz1mxJcUCMesq0yXQQypVQfoFEciSeRuz4cBRMbNgE-TDMsjtCom2gkTJqnlR7ka4RcGKItp0MLo6fOBug-50k2xezs2yH2NyBUZYlc

	OpenAIKey, err = aws.AwsSsmGetParam("openai_api_key")
	if err != nil {
		return err
	}
	return nil
}

func getClientID() (string, error) {
	if Env.IsLocal {
		clientID, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
		if !ok {
			return "", fmt.Errorf("GOOGLE_CLIENT_ID not found")
		}
		return clientID, nil
	} else {
		ClientID, err := aws.AwsSsmGetParam("food_google_client_id")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}

func getClientSecret() (string, error) {
	if Env.IsLocal {

		clientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
		if !ok {
			return "", fmt.Errorf("GOOGLE_CLIENT_ID not found")
		}
		return clientSecret, nil

	} else {
		ClientID, err := aws.AwsSsmGetParam("food_google_client_secret")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}

func GoogleValidate(ctx context.Context, token string) (OAuthData, error) {
	claims, iErr := jwtVerifyWithKeySet(ctx, AuthProviderGoogle, token, "https://www.googleapis.com/oauth2/v3/certs")
	if iErr != nil {
		return OAuthData{}, iErr
	}

	fmt.Println(claims)
	aud, okAud := claims["aud"].(string)
	azp, _ := claims["azp"].(string)
	iss, okIss := claims["iss"].(string)
	sub, okSub := claims["sub"].(string)
	email, okEmail := claims["email"].(string)
	if !okAud || !okIss || !okSub || !okEmail ||
		(aud != authMeta.GoogleIosID && isGoogleIdAndNotExisted(aud) && azp != authMeta.GoogleIosID && isGoogleIdAndNotExisted(azp)) ||
		(iss != "accounts.google.com" && iss != "https://accounts.google.com") {
		return OAuthData{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("not valid token claims from provider Google - %+v", claims), ErrFromClient)
	}
	oauthData := OAuthData{
		ID:       sub,
		Email:    email,
		Provider: "google",
	}
	return oauthData, nil
}

func isGoogleIdAndNotExisted(key string) bool {
	for _, andKey := range authMeta.GoogleAndID {
		if key == andKey {
			return false
		}
	}
	return true
}
