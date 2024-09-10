package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/fwt"
)

func (s *Server) createProfile() gin.HandlerFunc {
	var req struct {
		Profile struct {
			FirstName   string    `json:"first_name"`
			LastName    string    `json:"last_name"`
			DateOfBirth time.Time `json:"dob"`
			Gender      string    `json:"gender"`
			Height      float64   `json:"height"`
			Weight      float64   `json:"weight"`
		} `json:"profile"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user := c.MustGet("user").(*fwt.User)
		newProfile := fwt.Profile{
			UserID:      user.ID,
			FirstName:   req.Profile.FirstName,
			LastName:    req.Profile.LastName,
			DateOfBirth: req.Profile.DateOfBirth,
			Gender:      req.Profile.Gender,
			Height:      req.Profile.Height,
			Weight:      req.Profile.Weight,
		}

		if err := s.ProfileService.CreateProfile(c, &newProfile); err != nil {
			if fwt.ErrorCode(err) == fwt.EINVALID {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}

			if fwt.ErrorCode(err) == fwt.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}

			log.Printf("error in create profile handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"profile": newProfile,
		})
	}
}
