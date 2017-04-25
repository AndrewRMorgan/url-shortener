package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"math/rand"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	check(err)

	defer db.Close()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/new", new)
	router.POST("/new/:url", createURL)
	router.GET("/:id", getURL)
	router.GET("/favicon.ico", ) //Don't quite know what should be here.
	http.ListenAndServe(":"+port, router)
}

func index(w http.ResponseWriter, r *http.Request) {

}

func new(w http.ResponseWriter, r *http.Request) {

}

func getURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := JsonResponse

	idStr := ps.ByName("id")
	idNum := strconv.Atoi(idStr)
	var originalUrl string
	var shortUrl string
	err = db.QueryRow("SELECT original_url, short_url FROM urls WHERE id = ?", idNum).Scan(&originalUrl, &shortUrl)
	check(err)

	response = JsonResponse{OriginalURL: originalURL, ShortURL: shortUrl}

	js, err := json.Marshal(reponse)
	check(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	originalUrl := vars["url"]
	idNum := random(0, 9999)
	idStr := strconv.Itoa(idNum)
	shortUrl := "https://morning-retreat-24523.herokuapp.com/" + idStr

	res, err := db.Exec("INSERT INTO urls(id, original_url, short_url) VALUES(?, ?, ?)", idNum, originalUrl, shortUrl)
	check(err)

	GetURLHandler(shortUrl)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
