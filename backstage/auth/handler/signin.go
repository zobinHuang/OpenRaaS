package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model/apperrors"
)

/*
	struct: signinReq
	description: format of html body (json) in http request that sent to endpoint "/api/auth/signin"
*/
type signinReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
	func: Signin
	description: handler for endpoint "/api/auth/signin"
*/
func (h *Handler) Signin(c *gin.Context) {
	var req signinReq
	if ok := bindData(c, &req); !ok {
		return
	}

	user := &model.User{}
	user.Email = req.Email
	user.Password = req.Password

	ctx := c.Request.Context()
	upFetched, err := h.UserService.Signin(ctx, user)
	if err != nil {
		log.WithFields(log.Fields{
			"Email": user.Email,
			"error": err,
		}).Warn("Failed to signin user")

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, upFetched, "")
	if err != nil {
		log.WithFields(log.Fields{
			"Email": upFetched.Email,
			"error": err,
		}).Warn("Failed to create tokens for user")

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens":   tokens,
		"username": upFetched.Email,
	})
}
