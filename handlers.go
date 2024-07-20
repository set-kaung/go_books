package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

type BookFile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Extension string `json:"ext"`
}

var ErrNoFiles = errors.New("no files found")

func (app *Application) CacheFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("cache update request processing...")
	_, err := os.Open(app.CacheFilePath)
	if err != nil {
		err = ServerErrorResp(w, "cache file not found", err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	log.Println("api request started")
	res, err := app.Service.Files.List().Q("trashed=false").Fields("nextPageToken, files(id, name, webContentLink)").Do()
	if err != nil {
		log.Println(err)
		err = ServerErrorResp(w, "api request failed from server", err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if len(res.Files) == 0 {
		err = NotFoundResp(w, "no files found", err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	pageToken := ""
	files := []BookFile{}
	for {
		res, err := app.Service.Files.List().Q("trashed=false").PageToken(pageToken).Fields("nextPageToken, files(id, name, webContentLink)").Do()
		if err != nil {
			err = NotFoundResp(w, "no files found", err)
			if err != nil {
				log.Println(err)
			}
			return
		}

		for _, f := range res.Files {
			if strings.Contains(f.Name, ".pdf") || strings.Contains(f.Name, ".epub") {
				file := BookFile{ID: f.Id, Name: f.Name, Link: f.WebContentLink}
				files = append(files, file)
			}
		}

		if res.NextPageToken == "" {
			break
		}

		pageToken = res.NextPageToken
	}
	file, err := os.OpenFile(app.CacheFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		err = ServerErrorResp(w, "faile to open cache file", err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	encoder := json.NewEncoder(file)
	log.Println("api request ended")
	err = encoder.Encode(map[string][]BookFile{"files": files})
	if err != nil {
		err = ServerErrorResp(w, "failed to encode json", err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	log.Println("caching finished")
	file.Close()
}

func (app *Application) FileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.URL.Path != "/" {
		err := BadRequestResp(w, "method not allowed or wrong path", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	fileData, err := os.ReadFile(app.CacheFilePath)
	if err != nil {
		err = ServerErrorResp(w, "can't open cache file", err)
		if err != nil {
			log.Println("File Handler", err)
		}
		return
	}
	w.Header().Set("content-type", "application/json")
	written, err := w.Write(fileData)
	if err != nil {
		err = ServerErrorResp(w, "error writing json", err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	log.Printf("written %d as response\n", written)
}
