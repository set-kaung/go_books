package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Application struct {
	Client     *http.Client
	Service    *drive.Service
	LocalCache JSONReadWriter
	Database   JSONReadWriter
}
type BookFile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Extension string `json:"ext"`
}

type JSONReadWriter interface {
	WriteData([]BookFile) error
	ReadData() ([]byte, error)
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
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

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
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

func GetApp(cfilePath string) (Application, error) {
	app := Application{LocalCache: &CacheFile{cfilePath}}
	app.Database = NewDatabase("./books.db")
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return app, err
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return app, err
	}
	client := getClient(config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return app, err
	}
	app.Service = srv
	app.Client = client
	return app, nil
}
