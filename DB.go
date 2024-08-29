package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

func NewDatabase(dbPath string) *Database {
	db, err := initDatabese(dbPath)
	if err != nil {
		log.Fatalln("failed to initialise database in DB.NewDatabase:", err)
	}
	return &Database{db}
}

func initDatabese(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	_, err = db.ExecContext(context.Background(),
		`create table if not exists books(
				id varchar(255)  not null,
				name varchar(255)  default null,
				link varchar(255)  default null,
				extension varchar(255)  default null,
				primary key (id)    
				)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) WriteData(books []BookFile) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// Prepare the statement outside the loop
	stmt, err := tx.PrepareContext(context.Background(),
		`INSERT INTO books (id, name, link, extension) VALUES (?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement for each book
	for _, b := range books {
		_, err = stmt.ExecContext(context.Background(), b.ID, b.Name, b.Link, b.Extension)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) ReadData() ([]byte, error) {
	books, err := db.getAllBooks()
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string][]BookFile{"files": books})
}

func (db *Database) getAllBooks() ([]BookFile, error) {
	books := []BookFile{}
	result, err := db.QueryContext(context.Background(), `SELECT * FROM books`)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		b := BookFile{}
		err = result.Scan(&b.ID, &b.Name, &b.Link, &b.Extension)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}
