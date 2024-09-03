package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/maliByatzes/fwt/http"
)

type config struct {
	port  string
	dbURL string
}

func main() {
	cfg := envConfig()

	srv := http.NewServer()
	log.Fatal(srv.Run(cfg.port))
}

func envConfig() config {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("PORT is not set!")
	}

	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("DATABASE_URL is not set!")
	}

	return config{port: port, dbURL: dbURL}
}
