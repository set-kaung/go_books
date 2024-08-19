package main

import (
	"fmt"
	"go_books/data"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type Application struct {
	Client     *http.Client
	Service    *drive.Service
	LocalCache data.JSONReadWriter
	Database   data.JSONReadWriter
	Port       string
	config     *oauth2.Config
}

// func getClient(config *oauth2.Config) *http.Client {

// }

func GetApp(cacheFilePath string, port string) (Application, error) {
	app := Application{LocalCache: data.NewCacheFile(cacheFilePath)}
	app.Database = data.NewDatabase("/Users/setkaung/Library/Application Support/go_books/books.db")
	app.Port = port
	b, err := os.ReadFile("credentials.json")
	if err != nil {
	}
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		panic(err)
	}
	// Set the redirect URI here
	config.RedirectURL = fmt.Sprintf("http://localhost:%s/callback", app.Port)
	app.config = config
	fmt.Println(config.RedirectURL)
	return app, nil
}
