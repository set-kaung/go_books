package main

import "testing"

func TestGetAllBooks(t *testing.T) {
	db := NewDatabase("./books.db")
	_, err := db.getAllBooks()
	if err != nil {
		t.Errorf("failed to get books: %v\n", err)
	}
}
