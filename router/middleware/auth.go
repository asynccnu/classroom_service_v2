package middleware

import (
	"github.com/asynccnu/classroom_service_v2/handler"
	"github.com/asynccnu/classroom_service_v2/pkg/errno"
	"github.com/asynccnu/classroom_service_v2/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
