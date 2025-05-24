package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"saas-go-postgres/internal/router"
	"saas-go-postgres/internal/config"
	"saas-go-postgres/internal/user"
)

func main() {
	port := ":8080"
	r := router.NewRouter()
	
	dbConfig := config.LoadDBConfig()
	db, err := sql.Open("postgres", dbConfig.DSN())
	if err != nil {
		log.Fatal(err)
	}

	user.DB = db

	fmt.Printf("Running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
