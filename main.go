package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func init() {
	db, err = spl.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", Router)
	http.HandleFunc("/new/", AddURL)
	http.ListenAndServe(":"+port, nil)
}

func Router(w http.ResponseWriter, r *http.Request) {

}

func GetURL(w http.ResponseWriter, r *http.Request) {
	var shortUrl string
	err = db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortUrl)
}

func AddURL(w http.ResponseWriter, r *http.Request) {

}
