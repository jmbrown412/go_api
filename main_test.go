package main_test

import (
	"bytes"
	"encoding/json"
	"go_api"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a main.App

func TestEmptyDrafts(t *testing.T) {
	a.Initialize()
	clearTable()

	req, _ := http.NewRequest("GET", "/api/drafts", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateDocument(t *testing.T) {
	a.Initialize()
	//clearTable()

	var jsonStr = []byte(`{"name":"test document", "text": "test draft text"}`)
	req, _ := http.NewRequest("POST", "/api/drafts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test document" {
		t.Errorf("Expected document name to be 'test document'. Got '%v'", m["name"])
	}

	if m["text"] != "test draft text" {
		t.Errorf("Expected document text to be 'test draft text'. Got '%v'", m["text"])
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM documents")
	a.DB.Exec("ALTER SEQUENCE documents_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
