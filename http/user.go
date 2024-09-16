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

		if err := s.UserService.CreateUser(c.Request.Context(), &newUser); err != nil {
			if fwt.ErrorCode(err) == fwt.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in create user handler: %v", err)
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

		user, err := s.UserService.Authenticate(c.Request.Context(), req.User.Username, req.User.Password)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		accessToken, accessPayload, err := s.TokenMaker.CreateToken(
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
		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

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

		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		newUser, err := s.UserService.UpdateUser(c.Request.Context(), user.ID, upd)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}

			if fwt.ErrorCode(err) == fwt.ENOTAUTHORIZED {
				c.JSON(http.StatusUnauthorized, gin.H{
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

func (s *Server) deleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		err := s.UserService.DeleteUser(c.Request.Context(), user.ID)
		if err != nil {
			log.Printf("error in delete user handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}

		c.SetCookie("access_token", "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "user deleted successfully",
		})
	}
}
