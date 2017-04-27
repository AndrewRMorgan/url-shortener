package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

var db *sql.DB
var err error

type UrlResponse struct {
	OriginalURL interface{} `json:"original_url"`
	ShortURL    interface{} `json:"short_url"`
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"database"`
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name"`
}

func main() {
	var config = loadConfig("config.json")

	db, err = sql.Open("mysql", ""+config.Database.User+":"+config.Database.Password+"@tcp("+config.Host+":"+config.Port+")/"+config.Name+"")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := httprouter.New()
	router.GET("/new/*url", createURL) ///
	router.GET("/get/:id", getURL)
	router.GET("/", index)
	http.ListenAndServe(":"+port, router)
}

func loadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	configFile.Close()
	return config
}

func index(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprint(res, "Welcome to the URL Shortener Service!\n")
}

func createURL(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	var response = UrlResponse{}
	var shortUrl string
	var originalUrl string = ps.ByName("url")
	var newUrl string = strings.Replace(originalUrl, "/", "", 1) //This removes the first forward slash: '/'
	reg, _ := regexp.Compile(`^(?:(?:https?|ftp):\/\/)(?:\S+(?::\S*)?@)?(?:(?:!(?:10|127)(?:\.\d{1,3}){3})(?:!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?:!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\x{00a1}-\x{ffff}0-9]-*)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.(?:[a-z\x{00a1}-\x{ffff}0-9]-*)*[a-z\x{00a1}-\x{ffff}0-9]+)*(?:\.(?:[a-z\x{00a1}-\x{ffff}]{2,}))\.?)(?::\d{2,5})?(?:[/?#]\S*)?$`)

	if reg.MatchString(newUrl) {
		err = db.QueryRow("SELECT original_url, short_url FROM urls WHERE original_url = ?", newUrl).Scan(&newUrl, &shortUrl)
		if err != nil {
			idNum := random(0, 9999)
			idStr := strconv.Itoa(idNum)
			check(err)
			shortUrl = "https://morning-retreat-24523.herokuapp.com/" + idStr
			_, err := db.Exec("INSERT INTO urls(id, original_url, short_url) VALUES(?, ?, ?)", idNum, newUrl, shortUrl)
			check(err)
			response = UrlResponse{OriginalURL: newUrl, ShortURL: shortUrl}
		} else {
			response = UrlResponse{OriginalURL: newUrl, ShortURL: shortUrl}
		}
		js, err := json.Marshal(response)
		check(err)
		res.Header().Set("Content-Type", "application/json")
		res.Write(js)
	} else {
		errorResponse := ErrorResponse{Error: "Sorry, wrong url format. Please make sure you have a valid protocol and real site."}
		js, err := json.Marshal(errorResponse)
		check(err)
		res.Header().Set("Content-Type", "application/json")
		res.Write(js)
	}

}

func getURL(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	var idStr string = ps.ByName("id")

	var idNum, err = strconv.Atoi(idStr)
	var originalUrl string
	err = db.QueryRow("SELECT original_url FROM urls WHERE id = ?", idNum).Scan(&originalUrl)
	check(err)

	http.Redirect(res, req, originalUrl, 301) //This needs to have a protocol, i.e. http or https
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
