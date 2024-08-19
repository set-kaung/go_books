package data

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

var ErrEmtpyCache = errors.New("empty cache file: cache.json")

type CacheFile struct {
	filePath string
}

func NewCacheFile(path string) *CacheFile {
	return &CacheFile{filePath: path}
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
	file, err := os.OpenFile(c.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file is closed after reading
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrEmtpyCache
	}
	return data, nil
}

func (c *CacheFile) GetBooks() ([]BookFile, error) {
	file, err := os.OpenFile(c.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	m := map[string][]BookFile{}
	err = decoder.Decode(&m)
	return m["files"], err
}
