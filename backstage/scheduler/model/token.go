package model

import "github.com/google/uuid"

/*
	model: TokenPair
	description: token used for user authorization
*/
type TokenPair struct {
	IDToken
	RefreshToken
}

type RefreshToken struct {
	ID  uuid.UUID `json:"-"`            // Refresh Token ID
	UID uint64    `json:"-"`            // User ID
	SS  string    `json:"refreshToken"` // Sign String
}

type IDToken struct {
	SS string `json:"idToken"` // Sign String
}
