package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	TOKEN_Key        = os.Getenv("JWT_SECRET")
	TOKEN_Expiration = 24 * time.Hour
)

func GenerateJwtToken(authId string) (string, error) {
	payload := Token{
		AuthId:         authId,
		ExpirationTime: time.Now().Add(TOKEN_Expiration),
	}
	claims := jwt.MapClaims{
		"payload": payload,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(TOKEN_Key))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ValidateJwtToken(tokenString string) (*Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(TOKEN_Key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payloadInterface := claims["payload"]

		payloadToken := Token{}

		payloadByte, err := json.Marshal(payloadInterface)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(payloadByte, &payloadToken)
		if err != nil {
			return nil, err
		}
		now := time.Now()
		if now.After(payloadToken.ExpirationTime) {
			return nil, errors.New("Token Expired")
		}
		return &payloadToken, nil
	} else {
		return nil, errors.New("Unauthorized")
	}
}
