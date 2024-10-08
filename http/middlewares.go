package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/fwt"
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

		payload, err := s.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Unauthorized - %v", err),
			})
			c.Abort()
			return
		}

		user, err := s.UserService.FindUserByID(c, payload.ID)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Unauthorized - %v", err),
			})
			c.Abort()
			return
		}

		ctx := fwt.NewContextWithUser(c.Request.Context(), user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
