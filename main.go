package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *sql.DB

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env", err.Error())
	}

	// get environment variables
	server := os.Getenv("SQL_HOST")
	port := 1433
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PASSWORD")
	database := os.Getenv("SQL_DATABASE")

	r := mux.NewRouter().StrictSlash(true)
	handler := cors.AllowAll().Handler(r)
	srv := &http.Server{
		Addr:    ":8008",
		Handler: handler,
	}

	// routes
	r.HandleFunc("/api/v2/dept/{dept}", getQueue)

	// set up database
	conn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)

	db, err = sql.Open("sqlserver", conn)
	if err != nil {
		log.Fatal("Error connection to database: ", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not connect to database: ", err.Error())
	}
	// connect
	fmt.Println("Server Running")

	log.Fatal(srv.ListenAndServe())
}
