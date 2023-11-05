package main

import (
	"time"
)

// ---------------------------------
// API Request Models
type CreateDocumentDraftRequest struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type CreateDraftCommentRequest struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// --------------------------------
// DB Entities
type Document struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdat"`
}

type Draft struct {
	ID         int       `json:"id"`
	DocumentID int       `json:"documemtid"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdat"`
}

type Comment struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userid"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdat"`
}
