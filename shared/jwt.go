package shared

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type UserClaims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (c UserClaims) Valid() error {
	return nil
}

func GenerateToken(claims *UserClaims) (*string, error) {
	key := []byte(jwtSecret())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(key)

	if err != nil {
		return nil, err
	}

	return &signedString, nil
}

func ParseToken(rawToken string) (*UserClaims, error) {
	token := strings.TrimPrefix(rawToken, "Bearer ")

	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, keyFunc)

	if err != nil {
		return nil, err
	}

	if err = jwtToken.Claims.Valid(); err != nil {
		return nil, err
	}

	return jwtToken.Claims.(*UserClaims), nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(jwtSecret()), nil
}

func jwtSecret() string {
	return viper.GetString("JWT_SECRET")
}
