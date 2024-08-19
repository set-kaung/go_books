package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (app *Application) Authorise(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Println("no token file.redirect")
		authURL := app.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		http.Redirect(w, r, authURL, http.StatusFound)
		return
	} else {
		app.Client = app.config.Client(r.Context(), tok)
		srv, err := drive.NewService(r.Context(), option.WithHTTPClient(app.Client))
		if err != nil {
			err = ServerErrorResp(w, "error creating new google service: ", err)
			if err != nil {
				log.Println(err)
			}
			return
		}
		app.Service = srv
	}
}

func (app *Application) Callback(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	tok, err := app.config.Exchange(context.TODO(), code)
	saveToken("./token.json", tok)
	app.Client = app.config.Client(context.Background(), tok)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	srv, err := drive.NewService(ctx, option.WithHTTPClient(app.Client))
	if err != nil {
	}
	app.Service = srv
	http.Redirect(w, r, "/", http.StatusFound)
}

func tokenFromFile(filename string) (*oauth2.Token, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handlePreflight(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}
