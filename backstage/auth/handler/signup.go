package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model/apperrors"
)

/*
	struct: rdmReq
	description: format of html body (json) in http request that sent to endpoint "/api/github.com/zobinHuang/OpenRaaS/backstage/auth/test"
*/
type signupReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
	func: Signup
	description: handler for endpoint "/api/auth/signup"
*/
func (h *Handler) Signup(c *gin.Context) {
	req := &signupReq{}

	if ok := bindData(c, &req); !ok {
		return
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	err := h.UserService.Signup(ctx, user)
	if err != nil {
		log.WithFields(log.Fields{
			"Email": user.Email,
			"error": err,
		}).Warn("Failed to sign up user")

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// ctx = c.Request.Context()
	// tokens, err := h.TokenService.NewPairFromUser(ctx, user, "")
	// if err != nil {
	// 	log.Printf("Failed to create tokens for user: %v\n", err.Error())
	// 	c.JSON(apperrors.Status(err), gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{
		/* Temp nothing */
	})
}
