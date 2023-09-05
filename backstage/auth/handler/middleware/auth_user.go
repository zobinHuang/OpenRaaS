package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model"
	"github.com/zobinHuang/OpenRaaS/backstage/auth/model/apperrors"
)

/*
	struct: authHeader
	description: expected http header: Authorization
*/
type authHeader struct {
	IDToken string `header:"Authorization"`
}

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
	func:	AuthUser
	description: AuthUser extracts a user from the Authorization header
			which is of the form "Bearer token"
			It sets the user to the context if the user exists
*/
func AuthUser(s model.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {

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

				err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")
				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			err := apperrors.NewInternal()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")
		if len(idTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		user, err := s.ValidateIDToken(idTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
