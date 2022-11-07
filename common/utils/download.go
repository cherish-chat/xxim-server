package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, saveTo string) error {
	// Create the file
	out, err := os.Create(saveTo)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
