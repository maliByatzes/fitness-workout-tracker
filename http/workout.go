package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/fwt"
)

func (s *Server) createWorkout() gin.HandlerFunc {
	var req struct {
		Workout struct {
			Name          string    `json:"name"`
			ScheduledDate time.Time `json:"scheduled_date"`
			Exercises     []string  `json:"exercises"`
		} `json:"workout"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user := c.MustGet("user").(*fwt.User)
		newWorkout := fwt.Workout{
			UserID:        user.ID,
			Name:          req.Workout.Name,
			ScheduledDate: req.Workout.ScheduledDate,
		}

		if err := s.WorkoutService.CreateWorkout(c, &newWorkout, req.Workout.Exercises); err != nil {
			if fwt.ErrorCode(err) == fwt.EINVALID {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}

			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}

			log.Printf("error in create workout handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"workout": newWorkout,
		})
	}
}
