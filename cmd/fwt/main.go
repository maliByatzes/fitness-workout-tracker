package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/maliByatzes/fwt/http"
	"github.com/maliByatzes/fwt/postgres"
)

type config struct {
	port      string
	dbURL     string
	jwtSecret string
}

func main() {
	cfg := envConfig()

	db := postgres.NewDB(cfg.dbURL)
	if err := db.Open(); err != nil {
		log.Fatalf("cannot open database: %v", err)
	}
	defer db.Close()

	srv, err := http.NewServer(db, cfg.jwtSecret)
	if err != nil {
		log.Fatalf("cannot create new server: %v", err)
	}
	defer srv.Close()
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

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		panic("JWT_SECRET is not set!")
	}

	return config{port: port, dbURL: dbURL, jwtSecret: jwtSecret}
}
