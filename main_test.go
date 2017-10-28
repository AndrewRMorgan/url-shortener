package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

//Need to change test function to either make a database connect or use a mock one.
func TestCreateURL(t *testing.T) {
	expected := `{"original_url":"https://www.google.com","short_url":"https://morning-retreat-24523.herokuapp.com/get/3578"}`

	router := httprouter.New()
	router.GET("/new/*url", createURL)

	req, _ := http.NewRequest("GET", "/new/https://www.google.com", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: Got %v but want %v",
			rr.Body.String(), expected)
	}
}
