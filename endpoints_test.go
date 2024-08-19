package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileHandlersOk(t *testing.T) {
	req, err := http.NewRequest("GET", "/files", nil)
	if err != nil {
		panic(err)
	}

	responseRecorder := httptest.NewRecorder()
	app, err := GetApp("./cache.json", "6543")
	if err != nil {
		panic(err)
	}
	handler := http.HandlerFunc(app.FileHandler)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler response was not ok: got %d\n", status)
	}
}

func TestFileHandlerFileError(t *testing.T) {
	req, err := http.NewRequest("GET", "/files", nil)
	if err != nil {
		panic(err)
	}

	responseRecorder := httptest.NewRecorder()
	app, err := GetApp("./deada.json", ":6868")
	if err != nil {
		panic(err)
	}
	handler := http.HandlerFunc(app.FileHandler)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler response was not ok: got %d\n", status)
	}
}

func TestCacheHandlerOK(t *testing.T) {
	req, err := http.NewRequest("GET", "/cache", nil)
	if err != nil {
		panic(err)
	}

	responseRecorder := httptest.NewRecorder()
	app, err := GetApp("./cache.json", "6543")
	if err != nil {
		panic(err)
	}
	handler := http.HandlerFunc(app.CacheFileHandler)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler response was not ok: got %d\n", status)
	}
}

func TestCacheHandlerSvrError(t *testing.T) {
	req, err := http.NewRequest("GET", "/cache", nil)
	if err != nil {
		panic(err)
	}

	responseRecorder := httptest.NewRecorder()
	app, err := GetApp("./cache.json", "6543")
	if err != nil {
		panic(err)
	}
	handler := http.HandlerFunc(app.CacheFileHandler)
	handler.ServeHTTP(responseRecorder, req)
	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler response was not ok: got %d\n", status)
	}
}
