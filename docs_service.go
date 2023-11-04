package main

import (
	"database/sql"
	"time"
)

type DocumentService struct {
	DB *sql.DB
}

func (d *DocumentService) GetDrafts() ([]Document, error) {
	rows, err := d.DB.Query(
		"SELECT id, name, text, createdat FROM documents")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := []Document{}

	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.Name, &doc.Text, &doc.CreatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (d *DocumentService) CreateDocument(name, text string) (*Document, error) {
	res, err := d.DB.Exec(
		"INSERT INTO documents(name, text, createdat) VALUES($1, $2, $3)", name, text, time.Now(),
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

func (d *DocumentService) GetDocumentById(id int) (*Document, error) {
	row := d.DB.QueryRow("SELECT id, name, text, createdat FROM documents WHERE id=?", id)

	doc := Document{}
	if err := row.Scan(&doc.ID, &doc.Name, &doc.Text, doc.CreatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	return &doc, nil
}
