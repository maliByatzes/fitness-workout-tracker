package http

import (
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := getAllowedOrigins()
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		if slices.Contains(allowedOrigins, origin) || origin == "" {
			ctx.Header("Access-Control-Allow-Origin", origin)
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Origin not allowed",
			})
			return
		}

		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func getAllowedOrigins() []string {
	log.Print(os.Getenv("ALLOWED_ORIGINS"))
	return strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
}
