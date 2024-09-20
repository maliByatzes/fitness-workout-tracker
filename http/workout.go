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

		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		exercises := make([]*fwt.Exercise, 0)
		for _, exName := range req.Workout.Exercises {
			exercise, err := s.ExerciseService.FindExerciseByName(c.Request.Context(), exName)
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

			exercises = append(exercises, exercise)
		}

		newWorkout := fwt.Workout{
			UserID:        user.ID,
			Name:          req.Workout.Name,
			ScheduledDate: req.Workout.ScheduledDate,
			Exercises:     exercises,
		}

		if err := s.WorkoutService.CreateWorkout(c.Request.Context(), &newWorkout); err != nil {
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

			if fwt.ErrorCode(err) == fwt.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
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
		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		workouts, n, err := s.WorkoutService.FindWorkouts(c.Request.Context(), fwt.WorkoutFilter{UserID: &user.ID})
		if err != nil {
			log.Printf("error in get all workouts workout handler: %v", err)
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

		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		workoutID2 := uint(workoutID)
		workout, err := s.WorkoutService.FindWorkoutByIDUserID(c.Request.Context(), workoutID2, user.ID)
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

func (s *Server) updateWorkout() gin.HandlerFunc {
	var req struct {
		Workout struct {
			Name          string    `json:"name"`
			ScheduledDate time.Time `json:"scheduled_date"`
		} `json:"workout"`
	}

	return func(c *gin.Context) {
		workoutIDstr := c.Param("id")
		workoutID, err := strconv.ParseUint(workoutIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid workout id param",
			})
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user := fwt.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		upd := fwt.WorkoutUpdate{}
		if req.Workout.Name != "" {
			upd.Name = &req.Workout.Name
		}
		if !req.Workout.ScheduledDate.IsZero() {
			upd.ScheduledDate = &req.Workout.ScheduledDate
		}

		workout, err := s.WorkoutService.UpdateWorkout(c.Request.Context(), uint(workoutID), upd)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTAUTHORIZED {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}

			log.Printf("error in update workout handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "workout updated successfully",
			"workout": workout,
		})
	}
}

func (s *Server) removeExercisesFromWorkout() gin.HandlerFunc {
	var req struct {
		Exercises []string `json:"exercises"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		workoutIDstr := c.Param("id")
		workoutID, err := strconv.ParseUint(workoutIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid workout id param",
			})
			return
		}

		w, err := s.WorkoutService.FindWorkoutByID(c.Request.Context(), uint(workoutID))
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": fwt.ErrorMessage(err),
				})
				return
			}
			log.Printf("error in remove exercises from workout handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		if len(req.Exercises) >= len(w.Exercises) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "There must be at least one exercise remaining in the workout.",
			})
			return
		}

		workout, err := s.WorkoutService.RemoveExercisesFromWorkout(c.Request.Context(), w.ID, req.Exercises)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
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

			log.Printf("error in remove exercises from workout handler: %v", err)
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

func (s *Server) addExercisesToWorkout() gin.HandlerFunc {
	var req struct {
		Exercises []string `json:"exercises"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		workoutIDstr := c.Param("id")
		workoutID, err := strconv.ParseUint(workoutIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid workout id param",
			})
			return
		}

		workout, err := s.WorkoutService.AddExercisesToWorkout(c.Request.Context(), uint(workoutID), req.Exercises)
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
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

			log.Printf("error in add exercises to workout handler: %v", err)
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

func (s *Server) deleteWorkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		workoutIDstr := c.Param("id")
		workoutID, err := strconv.ParseUint(workoutIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid workout id param",
			})
			return
		}

		err = s.WorkoutService.DeleteWorkout(c.Request.Context(), uint(workoutID))
		if err != nil {
			if fwt.ErrorCode(err) == fwt.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
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

			log.Printf("error in delete workout handler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "workout deleted successfully",
		})
	}
}
