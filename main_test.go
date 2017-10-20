package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestCreateUrl(t *testing.T) {
	req, err := http.NewRequest("GET", "/new/https://www.google.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createURL(w, r, httprouter.Params{})
	})

	handler.ServeHTTP(res, req)

	expected := `{"original_url":"https://www.google.com", "short_url":"https://morning-retreat-24523.herokuapp.com/get/3578"}`
	if res.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: Got %v but want %v",
			res.Body.String(), expected)
	}
}
