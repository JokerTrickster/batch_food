package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"github.com/lestrrat-go/jwx/jwk"
)

type JwtCustomClaims struct {
	CreateTime int64  `json:"createTime"`
	UserID     uint   `json:"userID"`
	Email      string `json:"email"`
	jwt.StandardClaims
}

var AccessTokenSecretKey []byte
var RefreshTokenSecretKey []byte
var JwtConfig middleware.JWTConfig

const (
	AccessTokenExpiredTime  = 24     // hour
	RefreshTokenExpiredTime = 24 * 7 // hour
)

func InitJwt() error {
	secret := "secret"
	AccessTokenSecretKey = []byte(secret)
	RefreshTokenSecretKey = []byte(secret)
	return nil
}

func GenerateToken(email string, userID uint) (string, int64, string, int64, error) {
	now := time.Now()
	accessToken, accessTknExpiredAt, err := GenerateAccessToken(email, now, userID)
	if err != nil {
		return "", 0, "", 0, err
	}
	refreshToken, refreshTknExpiredAt, err := GenerateRefreshToken(email, now, userID)
	if err != nil {
		return "", 0, "", 0, err
	}
	return accessToken, accessTknExpiredAt, refreshToken, refreshTknExpiredAt, nil
}

func GenerateAccessToken(email string, now time.Time, userID uint) (string, int64, error) {
	// Set custom claims
	expiredAt := now.Add(time.Hour * AccessTokenExpiredTime).Unix()
	claims := &JwtCustomClaims{
		TimeToEpochMillis(now),
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString(AccessTokenSecretKey)
	if err != nil {
		return "", 0, err
	}
	return accessToken, expiredAt, nil
}

func GenerateRefreshToken(email string, now time.Time, userID uint) (string, int64, error) {
	expiredAt := now.Add(time.Hour * RefreshTokenExpiredTime).Unix()
	claims := &JwtCustomClaims{
		TimeToEpochMillis(now),
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	refreshToken, err := token.SignedString(RefreshTokenSecretKey)
	if err != nil {
		return "", 0, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to generate refresh token - %v", err), ErrFromInternal)
	}
	return refreshToken, expiredAt, nil
}
func VerifyToken(tokenString string) error {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	if err != nil {
		return ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to parse token - %v", token), ErrFromClient)
	}

	// Check token validity
	if !token.Valid {
		return ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("invalid token - %v", token), ErrFromClient)
	}
	return nil
}
func ParseToken(tokenString string) (uint, string, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	// if err != nil {
	// 	return 0, "", ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to parse token - %v", token), ErrFromClient)
	// }
	// Extract claims
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return 0, "", ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to extract claims - %v", token), ErrFromClient)
	}

	// Extract email and userID
	email := claims.Email
	userID := claims.UserID
	return userID, email, nil
}

// JwtVerify verify data
func jwtVerifyWithKeySet(ctx context.Context, p AuthProvider, tokenString string, keySetUrl string) (jwt.MapClaims, error) {

	ctxHttp, ctxHttpCancel := context.WithTimeout(ctx, time.Second*10)
	defer ctxHttpCancel()

	req, err := http.NewRequestWithContext(ctxHttp, http.MethodGet, keySetUrl, nil)
	if err != nil {
		return jwt.MapClaims{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("request jwt pub key for %s - url : %s", authProviderName[p], keySetUrl), ErrFromClient)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return jwt.MapClaims{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("request jwt pub key for %s - url : %s", authProviderName[p], keySetUrl), ErrFromClient)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return jwt.MapClaims{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("request jwt pub key for %s resCode not 200 - url : %s / resCode : %d", authProviderName[p], keySetUrl, res.StatusCode), ErrFromClient)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return jwt.MapClaims{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("read response jwt pub key for %s - url : %s", authProviderName[p], keySetUrl), ErrFromClient)
	}

	set, err := jwk.Parse(bytes)
	if err != nil {
		return jwt.MapClaims{}, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("parse response jwt pub key for %s - %s", authProviderName[p], keySetUrl), ErrFromClient)
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve the key ID from the token header
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing key ID (kid) in token header")
		}

		// Look up the key in the key set
		key, exists := set.LookupKeyID(keyID)
		if !exists {
			return nil, fmt.Errorf("key ID %s not found", keyID)
		}

		var pubKey interface{}
		if err := key.Raw(&pubKey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %w", err)
		}
		return pubKey, nil
	})
	fmt.Println(token.Valid)
	return claims, nil
}

func aesEncrypt(ctx context.Context, byteToEncrypt []byte, keyString string) (string, error) {

	// since the key is in string, we need to convert decode it to bytes
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), "decode aes encrypt key string", ErrFromClient)
	}

	// create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), fmt.Sprintf("create aes encrypt cipher - %v", err), ErrFromInternal)
	}

	// create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), fmt.Sprintf("decode aes encrypt key gcm - %v", err), ErrFromInternal)
	}

	// create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), fmt.Sprintf("create aes encrypt nonce - %v", err), ErrFromInternal)
	}

	// encrypt the data using aesGCM.Seal
	return base64.StdEncoding.EncodeToString(aesGCM.Seal(nonce, nonce, byteToEncrypt, nil)), nil
}

func aesDecrypt(ctx context.Context, stringToDecrypt string, keyString string) (string, error) {

	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), fmt.Sprintf("decode aes decode key string"), ErrFromInternal)
	}
	enc, err := base64.StdEncoding.DecodeString(stringToDecrypt)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("decode encrypted string"), ErrFromClient)
	}

	// create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), fmt.Sprintf("create aes decrypt cipher - %v", err), ErrFromInternal)

	}

	// create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrInternalServer, Trace(), fmt.Sprintf("decode aes decrypt key gcm - %v", err), ErrFromInternal)

	}

	// get the nonce size
	nonceSize := aesGCM.NonceSize()

	// extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("decrypt token data - %v", err), ErrFromClient)
	}

	return string(plaintext), nil
}
