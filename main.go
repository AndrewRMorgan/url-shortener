package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"math/rand"
	"strconv"

	"github.com/gorilla/mux"
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

	defer db.Close()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/new/", NewHandler)
	r.HandleFunc("/new/{url}", CreateURLHandler)
	r.HandleFunc("/{id}", GetURLHandler)
	http.Handle("/", r)

	http.ListenAndServe(":"+port, nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(w http.ResponseWriter, r *http.Request) {

}

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var shortUrl string = vars["id"]
	var originalUrl string
	err = db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortUrl).Scan(&originalUrl)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(originalUrl)
}

func CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	originalURL := vars["url"]
	randNum := random(0, 9999)
	id := strconv.Itoa(randNum)
	shortURL := "https://morning-retreat-24523.herokuapp.com/" + id

	res, err := db.Exec("INSERT INTO urls(original_url, short_url) VALUES(?, ?)", originalURL, shortURL)
	if err != nil {
		fmt.Println(err)
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}
