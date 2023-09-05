package utils

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
	struct: IDTokenCustomClaims
	description: Claims for id token,
		contains basic jwt claims
		and user information
*/
type IDTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

/*
	func: ValidateIDToken
	description: returns the token's claims if the token is valid
*/
func ValidateIDToken(tokenString string, key *rsa.PublicKey) (*IDTokenCustomClaims, error) {
	claims := &IDTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID Token is invalid")
	}

	claims, ok := token.Claims.(*IDTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}
