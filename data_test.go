package main

import (
	"go_books/data"
	"testing"
)

func TestReadData(t *testing.T) {
	db := data.NewDatabase("./books.db")
	_, err := db.ReadData()
	if err != nil {
		t.Errorf("failed to get books: %v\n", err)
	}
}

func TestGetBooks(t *testing.T) {
	cache := data.NewCacheFile("./cache.json")
	files, err := cache.GetBooks()
	if err != nil || len(files) == 0 {
		t.Errorf("failed to get books: %v\n", err)
	}
}
