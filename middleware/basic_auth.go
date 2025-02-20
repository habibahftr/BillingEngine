package middleware

import "github.com/gin-gonic/gin"

func BasicAuth(
	username, pass string,
) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		username: pass,
	})

}
