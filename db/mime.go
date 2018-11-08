package db

import (
	"net/http"
	"os"
)

func mimeType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buff := make([]byte, 500)
	_, err = file.Read(buff)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buff), nil
}
