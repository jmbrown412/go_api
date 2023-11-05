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

func TestCreateDraft(t *testing.T) {
	a.Initialize()
	clearTable()

	request := main.CreateDocumentDraftRequest{Name: "test document", Text: "test draft text"}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	req, _ := http.NewRequest("POST", "/api/drafts", bytes.NewBuffer(requestBytes))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["text"] != "test draft text" {
		t.Errorf("Expected document text to be 'test draft text'. Got '%v'", m["text"])
	}
}

func TestCreateDocumentMultipleDrafts(t *testing.T) {
	a.Initialize()
	clearTable()

	request := main.CreateDocumentDraftRequest{Name: "test document", Text: "test draft text"}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}
	req, _ := http.NewRequest("POST", "/api/drafts", bytes.NewBuffer(requestBytes))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["text"] != "test draft text" {
		t.Errorf("Expected document text to be 'test draft text'. Got '%v'", m["text"])
	}

	request2 := main.CreateDocumentDraftRequest{Name: "test document", Text: "test draft text 2"}
	requestBytes2, err := json.Marshal(request2)
	if err != nil {
		t.Error(err)
	}
	req2, _ := http.NewRequest("POST", "/api/drafts", bytes.NewBuffer(requestBytes2))
	req2.Header.Set("Content-Type", "application/json")

	response2 := executeRequest(req2)
	checkResponseCode(t, http.StatusCreated, response2.Code)

	var m2 map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &m2)

	if m2["text"] != "test draft text 2" {
		t.Errorf("Expected document text to be 'test draft text 2'. Got '%v'", m2["text"])
	}

	if m2["documentid"] != m["documentid"] {
		t.Errorf("Expected documentid to be same for both drafts. Got '%v'", m2["documentid"])
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM documents")
	a.DB.Exec("DELETE FROM drafts")
	a.DB.Exec("ALTER SEQUENCE drafts_id_seq RESTART WITH 1")
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
