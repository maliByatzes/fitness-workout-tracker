package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string

		if v := c.GetHeader("Authorization"); strings.HasPrefix(v, "Bearer ") {
			accessToken = strings.TrimPrefix(v, "Bearer ")
		} else {
			accessToken, _ = c.Cookie("access_token")
		}

		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - No access token",
			})
			c.Abort()
			return
		}

		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Unauthorized - %v", err),
			})
			c.Abort()
			return
		}

		user, err := s.userService.FindUserByID(c, payload.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Unauthorized - %v", err),
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
