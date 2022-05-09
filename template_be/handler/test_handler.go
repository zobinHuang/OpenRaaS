package handler

import (
	"business/model"
	"business/model/apperrors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/*
	struct: rdmReq
	description: format of html body (json) in http request that sent to endpoint "/api/business/test"
*/
type rdmReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
	func: Test
	description: handler for endpoint "/api/business/test"
*/
func (h *Handler) Test(c *gin.Context) {
	// bind request
	var req rdmReq
	if ok := bindData(c, &req); !ok {
		return
	}

	// initialize model
	rdbm := &model.RDbModel{}
	rdbm.UserName = req.Username
	rdbm.Password = req.Password

	// invoke related service
	ctx := c.Request.Context()
	err := h.RDbService.GetRDbModel(ctx, rdbm)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to get test model")
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
	}

	// return http_ok if success
	c.JSON(http.StatusOK, gin.H{})
}
