package main

import (
	"encoding/json"
	"io"
	"os"
)

type CacheFile struct {
	filePath string
}

func (c *CacheFile) WriteData(files []BookFile) error {
	_, err := os.Open(c.filePath)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(c.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(map[string][]BookFile{"files": files})
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheFile) ReadData() ([]byte, error) {
	file, err := os.OpenFile(c.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}
