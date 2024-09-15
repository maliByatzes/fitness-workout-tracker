package http

import (
	"log"
	"net/http"
	"strconv"
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

func (s *Server) getAllWorkouts() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*fwt.User)

		workouts, n, err := s.WorkoutService.FindWorkouts(c, fwt.WorkoutFilter{UserID: &user.ID})
		if err != nil {
			log.Printf("error in create workout handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    n,
			"workouts": workouts,
		})
	}
}

func (s *Server) getOneWorkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		workoutIDstr := c.Param("id")
		workoutID, err := strconv.ParseUint(workoutIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid workout id param",
			})
			return
		}

		user := c.MustGet("user").(*fwt.User)

		workoutID2 := uint(workoutID)
		workout, err := s.WorkoutService.FindWorkoutByIDUserID(c, workoutID2, user.ID)
		if err != nil {
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

		c.JSON(http.StatusOK, gin.H{
			"workout": workout,
		})
	}
}
