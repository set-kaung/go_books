package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fileSys http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fileSys.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fileSys.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

//go:embed ui/dist
var ui embed.FS

func main() {
	mux := http.NewServeMux()
	app, err := GetApp()
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	dir, err := fs.Sub(ui, "ui/dist")
	if err != nil {
		log.Fatalln(err)
	}
	nfs := neuteredFileSystem{fileSys: http.FS(dir)}
	mux.HandleFunc("/files", app.FileHandler)
	mux.HandleFunc("GET /cache", app.CacheFiles)
	mux.Handle("/", http.FileServer(nfs))

	log.Println("running server on port 6543")
	fmt.Println("Clicked the link to go to the interface:\nhttp://localhost:6543")
	err = http.ListenAndServe(":6543", mux)
	if err != nil {
		log.Fatalln("error starting server:", err)
	}

}
