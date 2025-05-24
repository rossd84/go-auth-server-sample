package main

import (
	"fmt"
	"log"
	"net/http"
	"saas-go-postgres/internal/router"
)

func main() {
	port := ":8080"
	r := router.NewRouter()
	
	fmt.Printf("Running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
