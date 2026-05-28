package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractBase64Recursive(t *testing.T) {
	testCases := []struct {
		name          string
		input         interface{}
		expectedImage string
		expectedFound bool
	}{
		{
			name:          "Direct string match",
			input:         "data:image/png;base64,helloworld",
			expectedImage: "data:image/png;base64,helloworld",
			expectedFound: true,
		},
		{
			name:          "Nested in map",
			input:         map[string]interface{}{"key1": "value1", "image": "data:image/jpeg;base64,anotherworld"},
			expectedImage: "data:image/jpeg;base64,anotherworld",
			expectedFound: true,
		},
		{
			name:          "Nested in slice",
			input:         []interface{}{"string1", 123, "data:image/gif;base64,gifworld"},
			expectedImage: "data:image/gif;base64,gifworld",
			expectedFound: true,
		},
		{
			name:          "Deeply nested",
			input:         map[string]interface{}{"data": []interface{}{map[string]interface{}{"image_data": "data:image/bmp;base64,deepworld"}}},
			expectedImage: "data:image/bmp;base64,deepworld",
			expectedFound: true,
		},
		{
			name:          "No base64 image",
			input:         map[string]interface{}{"key1": "value1", "key2": 123},
			expectedImage: "",
			expectedFound: false,
		},
		{
			name:          "Not a data image URL",
			input:         "http://example.com/image.png",
			expectedImage: "",
			expectedFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img, found := extractBase64Recursive(tc.input)
			assert.Equal(t, tc.expectedFound, found)
			assert.Equal(t, tc.expectedImage, img)
		})
	}
}

func TestSendRequestWithResponse(t *testing.T) {
	t.Run("Successful response with base64 image", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			response := APIResponse{
				Data:    map[string]interface{}{"image": "data:image/png;base64,test"},
				Message: "Success",
				Status:  http.StatusOK,
				IsError: false,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		img, err := SendRequestWithResponse(server.Client(), server.URL, nil)
		assert.NoError(t, err)
		assert.Equal(t, "data:image/png;base64,test", img)
	})

	t.Run("API returns an error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			response := APIResponse{
				Message: "Bad Request",
				Status:  http.StatusBadRequest,
				IsError: true,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		_, err := SendRequestWithResponse(server.Client(), server.URL, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Bad Request (status 400)")
	})

	t.Run("No face found error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			response := APIResponse{
				Message: "No face found in image",
				Status:  http.StatusOK,
				IsError: false,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		_, err := SendRequestWithResponse(server.Client(), server.URL, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error: No face found in image")
	})

	t.Run("No base64 image in response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			response := APIResponse{
				Data:    map[string]interface{}{"message": "just a regular message"},
				Status:  http.StatusOK,
				IsError: false,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		_, err := SendRequestWithResponse(server.Client(), server.URL, nil)
		assert.Error(t, err)
		assert.Equal(t, "failed to extract base64 image from response", err.Error())
	})
}

func TestSendRequest(t *testing.T) {
	t.Run("Successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			response := APIResponse{
				Status:  http.StatusOK,
				IsError: false,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		err := SendRequest(server.Client(), server.URL, nil)
		assert.NoError(t, err)
	})

	t.Run("API returns an error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			response := APIResponse{
				Message: "Forbidden",
				Status:  http.StatusForbidden,
				IsError: true,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		err := SendRequest(server.Client(), server.URL, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Forbidden (status: 403)")
	})
}

// Note: InitMinioClient, UploadToMinio, UploadExcelizeToMinio, and GetObjectFromMinio are not unit-tested
// as they directly interact with an external service (Minio). These are better suited for integration tests.
