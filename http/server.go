package http

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/fwt"
	"github.com/maliByatzes/fwt/postgres"
	"github.com/maliByatzes/fwt/token"
)

const TimeOut = 5 * time.Second

type Server struct {
	Server                 *http.Server
	Router                 *gin.Engine
	TokenMaker             token.Maker
	UserService            fwt.UserService
	ProfileService         fwt.ProfileService
	WorkoutService         fwt.WorkoutService
	ExerciseService        fwt.ExerciseService
	WorkoutExerciseService fwt.WorkoutExerciseService
	WEStatusService        fwt.WEStatusService
}

func NewServer(db *postgres.DB, secretKey string) (*Server, error) {
	s := Server{
		Server: &http.Server{
			WriteTimeout: TimeOut,
			ReadTimeout:  TimeOut,
			IdleTimeout:  TimeOut,
		},
		Router: gin.Default(),
	}

	tkMaker, err := token.NewJWTMaker(secretKey)
	if err != nil {
		return nil, err
	}
	s.TokenMaker = tkMaker

	s.routes()
	s.UserService = postgres.NewUserService(db)
	s.ProfileService = postgres.NewProfileService(db)
	s.WorkoutService = postgres.NewWorkoutService(db)
	s.ExerciseService = postgres.NewExerciseService(db)
	s.WorkoutExerciseService = postgres.NewWorkoutExerciseService(db)
	s.WEStatusService = postgres.NewWEStatusService(db)
	s.Server.Handler = s.Router

	return &s, nil
}

func (s *Server) Run(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	s.Server.Addr = port
	log.Printf("🚀 Server starting on port %s", port)
	return s.Server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}

func healthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "available",
		})
	}
}
