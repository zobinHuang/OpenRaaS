package service

import (
	"context"
	"crypto/rsa"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model/apperrors"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/utils"
)

/*
	struct: tokenService
	description: service layer
*/
type tokenService struct {
	TokenDAL              model.TokenDAL
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

/*
	struct: TSConfig
	description: used for config instance of
			struct tokenService
*/
type TSConfig struct {
	TokenDAL              model.TokenDAL
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

/*
	func: NewHandler
	description: create, config and return an instance
			of struct TokenService
*/
func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenDAL:              c.TokenDAL,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

/*
	func: NewPairFromUser
	description: service that return a token pair for a user
*/
func (s *tokenService) NewPairFromUser(ctx context.Context, up *model.User, prevTokenID string) (*model.TokenPair, error) {
	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := s.TokenDAL.DeleteRefreshToken(ctx, strconv.FormatUint(up.Id, 10), prevTokenID); err != nil {
			log.WithFields(log.Fields{
				"User ID":  strconv.FormatUint(up.Id, 10),
				"Token ID": prevTokenID,
				"error":    err,
			}).Warn("Could not delete previous refreshToken")
			return nil, err
		}
	}

	idToken, err := utils.GenerateIDToken(up, s.PrivKey, s.IDExpirationSecs)
	if err != nil {
		log.WithFields(log.Fields{
			"User ID": strconv.FormatUint(up.Id, 10),
			"error":   err,
		}).Warn("Error generating id token")
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := utils.GenerateRefreshToken(up.Id, s.RefreshSecret, s.RefreshExpirationSecs)
	if err != nil {
		log.WithFields(log.Fields{
			"User ID": strconv.FormatUint(up.Id, 10),
			"error":   err,
		}).Warn("Error generating refresh token")
		return nil, apperrors.NewInternal()
	}

	// store refresh token
	if err := s.TokenDAL.SetRefreshToken(ctx, strconv.FormatUint(up.Id, 10), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.WithFields(log.Fields{
			"User ID": strconv.FormatUint(up.Id, 10),
			"error":   err,
		}).Warn("Error storing tokenID")
		return nil, apperrors.NewInternal()
	}

	return &model.TokenPair{
		IDToken: model.IDToken{SS: idToken},
		RefreshToken: model.RefreshToken{
			SS:  refreshToken.SS,
			ID:  refreshToken.ID,
			UID: up.Id,
		},
	}, nil
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

/*
	func: ValidateRefreshToken
	description: validates the refresh token jwt string
		It returns the RefreshToken struct if the provided refresh
		token is valid
*/
func (s *tokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error) {
	claims, err := utils.ValidateRefreshToken(refreshTokenString, s.RefreshSecret)
	if err != nil {
		log.WithFields(log.Fields{
			"Refresh Token String": refreshTokenString,
			"error":                err,
		}).Warn("Unable to validate or parse refreshToken for token string")
		return nil, apperrors.NewAuthorization("Unable to Verify user from refresh token")
	}

	tokenID, err := uuid.Parse(claims.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"Claim ID": claims.Id,
			"error":    err,
		}).Warn("Claims ID could not be parsed as UUID")
		return nil, apperrors.NewAuthorization("Unable to Verify user from refresh token")
	}

	return &model.RefreshToken{
		SS:  refreshTokenString,
		ID:  tokenID,
		UID: claims.UID,
	}, nil
}

/*
	func: Signout
	description: used for signing out a user
*/
func (s *tokenService) Signout(ctx context.Context, uid uint64) error {
	return s.TokenDAL.DeleteRefreshTokens(ctx, strconv.FormatUint(uid, 10))
}
