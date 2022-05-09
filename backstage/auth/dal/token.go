package dal

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model/apperrors"
)

/*
	struct: redisTokenDAL
	description: dal layer
*/
type redisTokenDAL struct {
	Redis *redis.Client
}

/*
	func: NewTokenDAL
	description: return an instance of struct redisTokenDAL
*/
func NewTokenDAL(redisClient *redis.Client) model.TokenDAL {
	return &redisTokenDAL{
		Redis: redisClient,
	}
}

/*
	func: SetRefreshToken
	description: SET new refresh token to redis
*/
func (r *redisTokenDAL) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.WithFields(log.Fields{
			"User ID":  userID,
			"Token ID": tokenID,
			"error":    err,
		}).Warn("Failed to set refresh token to redis")
		return apperrors.NewInternal()
	}

	return nil
}

/*
	func: DeleteRefreshToken
	description: DELETE specified refresh token in redis
*/
func (r *redisTokenDAL) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.WithFields(log.Fields{
			"User ID":  userID,
			"Token ID": tokenID,
			"error":    err,
		}).Warn("Failed to delete refresh token")
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.WithFields(log.Fields{
			"User ID":  userID,
			"Token ID": tokenID,
		}).Warn("Refresh token to redis doesn't exist, failed to delete refresh token")
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

/*
	func: DeleteRefreshTokens
	description: DELETE all refresh tokens related to a user in redis
*/
func (r *redisTokenDAL) DeleteRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.WithFields(log.Fields{
				"Refresh Token": iter.Val(),
				"error":         err,
			}).Warn("Failed to delete refresh token")
			failCount++
		}
	}

	if err := iter.Err(); err != nil {
		log.WithFields(log.Fields{
			"Refresh Token": iter.Val(),
			"error":         err,
		}).Warn("Failed to delete refresh token")
	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}

	return nil
}
