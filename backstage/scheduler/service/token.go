package service

import (
	"crypto/rsa"

	log "github.com/sirupsen/logrus"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model/apperrors"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/utils"
)

/*
	struct: tokenService
	description: service layer
*/
type tokenService struct {
	PubKey           *rsa.PublicKey
	IDExpirationSecs int64
}

/*
	struct: TSConfig
	description: used for config instance of
			struct tokenService
*/
type TSConfig struct {
	PubKey           *rsa.PublicKey
	IDExpirationSecs int64
}

/*
	func: NewTokenService
	description: create, config and return an instance
			of struct TokenService
*/
func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		PubKey:           c.PubKey,
		IDExpirationSecs: c.IDExpirationSecs,
	}
}

/*
	func: ValidateIDToken
	description: validates the id token jwt string
		It returns the user extract from the IDTokenCustomClaims
		if the provided refresh token is valid
*/
func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := utils.ValidateIDToken(tokenString, s.PubKey)
	if err != nil {
		log.WithFields(log.Fields{
			"ID Token": tokenString,
			"error":    err,
		}).Warn("Unable to validate or parse idToken")
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, err
}
