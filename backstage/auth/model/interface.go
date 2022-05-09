package model

import (
	"context"
	"time"
)

/*
	interface: UserService
	description: interface of service layer for user authetication
*/
type UserService interface {
	Signin(ctx context.Context, u *User) (*User, error)
	Signup(ctx context.Context, u *User) error
}

/*
	interface: TokenService
	description: interface of token service for
			user authorization
*/
type TokenService interface {
	NewPairFromUser(ctx context.Context, up *User, prevTokenID string) (*TokenPair, error)
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error)
	Signout(ctx context.Context, uid uint64) error
}

/*
	interface: UserDAL
	description: interface of data access layer for user authetication
*/
type UserDAL interface {
	CreateNewUser(ctx context.Context, u *User) error
	FindUserProfileByEmail(ctx context.Context, email string) (*User, error)
}

/*
	interface: TokenDAL
	description: interface of data access layer
			for token processing
*/
type TokenDAL interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
	DeleteRefreshTokens(ctx context.Context, userID string) error
}
