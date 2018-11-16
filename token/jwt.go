package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func parseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("Invalid Pri RSA KEY")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Invalid Pri RSA KEY")
	}

	return priv, nil
}

func parseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("Invalid Pub RSA KEY")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Invalid Pub RSA KEY")
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("Invalid Pub RSA KEY")
	}
}

func (t *TokenManager) CreateAccessToken(id uint, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":   id,
		"name": name,
		"exp":  time.Now().Unix() + t.maxAccessTokenLife,
	})
	tokenString, err := token.SignedString(t.privateKey)
	return tokenString, err
}

func (t *TokenManager) CreateRefreshToken(id uint, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":   id,
		"name": name,
		"exp":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString(t.privateKey)
	if err == nil {
		strs := strings.Split(tokenString, ".")
		return strs[2], nil
	} else {
		return "", err
	}
}

func (t *TokenManager) ValidateAccessToken(tokenString string) (uint, string, error) {
	claim, err := t.ExtractUserInfo(tokenString)
	if err == nil {
		now := time.Now().Unix()
		if claim.VerifyExpiresAt(now, true) {
			return uint(claim["id"].(float64)), claim["name"].(string), nil
		}
	}
	return 0, "", errors.New("Access Token Invalid")
}

func (t *TokenManager) ExtractUserInfo(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return t.publicKey, nil
	})

	if err != nil {
		return nil, errors.New("Access Token Invalid")
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claim, nil
	}
	return nil, errors.New("Access Token Invalid")
}
