package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/your-username/storage-service/internal/api"
)

const portNum string = ":8080"

func main() {
	log.Println("Starting API server.")

	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	log.Println("DATABASE_URL:", dsn)
	if dsn == "" {
		log.Fatal("DATABASE_URL is required (e.g. postgres://user:pass@localhost:5432/storage?sslmode=disable)")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	router := api.NewRouter(db)

	server := &http.Server{
		Addr:    portNum,
		Handler: router,
	}

	log.Println("Started on http://localhost" + portNum)
	log.Println("To close connection CTRL+C")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
