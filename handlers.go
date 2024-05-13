package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

var ErrNoFiles = errors.New("no files found")

func (app *Application) CacheFiles(w http.ResponseWriter, r *http.Request) {
	log.Println("cache update request processing...")
	file, err := os.OpenFile("cache.json", os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "can't find cache file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	log.Println("api request started")
	res, err := app.Service.Files.List().Q("trashed=false").Fields("nextPageToken, files(id, name, webContentLink)").Do()
	if err != nil {
		http.Error(w, "no response from api", http.StatusInternalServerError)
		return
	}
	log.Println("api request ended")

	encoder := json.NewEncoder(file)
	if len(res.Files) == 0 {
		http.Error(w, "empty file list from api", http.StatusInternalServerError)
		return
	}
	pageToken := ""
	files := []File{}
	for {
		res, err := app.Service.Files.List().Q("trashed=false").PageToken(pageToken).Fields("nextPageToken, files(id, name, webContentLink)").Do()
		if err != nil {
			http.Error(w, "no response from api", http.StatusInternalServerError)
			return
		}

		for _, f := range res.Files {
			if strings.Contains(f.Name, ".pdf") || strings.Contains(f.Name, ".epub") {
				file := File{ID: f.Id, Name: f.Name, Link: f.WebContentLink}
				files = append(files, file)
			}
		}

		if res.NextPageToken == "" {
			break
		}

		pageToken = res.NextPageToken
	}
	err = encoder.Encode(map[string][]File{"files": files})
	if err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		return
	}
	log.Println("caching finished")
}

func (app *Application) FileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	file, err := os.ReadFile("cache.json")
	if err != nil {
		http.Error(w, "Unable to open JSON file", http.StatusInternalServerError)
		log.Printf("Error opening JSON file: %v", err)
		return
	}
	w.Header().Set("content-type", "application/json")
	written, err := w.Write(file)
	if err != nil {
		http.Error(w, "error writing json data", http.StatusInternalServerError)
		return
	}
	log.Printf("written %d as response\n", written)
}
