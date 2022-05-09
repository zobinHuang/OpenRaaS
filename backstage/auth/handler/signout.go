package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model/apperrors"
)

/*
	struct: signoutReq
	description: format of html body (json) in http request
			that sent to endpoint "/api/account/signout"
*/
type signoutReq struct {
}

/*
	func: SignOut
	description: handler for endpoint "/api/account/signout"
*/
func (h *Handler) SignOut(c *gin.Context) {
	user := c.MustGet("user")

	var req signoutReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	if err := h.TokenService.Signout(ctx, user.(*model.User).Id); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully",
	})
}
