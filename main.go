package main

import (
	"encoding/json"
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

type JsonResponse struct {
	OriginalURL interface{} `json:"original_url"`
	ShortURL interface{} `json:"short_url"`
}

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

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(w http.ResponseWriter, r *http.Request) {

}

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	response := JsonResponse

	vars := mux.Vars(r)
	idStr := vars["id"]
	idNum := strconv.Atoi(idStr)
	var originalUrl string
	var shortUrl string
	err = db.QueryRow("SELECT original_url, short_url FROM urls WHERE id = ?", idNum).Scan(&originalUrl, &shortUrl)
	if err != nil {
		fmt.Println(err)
	}

	response = JsonResponse{OriginalURL: originalURL, ShortURL: shortUrl}

	js, err := json.Marshal(reponse)
	if err != nil {
		fmt.Println("Json Marshal returned nil")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	originalUrl := vars["url"]
	idNum := random(0, 9999)
	idStr := strconv.Itoa(idNum)
	shortUrl := "https://morning-retreat-24523.herokuapp.com/" + idStr

	res, err := db.Exec("INSERT INTO urls(id, original_url, short_url) VALUES(?, ?, ?)", idNum, originalUrl, shortUrl)
	if err != nil {
		fmt.Println(err)
	}

	GetURLHandler(shortUrl)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}
