package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type TokenCreator struct {
	signature string
}

type TokenClaims struct {
	Contract []string `json:"contract"`
	jwt.StandardClaims
}

func NewTokenCreator(signature string) *TokenCreator {
	return &TokenCreator{
		signature: signature,
	}
}

func (tc *TokenCreator) CreateToken(contract []string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		Contract: contract,
	})
	tokenString, err := token.SignedString([]byte(tc.signature))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (tc *TokenCreator) ContractFromToken(t string) ([]string, error) {
	token, err := jwt.ParseWithClaims(t, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error")
		}
		return []byte(tc.signature), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.Contract, nil
	}
	return nil, errors.New("invalid token")
}
