package main_test

import (
	"bytes"
	"encoding/json"
	"go_api"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var a main.App

func TestEmptyDrafts(t *testing.T) {
	a.Initialize()
	clearTables()

	req, _ := http.NewRequest("GET", "/api/drafts", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetAllDocumentsWithRecentDrafts(t *testing.T) {
	a.Initialize()
	clearTables()

	// Setup the DB with a couple drafts with some comments
	a.DB.Exec("INSERT INTO documents(name, createdat) VALUES($1, $2)", "Test document 1", time.Now())
	a.DB.Exec("INSERT INTO documents(name, createdat) VALUES($1, $2)", "Test document 2", time.Now())

	a.DB.Exec("INSERT INTO drafts(documentid, text, createdat) VALUES($1, $2, $3)", 1, "Doc1 comment1", time.Now())
	a.DB.Exec("INSERT INTO drafts(documentid, text, createdat) VALUES($1, $2, $3)", 1, "Doc1 comment2", time.Now())

	a.DB.Exec("INSERT INTO drafts(documentid, text, createdat) VALUES($1, $2, $3)", 2, "Doc2 comment1", time.Now())
	a.DB.Exec("INSERT INTO drafts(documentid, text, createdat) VALUES($1, $2, $3)", 2, "Doc2 comment2", time.Now())

	req, _ := http.NewRequest("GET", "/api/documents", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array with 2 documents. Got %s", body)
	}

	// TODO - Ran out of steam. Make this test more explicit checks
}

func TestCreateDraft(t *testing.T) {
	a.Initialize()
	clearTables()

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
	clearTables()

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

func TestCreateDraftComment(t *testing.T) {
	a.Initialize()
	clearTables()

	// First, create a draft
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

	// Now let's create a comment on the draft
	commentRequest := main.CreateDraftCommentRequest{ID: 1, Text: "test draft comment text"}
	commentRequestBytes, err := json.Marshal(commentRequest)
	if err != nil {
		t.Error(err)
	}
	req2, _ := http.NewRequest("POST", "/api/drafts/1/comments", bytes.NewBuffer(commentRequestBytes))
	req2.Header.Set("Content-Type", "application/json")

	response2 := executeRequest(req2)
	checkResponseCode(t, http.StatusCreated, response2.Code)

	var m2 map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &m2)

	if m2["text"] != commentRequest.Text {
		t.Errorf("Expected document text to be 'test draft comment text'. Got '%v'", m["text"])
	}
}

func TestSearchDrafts(t *testing.T) {
	a.Initialize()
	clearTables()

	// First, create a draft
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

	// Now let's create a comment on the draft
	commentRequest := main.CreateDraftCommentRequest{ID: 1, Text: "test draft comment text"}
	commentRequestBytes, err := json.Marshal(commentRequest)
	if err != nil {
		t.Error(err)
	}
	req2, _ := http.NewRequest("GET", "/api/drafts/comments?text=draft", bytes.NewBuffer(commentRequestBytes))
	req2.Header.Set("Content-Type", "application/json")

	response2 := executeRequest(req2)
	checkResponseCode(t, http.StatusOK, response2.Code)

	var m2 map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &m2)

	if body := response2.Body.String(); body == "[]" {
		t.Errorf("Expected an array with one draft. Got %s", body)
	}
}

func clearTables() {
	a.DB.Exec("DELETE FROM documents")
	a.DB.Exec("DELETE FROM drafts")
	a.DB.Exec("DELETE FROM comments")
	a.DB.Exec("ALTER SEQUENCE drafts_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE documents_id_seq RESTART WITH 1")
	a.DB.Exec("ALTER SEQUENCE comments_id_seq RESTART WITH 1")

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
