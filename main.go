package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", Router)
	http.ListenAndServe(":"+port, nil)
}

func Router(w http.ResponseWriter, r *http.Request) {

}
