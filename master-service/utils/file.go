package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/xuri/excelize/v2"
)

func OpenExcelFromURL(url string) (*excelize.File, error) {
	// 1. GET the file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 2. Read bytes
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 3. Load excel from []byte
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to open excel: %w", err)
	}

	return f, nil
}
