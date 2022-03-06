package apple

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
)

type AppleClaims struct {
	CHash          string `json:"c_hash"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	AuthTime       int64  `json:"auth_time"`
	NonceSupported bool   `json:"nonce_supported"`
	jwt.StandardClaims
}

type Key struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type Keys struct {
	Keys []*Key `json:"keys"`
}

var (
	pubKeys []*Key
)

func Base64URLToBase64(in string) string {
	out := strings.ReplaceAll(in, "-", "+")
	out = strings.ReplaceAll(out, "_", "/")
	if len(out)%4 != 0 {
		out = out + strings.Repeat("=", 4-len(out)%4)
	}
	return out
}

func GetPublicKeys() ([]*Key, error) {
	client := resty.New()
	resp, err := client.R().EnableTrace().Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}

	var keys Keys
	err = json.Unmarshal(resp.Body(), &keys)
	if err != nil {
		return nil, err
	}

	return keys.Keys, nil
}

func AuthFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid := token.Header["kid"].(string)

	if pubKeys == nil {
		keys, err := GetPublicKeys()
		if err != nil {
			return nil, err
		}
		pubKeys = keys
	}

	var key *Key
	for _, k := range pubKeys {
		if k.Kid == kid {
			key = k
			break
		}
	}

	result := getPublicKeyObject(key.N, key.E)

	return result, nil
}

func getPublicKeyObject(base64urlEncodedN string, base64urlEncodedE string) *rsa.PublicKey {

	var pub rsa.PublicKey
	var decE, decN []byte
	var eInt int
	var err error

	//get the modulo
	decN, err = base64.RawURLEncoding.DecodeString(base64urlEncodedN)
	if err != nil {
		return nil
	}
	pub.N = new(big.Int)
	pub.N.SetBytes(decN)
	//get exponent
	decE, err = base64.RawURLEncoding.DecodeString(base64urlEncodedE)
	if err != nil {
		return nil
	}
	//convert the bytes into int
	for _, v := range decE {
		eInt = eInt << 8
		eInt = eInt | int(v)
	}
	pub.E = eInt

	return &pub
}
