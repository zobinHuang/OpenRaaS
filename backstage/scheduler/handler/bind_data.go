package handler

import (
	"fmt"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

/*
	struct: invalidArgument
	description: used to record invalid arguments
		which detected during binding stage
*/
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

/*
	func: bindData
	description: bind incoming data to specified struct
*/
func bindData(c *gin.Context, req interface{}) bool {
	// check content type
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())
		err := apperrors.NewUnsupportedMediaType(msg)

		c.JSON(err.Status(), gin.H{
			"errors": err,
		})

		return false
	}

	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(&req); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Error binding data")

		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			err := apperrors.NewBadRequest("Invalid request parameters, See invalidArgs")

			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})

			return false
		}

		fallback := apperrors.NewInternal()

		c.JSON(fallback.Status(), gin.H{
			"error": fallback,
		})

		return false
	}

	return true
}
