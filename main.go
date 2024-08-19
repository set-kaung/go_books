package main

import (
	"embed"
	"flag"
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

	port := flag.String("port", "6543", "port number to run the server on")

	flag.Parse()
	entries, err := createLocalSources()
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	app, err := GetApp(entries["cache.json"], *port)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	dir, err := fs.Sub(ui, "ui/dist")
	if err != nil {
		log.Fatalln(err)
	}

	nfs := neuteredFileSystem{fileSys: http.FS(dir)}
	mux.HandleFunc("/cache", handlePreflight(app.CacheFileHandler))
	mux.HandleFunc("/files", handlePreflight(app.FileHandler))
	mux.HandleFunc("GET /auth", app.Authorise)
	mux.HandleFunc("GET /callback", app.Callback)
	mux.Handle("/", http.FileServer(nfs))

	log.Printf("running server on port :%s\n", *port)
	fmt.Printf("Follow the link to go to the interface:\nhttp://localhost:%s\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", *port), mux)
	if err != nil {
		log.Fatalln("error starting server:", err)
	}

}
