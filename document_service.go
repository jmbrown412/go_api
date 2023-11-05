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

func (d *DocumentService) CreateDraft(name, text string) (*Draft, error) {
	doc, err := d.CreateDocument(name)
	if err != nil {
		return nil, err
	}

	res, err := d.DB.Exec(
		"INSERT INTO drafts(documentid, text, createdat) VALUES($1, $2, $3)", doc.ID, text, time.Now(),
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
