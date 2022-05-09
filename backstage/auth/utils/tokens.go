package utils

import (
	"crypto/rsa"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model"
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
	func: GenerateIDToken
	description: generate a jwt id token according to
		giving user and rsa private key
*/
func GenerateIDToken(u *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp // 15 minutes from current unixTime

	// payload
	claims := IDTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	// create token header and payload
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// signed token and get toke final string
	ss, err := token.SignedString(key)
	if err != nil {
		log.Warn("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

/*
	struct: RefreshTokenData
	description: Struct to descripte a refresh token, contains a
		refresh token string, a refresh token id and expire
		timestamp
*/
type RefreshTokenData struct {
	SS        string
	ID        uuid.UUID
	ExpiresIn time.Duration
}

/*
	struct: RefreshTokenCustomClaims
	description: Claims for refresh token, contains basic jwt claims,
		and a user id
*/
type RefreshTokenCustomClaims struct {
	UID uint64 `json:"uid"`
	jwt.StandardClaims
}

/*
	func: GenerateRefreshToken
	description: generate a jwt refresh token according to
		giving user id and hmac key string
*/
func GenerateRefreshToken(uid uint64, key string, exp int64) (*RefreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Warn("Failed to generate refresh token ID")
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		log.Warn("Failed to sign refresh token string")
		return nil, err
	}

	return &RefreshTokenData{
		SS:        ss,
		ID:        tokenID,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
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

/*
	func: ValidateRefreshToken
	description:  returns the refresh-token's claims if the token is valid
*/
func ValidateRefreshToken(tokenString string, key string) (*RefreshTokenCustomClaims, error) {
	claims := &RefreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("refresh token is invalid")
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("refresh token valid but couldn't parse claims")
	}

	return claims, nil
}
