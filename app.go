package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Router     *mux.Router
	DocService *DocumentService
	DB         *sql.DB
}

func (a *App) Initialize() (bool, *error) {
	const file string = "documents.db"
	const create string = `
		CREATE TABLE IF NOT EXISTS documents (
			id INTEGER NOT NULL PRIMARY KEY,
			name text UNIQUE,
			createdat DATETIME NOT NULL
		);
		CREATE TABLE IF NOT EXISTS drafts (
			id INTEGER NOT NULL PRIMARY KEY,
			documentid INTEGER,
			text TEXT,
			createdat DATETIME NOT NULL,
  			FOREIGN KEY(documentid) REFERENCES documents(id)
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER NOT NULL PRIMARY KEY,
			userid INTEGER,
			text TEXT,
			draftid INTEGER,
			createdat DATETIME NOT NULL,
			FOREIGN KEY(draftid) REFERENCES document_drafts(id)
		);
		`
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return false, &err
	}
	if _, err := db.Exec(create); err != nil {
		return false, &err
	}

	a.DocService = &DocumentService{DB: db}
	a.DB = db
	a.Router = mux.NewRouter()
	a.initializeRoutes()

	return true, nil
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/drafts", a.getDraftsHandler).Methods("GET")
	a.Router.HandleFunc("/api/drafts", a.createDraftHandler).Methods("POST")
	a.Router.HandleFunc("/api/drafts/{id:[0-9]+}/comments", a.createCommentHandler).Methods("POST")
	a.Router.HandleFunc("/api/drafts/comments", a.searchDraftsHandler).Methods("GET")
}

func (a *App) getDraftsHandler(w http.ResponseWriter, r *http.Request) {
	docs, err := a.DocService.GetDrafts()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "No drafts found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, docs)
}

func (a *App) searchDraftsHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	text := v.Get("text")
	drafts, err := a.DocService.SearchDrafts(text)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "No drafts found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, drafts)
}

func (a *App) createDraftHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateDocumentDraftRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	doc, err := a.DocService.CreateDraft(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, doc)
}

func (a *App) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateDraftCommentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	draft, err := a.DocService.CreateDraftComment(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, draft)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) Run(addr string) {
	a.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		fmt.Println(tpl, err1, met, err2)
		return nil
	})
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}
