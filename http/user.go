package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/fwt"
)

func (s *Server) createUser() gin.HandlerFunc {
	var req struct {
		User struct {
			Username string `json:"username" binding:"required,min=3"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8,max=72"`
		} `json:"user" binding:"required"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		newUser := fwt.User{
			Username: req.User.Username,
			Email:    req.User.Email,
		}
		newUser.SetPassword(req.User.Password)

		if err := s.userService.CreateUser(c, &newUser); err != nil {
			switch {
			case fwt.ErrorCode(err) == fwt.ECONFLICT && fwt.ErrorMessage(err) == "This username already exists.":
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
			case fwt.ErrorCode(err) == fwt.ECONFLICT && fwt.ErrorMessage(err) == "This email already exists.":
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
			default:
				log.Printf("Error in createUser: %v\n", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"user": newUser,
		})
	}
}
