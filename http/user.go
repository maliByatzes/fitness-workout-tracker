package http

import (
	"log"
	"net/http"
	"time"

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
			return
		}

		newUser := fwt.User{
			Username: req.User.Username,
			Email:    req.User.Email,
		}
		newUser.SetPassword(req.User.Password)

		if err := s.userService.CreateUser(c, &newUser); err != nil {
			if fwt.ErrorCode(err) == fwt.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in update user handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"user": newUser,
		})
	}
}

func (s *Server) loginUser() gin.HandlerFunc {
	var req struct {
		User struct {
			Username string `json:"username" binding:"required,min=3"`
			Password string `json:"password" binding:"required,min=8,max=72"`
		} `json:"user" binding:"required"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.userService.Authenticate(c, req.User.Username, req.User.Password)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		accessToken, accessPayload, err := s.tokenMaker.CreateToken(
			user.ID,
			user.Username,
			time.Hour*24,
		)
		if err != nil {
			log.Printf("error in create token in login user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		duration := accessPayload.ExpiredAt.Sub(time.Now())

		c.SetCookie(
			"access_token",
			accessToken,
			int(duration.Seconds()),
			"/",
			"localhost",
			false,
			true)

		c.JSON(http.StatusOK, gin.H{
			"user":         user,
			"access_token": accessToken,
		})
	}
}

func (s *Server) logoutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "logged out successfully",
		})
	}
}

func (s *Server) getCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*fwt.User)

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (s *Server) updateUser() gin.HandlerFunc {
	var req struct {
		User struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		upd := fwt.UserUpdate{}
		if req.User.Username != "" {
			upd.Username = &req.User.Username
		}
		if req.User.Email != "" {
			upd.Email = &req.User.Email
		}

		user := c.MustGet("user").(*fwt.User)

		newUser, err := s.userService.UpdateUser(c, user.ID, upd)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in update user handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user updated successfully",
			"user":    newUser,
		})
	}
}
