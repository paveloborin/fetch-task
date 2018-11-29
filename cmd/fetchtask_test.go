package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paveloborin/fetch-task/pkg/model"
)

func TestStatusNotFound(t *testing.T) {
	handler := router()
	req, err := http.NewRequest("GET", "/", nil)
	if nil != err {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("path '%s' expected %d http code got %d", "/", http.StatusNotFound, rr.Code)
	}
}

func TestAPI(t *testing.T) {
	handler := router()
	//POST
	var jsonStr = []byte(`{"method": "GET", "uri": "https://google.com"}`)
	req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	if nil != err {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("path '%s' expected %d http code got %d", "POST /task", http.StatusOK, rr.Code)
	}

	body, _ := ioutil.ReadAll(rr.Body)
	task := model.Task{}

	if err := json.Unmarshal(body, &task); err != nil {
		t.Errorf("invalid response")
	}

	if task.ID.String() == "" {
		t.Errorf("invalid response")
	}

	//GET by id
	req, err = http.NewRequest("GET", "/task/"+task.ID.String(), nil)
	if nil != err {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("path '%s' expected %d http code got %d", "GET /task/"+task.ID.String(), http.StatusOK, rr.Code)
	}

	//GET All
	req, err = http.NewRequest("GET", "/task", nil)
	if nil != err {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("path '%s' expected %d http code got %d", "GET /task", http.StatusOK, rr.Code)
	}

	//DELETE
	req, err = http.NewRequest("DELETE", "/task/"+task.ID.String(), nil)
	if nil != err {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("path '%s' expected %d http code got %d", "DELETE /task/"+task.ID.String(), http.StatusNoContent, rr.Code)
	}

	//GET by id non existing task
	req, err = http.NewRequest("GET", "/task/"+task.ID.String(), nil)
	if nil != err {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("path '%s' expected %d http code got %d", "GET /task/"+task.ID.String(), http.StatusNotFound, rr.Code)
	}
}
