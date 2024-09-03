package http

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const TimeOut = 5 * time.Second

type Server struct {
	server *http.Server
	router *gin.Engine
	// services to used by various HTTP routes
	// ...
}

func NewServer() *Server {
	s := Server{
		server: &http.Server{
			WriteTimeout: TimeOut,
			ReadTimeout:  TimeOut,
			IdleTimeout:  TimeOut,
		},
		router: gin.Default(),
	}

	s.routes()
	s.server.Handler = s.router

	return &s
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
