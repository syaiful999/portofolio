package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

// TestOpenExcelFromURL tests the OpenExcelFromURL function.
func TestOpenExcelFromURL(t *testing.T) {
	// Helper function to create a dummy Excel file in memory for successful test cases.
	createDummyExcel := func() []byte {
		f := excelize.NewFile()
		// Add some content to make it a valid file.
		f.SetCellValue("Sheet1", "A1", "Hello World")
		buf, _ := f.WriteToBuffer()
		return buf.Bytes()
	}

	dummyExcelBytes := createDummyExcel()

	testCases := []struct {
		name        string
		handler     http.HandlerFunc
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success - Valid Excel file",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(dummyExcelBytes)
			},
			expectError: false,
		},
		{
			name: "Error - Server returns 404 Not Found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
			errorMsg:    "unexpected status code: 404",
		},
		{
			name: "Error - Server returns non-Excel content",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "this is not an excel file")
			},
			expectError: true,
			errorMsg:    "failed to open excel",
		},
		{
			name: "Error - Empty response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectError: true,
			errorMsg:    "failed to open excel",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			file, err := OpenExcelFromURL(server.URL)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, file)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, file)
				// You can add more assertions here, like checking the sheet name or cell value.
				val, _ := file.GetCellValue("Sheet1", "A1")
				assert.Equal(t, "Hello World", val)
			}
		})
	}
}
