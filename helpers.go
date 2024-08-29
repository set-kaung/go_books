package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrWrongFileFormat = errors.New("the file must be in .json format")

func checkCacheFile(path string) (string, error) {
	if !strings.HasSuffix(path, "json") {
		return "", ErrWrongFileFormat
	}
	_, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%s not found. do you want to create the file?[y/n]:", path)
		} else {
			return path, err
		}
	}
	return path, nil
}
