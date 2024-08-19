package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

var ErrEmptyDatabase = errors.New("empty database: books.db")

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
	ctx := context.Background()

	// Step 1: Fetch existing book IDs from the database
	existingBooks, err := db.GetBooks()
	if err != nil {
		return err
	}

	existingBookIDs := make(map[string]struct{})
	for _, b := range existingBooks {
		existingBookIDs[b.ID] = struct{}{}
	}

	// Step 2: Create a list of new books to insert
	newBooks := []BookFile{}
	for _, b := range books {
		if _, exists := existingBookIDs[b.ID]; !exists {
			newBooks = append(newBooks, b)
		}
	}

	// If there are no new books to insert, return early
	if len(newBooks) == 0 {
		return nil
	}
	tx, err := db.BeginTx(ctx, nil)
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

	// Execute the prepared statement for each new book
	for _, b := range newBooks {
		_, err = stmt.ExecContext(ctx, b.ID, b.Name, b.Link, b.Extension)
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
	books, err := db.GetBooks()
	if err != nil {
		return nil, err
	}
	if len(books) == 0 {
		return nil, ErrEmptyDatabase
	}
	return json.Marshal(map[string][]BookFile{"files": books})
}

func (db *Database) GetBooks() ([]BookFile, error) {
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
