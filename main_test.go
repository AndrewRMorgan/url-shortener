package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestCreateURL(t *testing.T) {
	expected := `{"original_url":"https://www.google.com", "short_url":"https://morning-retreat-24523.herokuapp.com/get/3578"}`
	handler := createURL
	router := httprouter.New()
	router.GET("/new/*url", handler)

	req, err := http.NewRequest("GET", "/new/https://www.google.com", nil)
	if err != nil {
		fmt.Println(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: Got %v but want %v",
			rr.Body.String(), expected)
	}
}
