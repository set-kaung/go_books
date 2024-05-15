package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func openBrowser(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}

func Run(app *Application) {
	bufReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the file name you want to search")
	str, err := bufReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	str = str[:len(str)-1]

	query := fmt.Sprintf("name contains '%s'", str)

	r, err := app.Service.Files.List().Q(query).Fields("nextPageToken, files(id, name, webContentLink)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for i, file := range r.Files {
			fmt.Printf("%d. %s (%s)\n", i+1, file.Name, file.Id)
		}
	}
	fmt.Println("The number of the file you want to download: ")
	str, err = bufReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	str = strings.TrimSuffix(str, "\n")
	index, err := strconv.ParseInt(str, 10, 64)
	if err != nil || (index < 1 || index > int64(len(r.Files))) {
		log.Fatalln("error: enter valid number:")
	}
	target := r.Files[index-1]
	openBrowser(target.WebContentLink)
}
