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
	server      *http.Server
	router      *gin.Engine
	tokenMaker  token.Maker
	userService fwt.UserService
}

func NewServer(db *postgres.DB, secretKey string) (*Server, error) {
	s := Server{
		server: &http.Server{
			WriteTimeout: TimeOut,
			ReadTimeout:  TimeOut,
			IdleTimeout:  TimeOut,
		},
		router: gin.Default(),
	}

	tkMaker, err := token.NewJWTMaker(secretKey)
	if err != nil {
		return nil, err
	}
	s.tokenMaker = tkMaker

	s.routes()
	s.userService = postgres.NewUserService(db)
	s.server.Handler = s.router

	return &s, nil
}

func (s *Server) Run(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	s.server.Addr = port
	log.Printf("🚀 Server starting on port %s", port)
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func healthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "available",
		})
	}
}
