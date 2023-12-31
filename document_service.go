package main

import (
	"database/sql"
	"time"
)

type DocumentService struct {
	DB *sql.DB
}

func (d *DocumentService) CreateDocument(name string) (*Document, error) {
	existingDocument, err := d.GetDocumentByName(name)
	if existingDocument != nil {
		return existingDocument, nil
	}

	res, err := d.DB.Exec(
		"INSERT INTO documents(name, createdat) VALUES($1, $2)", name, time.Now(),
	)

	if err != nil {
		return nil, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return d.GetDocumentById(int(id))
}

func (d *DocumentService) GetDocumentByName(name string) (*Document, error) {
	row := d.DB.QueryRow("SELECT id, name, createdat FROM documents WHERE name=?", name)

	document := Document{}
	if err := row.Scan(&document.ID, &document.Name, &document.CreatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	return &document, nil
}

func (d *DocumentService) GetDocumentById(id int) (*Document, error) {
	row := d.DB.QueryRow("SELECT id, name, createdat FROM documents WHERE id=?", id)

	document := Document{}
	if err := row.Scan(&document.ID, &document.Name, &document.CreatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	return &document, nil
}

func (d *DocumentService) GetDrafts() ([]Draft, error) {
	rows, err := d.DB.Query(
		"SELECT id, documentid, text, createdat FROM drafts")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	drafts := []Draft{}

	for rows.Next() {
		var draft Draft
		if err := rows.Scan(&draft.ID, &draft.DocumentID, &draft.Text, &draft.CreatedAt); err != nil {
			return nil, err
		}
		drafts = append(drafts, draft)
	}

	return drafts, nil
}

func (d *DocumentService) GetDocuments() ([]DocumentWithRecentDraft, error) {
	docsWithDrafts := []DocumentWithRecentDraft{}
	rows, err := d.DB.Query(
		"SELECT id, name, createdat FROM documents")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := []Document{}
	for rows.Next() {
		var document Document
		if err := rows.Scan(&document.ID, &document.Name, &document.CreatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, document)
	}

	for i := 0; i < len(docs); i++ {
		row := d.DB.QueryRow("SELECT id, documentid, text, createdat FROM drafts WHERE documentid=? ORDER BY createdat desc LIMIT 1", docs[i].ID)

		draft := Draft{}
		if err := row.Scan(&draft.ID, &draft.DocumentID, &draft.Text, &draft.CreatedAt); err == sql.ErrNoRows {
			return nil, err
		}
		docWithDraft := DocumentWithRecentDraft{Name: docs[i].Name, Text: draft.Text}
		docsWithDrafts = append(docsWithDrafts, docWithDraft)
	}

	return docsWithDrafts, nil
}

func (d *DocumentService) CreateDraft(req CreateDocumentDraftRequest) (*Draft, error) {
	doc, err := d.CreateDocument(req.Name)
	if err != nil {
		return nil, err
	}

	res, err := d.DB.Exec(
		"INSERT INTO drafts(documentid, text, createdat) VALUES($1, $2, $3)", doc.ID, req.Text, time.Now(),
	)

	if err != nil {
		return nil, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return d.GetDraftById(int(id))
}

func (d *DocumentService) GetDraftById(id int) (*Draft, error) {
	row := d.DB.QueryRow("SELECT id, documentid, text, createdat FROM drafts WHERE id=?", id)

	draft := Draft{}
	if err := row.Scan(&draft.ID, &draft.DocumentID, &draft.Text, &draft.CreatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	return &draft, nil
}

func (d *DocumentService) CreateDraftComment(req CreateDraftCommentRequest) (*Comment, error) {
	// Lookup the draft
	draft, err := d.GetDraftById(req.ID)
	if draft == nil || err != nil {
		return nil, err
	}

	userId := 1 // TODO Implement a Users table

	res, err := d.DB.Exec(
		"INSERT INTO comments(userid, text, draftid, createdat) VALUES($1, $2, $3, $4)", userId, req.Text, draft.ID, time.Now(),
	)

	if err != nil {
		return nil, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return d.GetDraftCommentById(int(id))
}

func (d *DocumentService) GetDraftCommentById(id int) (*Comment, error) {
	row := d.DB.QueryRow("SELECT id, userid, text, draftid, createdat FROM comments WHERE id=?", id)

	comment := Comment{}
	if err := row.Scan(&comment.ID, &comment.UserID, &comment.Text, &comment.DraftID, &comment.CreatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	return &comment, nil
}

func (d *DocumentService) SearchDrafts(text string) ([]Draft, error) {
	query := "%" + text + "%"
	rows, err := d.DB.Query(
		"SELECT id, documentid, text, createdat FROM drafts WHERE text LIKE ?", query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	drafts := []Draft{}

	for rows.Next() {
		var draft Draft
		if err := rows.Scan(&draft.ID, &draft.DocumentID, &draft.Text, &draft.CreatedAt); err != nil {
			return nil, err
		}
		drafts = append(drafts, draft)
	}

	return drafts, nil
}
