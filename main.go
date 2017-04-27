package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

	db, err = sql.Open("mysql", "" + config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Name + "")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	router.SkipClean(true)
	router.HandleFunc("/new/", index)
	router.HandleFunc("/", index)
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

func index(res http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.Replace(req.URL.Path, ":/", "://", -1)
	fmt.Println(req.URL.Path)
	param := strings.Split(req.URL.Path, "/")
	if len(param) > 4 {
		createURL(res, req, param)
	} else {
		fmt.Fprint(res, "Welcome to the URL Shortener Service!\n")
	}
}

func createURL(res http.ResponseWriter, req *http.Request, param []string) {

	fmt.Println("createURL")

	posURL := strings.Split(req.URL.Path, "new/")
	originalUrl := posURL[1]

	reg, _ := regexp.Compile(`https?:\/\/(www\.)[a-zA-Z0-9_\-]+\.[(com|net|org|edu|gov|mil|aero|asia|biz|cat|coop|info|int|jobs|mobi|museum|name|post|pro|tel|travel|xxx|ac|ad|ae|af|ag|ai|al|am|an|ao|aq|ar|as|at|au|aw|ax|az|ba|bb|bd|be|bf|bg|bh|bi|bj|bm|bn|bo|br|bs|bt|bv|bw|by|bz|ca|cc|cd|cf|cg|ch|ci|ck|cl|cm|cn|co|cr|cs|cu|cv|cx|cy|cz|dd|de|dj|dk|dm|do|dz|ec|ee|eg|eh|er|es|et|eu|fi|fj|fk|fm|fo|fr|ga|gb|gd|ge|gf|gg|gh|gi|gl|gm|gn|gp|gq|gr|gs|gt|gu|gw|gy|hk|hm|hn|hr|ht|hu|id|ie|il|im|in|io|iq|ir|is|it|je|jm|jo|jp|ke|kg|kh|ki|km|kn|kp|kr|kw|ky|kz|la|lb|lc|li|lk|lr|ls|lt|lu|lv|ly|ma|mc|md|me|mg|mh|mk|ml|mm|mn|mo|mp|mq|mr|ms|mt|mu|mv|mw|mx|my|mz|na|nc|ne|nf|ng|ni|nl|no|np|nr|nu|nz|om|pa|pe|pf|pg|ph|pk|pl|pm|pn|pr|ps|pt|pw|py|qa|re|ro|rs|ru|rw|sa|sb|sc|sd|se|sg|sh|si|sj|Ja|sk|sl|sm|sn|so|sr|ss|st|su|sv|sx|sy|sz|tc|td|tf|tg|th|tj|tk|tl|tm|tn|to|tp|tr|tt|tv|tw|tz|ua|ug|uk|us|uy|uz|va|vc|ve|vg|vi|vn|vu|wf|ws|ye|yt|yu|za|zm|zw)]\/?[/a-zA-Z0-9_\-]+$`)

	if reg.MatchString(originalUrl) {
		idNum := random(0, 9999)
		idStr := strconv.Itoa(idNum)
		check(err)
		shortUrl := "https://morning-retreat-24523.herokuapp.com/" + idStr
		_, err := db.Exec("INSERT INTO urls(id, original_url, short_url) VALUES(?, ?, ?)", idNum, originalUrl, shortUrl)
		check(err)
		response := UrlResponse{OriginalURL: originalUrl, ShortURL: shortUrl}
		js, err := json.Marshal(response)
		check(err)
		res.Header().Set("Content-Type", "application/json")
		res.Write(js)
	} else {
		response := ErrorResponse{Error: "Wrong url format, make sure you have a valid protocol and real site."}
		js, err := json.Marshal(response)
		check(err)
		res.Header().Set("Content-Type", "application/json")
		res.Write(js)
	}


}

func getURL(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var idStr string = vars["id"]
	var idNum, err = strconv.Atoi(idStr)
	var originalUrl string
	err = db.QueryRow("SELECT original_url FROM urls WHERE id = ?", idNum).Scan(&originalUrl)
	check(err)

	http.Redirect(w, r, originalUrl, 301) //This needs to have a protocol: http or https
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
