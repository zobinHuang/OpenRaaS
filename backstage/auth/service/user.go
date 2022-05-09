package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/zobinHuang/BrosCloud/backstage/auth/model"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model/apperrors"
	"github.com/zobinHuang/BrosCloud/backstage/auth/utils"
)

/*
	struct: userService
	description: service layer
*/
type userService struct {
	UserDAL model.UserDAL
}

/*
	struct: UserServiceConfig
	description: used for config instance of struct rdbService
*/
type UserServiceConfig struct {
	UserDAL model.UserDAL
}

/*
	func: NewRDbService
	description: create, config and return an instance of struct rdbService
*/
func NewUserService(c *UserServiceConfig) model.UserService {
	return &userService{
		UserDAL: c.UserDAL,
	}
}

/*
	func: Signin
	description: signin authentication service
*/
func (us *userService) Signin(ctx context.Context, u *model.User) (*model.User, error) {
	uFetched, err := us.UserDAL.FindUserProfileByEmail(ctx, u.Email)
	if err != nil {
		return nil, apperrors.NewAuthorization(fmt.Sprintf("No such user with email: %v", u.Email))
	}

	match, err := utils.ComparePassword(uFetched.Password, u.Password)
	if err != nil {
		return nil, apperrors.NewInternal()
	}
	if !match {
		return nil, apperrors.NewAuthorization("Invalid email and password combination")
	}

	return uFetched, nil
}

/*
	func: Signup
	description: service that used for signing up a new user
*/
func (us *userService) Signup(ctx context.Context, u *model.User) error {
	// update password field to hash password
	pw, err := utils.HashPassword(u.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"Email": u.Email,
			"error": err,
		}).Warn("Failed to hash password for the new user")
		return apperrors.NewInternal()
	}

	u.Password = pw

	if err := us.UserDAL.CreateNewUser(ctx, u); err != nil {
		return err
	}

	return nil
}
