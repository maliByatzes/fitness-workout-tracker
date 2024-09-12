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

func (s *Server) getUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*fwt.User)

		profile, err := s.ProfileService.FindProfileByUserID(c, user.ID)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in get user profile handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"profile": profile,
		})
	}
}

func (s *Server) updateProfile() gin.HandlerFunc {
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
		upd := fwt.ProfileUpdate{}
		if req.Profile.FirstName != "" {
			upd.FirstName = &req.Profile.FirstName
		}
		if req.Profile.LastName != "" {
			upd.LastName = &req.Profile.LastName
		}
		if !req.Profile.DateOfBirth.IsZero() {
			upd.DateOfBirth = &req.Profile.DateOfBirth
		}
		if req.Profile.Gender != "" {
			upd.Gender = &req.Profile.Gender
		}
		if req.Profile.Height != 0 {
			upd.Height = &req.Profile.Height
		}
		if req.Profile.Weight != 0 {
			upd.Weight = &req.Profile.Weight
		}

		profile, err := s.ProfileService.FindProfileByUserID(c, user.ID)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND && fwt.ErrorMessage(err) == "Profile not found." {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in update profile handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		updatedProfile, err := s.ProfileService.UpdateProfile(c, profile.ID, upd)
		if err != nil {
			log.Printf("error in update profile handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "profile updated successfully",
			"profile": updatedProfile,
		})
	}
}

func (s *Server) deleteProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*fwt.User)

		profile, err := s.ProfileService.FindProfileByUserID(c, user.ID)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND && fwt.ErrorMessage(err) == "Profile not found." {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in update profile handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		err = s.ProfileService.DeleteProfile(c, profile.ID)
		if err != nil {
			log.Printf("error in update profile handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "profile deleted successfully",
		})
	}
}
