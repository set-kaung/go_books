package main

import (
	"context"
	"errors"
	"go_books/data"
	"log"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var ErrNoFiles = errors.New("no files found")

// CacheFileHandler handles the request to update the cache with the latest files from Google Drive.
func (app *Application) CacheFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(app.Client)
	if app.Client == nil {
		log.Println("Client is nil, attempting to load token from file")
		tok, err := tokenFromFile("./token.json")
		if err != nil {
			log.Printf("Failed to load token from file: %v. Redirecting to authentication.", err)
			http.Redirect(w, r, "http://localhost:"+app.Port+"/auth", http.StatusTemporaryRedirect)
			return
		}
		app.Client = app.config.Client(context.Background(), tok)
		srv, err := drive.NewService(r.Context(), option.WithHTTPClient(app.Client))
		if err != nil {
			log.Printf("Failed to create new Google Drive service: %v", err)
			err = ServerErrorResp(w, "failed to create new Google Drive service", err)
			if err != nil {
				log.Printf("Failed to send server error response: %v", err)
			}
			return
		}
		app.Service = srv
	}
	log.Println("Starting API request to fetch files from Google Drive")
	res, err := app.Service.Files.List().Q("trashed=false").Fields("nextPageToken, files(id, name, webContentLink)").Do()
	if err != nil {
		log.Printf("API request failed: %v", err)
		err = ServerErrorResp(w, "API request to Google Drive failed", err)
		if err != nil {
			log.Printf("Failed to send server error response: %v", err)
		}
		return
	}

	if len(res.Files) == 0 {
		log.Println("No files found in Google Drive")
		err = NotFoundResp(w, "no files found in Google Drive", ErrNoFiles)
		if err != nil {
			log.Printf("Failed to send not found response: %v", err)
		}
		return
	}
	pageToken := ""
	files := []data.BookFile{}
	for {
		res, err = app.Service.Files.List().Q("trashed=false").PageToken(pageToken).Fields("nextPageToken, files(id, name, webContentLink)").Do()
		if err != nil {
			log.Printf("API request failed while paginating: %v", err)
			err = ServerErrorResp(w, "failed to paginate through Google Drive files", err)
			if err != nil {
				log.Printf("Failed to send server error response: %v", err)
			}
			return
		}

		for _, f := range res.Files {
			if strings.Contains(f.Name, ".pdf") || strings.Contains(f.Name, ".epub") {
				file := data.BookFile{ID: f.Id, Name: f.Name, Link: f.WebContentLink}
				files = append(files, file)
			}
		}

		if res.NextPageToken == "" {
			break
		}

		pageToken = res.NextPageToken
	}

	log.Println("API request completed, updating cache")

	err = app.LocalCache.WriteData(files)
	if err != nil {
		log.Printf("Failed to write to cache: %v", err)
		err = ServerErrorResp(w, "failed to write to cache file", err)
		if err != nil {
			log.Printf("Failed to send server error response: %v", err)
		}
		return
	}
	err = app.Database.WriteData(files)
	if err != nil {
		log.Printf("Failed to write to database: %v", err)
		err = ServerErrorResp(w, "failed to write to database", err)
		if err != nil {
			log.Printf("Failed to send server error response: %v", err)
		}
		return
	}
	log.Println("Caching finished successfully")
	http.Redirect(w, r, "/", http.StatusFound)
}

// FileHandler serves the database-stored or cached files to the client.
func (app *Application) FileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.URL.Path != "/" {
		err := BadRequestResp(w, "method not allowed or wrong path", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	fileData, err := app.Database.ReadData()
	if err != nil && !errors.Is(err, data.ErrEmptyDatabase) {
		log.Println(err)
		err = ServerErrorResp(w, "can't read from database file", err)
		if err != nil {
			log.Println("File Handler", err)
		}
		return
	}
	if len(fileData) == 0 {
		fileData, err = app.LocalCache.ReadData()
		if err != nil && !errors.Is(err, data.ErrEmtpyCache) {
			log.Println("failed to read local cache file: ", err)
			err = ServerErrorResp(w, "can't read from database file", err)
			if err != nil {
				log.Println("File Handler", err)
			}
			return
		}

		if len(fileData) == 0 {
			http.Redirect(w, r, "/cache", http.StatusTemporaryRedirect)
			return
		}
	}
	wg := &sync.WaitGroup{}
	if errors.Is(err, data.ErrEmptyDatabase) && errors.Is(err, data.ErrEmtpyCache) {
		http.Redirect(w, r, "/cache", http.StatusTemporaryRedirect)
		return
	}
	if errors.Is(err, data.ErrEmptyDatabase) && !errors.Is(err, data.ErrEmtpyCache) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			files, err := app.LocalCache.GetBooks()
			if err != nil {
				log.Println("error reading local cache in FileHandler: ", err)
				return
			}
			err = app.Database.WriteData(files)
			if err != nil {
				log.Println("error writting local cache in FileHandler: ", err)
				return
			}
		}()
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
	wg.Wait()
}
