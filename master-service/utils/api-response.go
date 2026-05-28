package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConfig struct
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}
type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	IsError bool        `json:"is_error"`
}

// Detect base64 image string
func isBase64Image(s string) bool {
	return strings.HasPrefix(s, "data:image/")
}

// Recursively search for base64 image in any JSON structure
func extractBase64Recursive(v interface{}) (string, bool) {

	switch val := v.(type) {

	case string:
		if isBase64Image(val) {
			return val, true
		}

	case []interface{}:
		for _, item := range val {
			if img, ok := extractBase64Recursive(item); ok {
				return img, true
			}
		}

	case map[string]interface{}:
		for _, item := range val {
			if img, ok := extractBase64Recursive(item); ok {
				return img, true
			}
		}
	}

	return "", false
}

func SendRequestWithResponse(client *http.Client, url string, payload interface{}) (string, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	if apiResp.IsError || apiResp.Status != 200 {
		return "", fmt.Errorf("%s (status %d)", apiResp.Message, apiResp.Status)
	}

	if strings.Contains(apiResp.Message, "No face found") {
		return "", fmt.Errorf("error: %s", apiResp.Message)
	}

	// UNIVERSAL BASE64 EXTRACTOR
	dataBytes, _ := json.Marshal(apiResp.Data)

	var generic interface{}
	if err := json.Unmarshal(dataBytes, &generic); err != nil {
		return "", fmt.Errorf("cannot parse response data: %w", err)
	}

	if img, found := extractBase64Recursive(generic); found {
		return img, nil
	}

	return "", fmt.Errorf("failed to extract base64 image from response")
}

func SendRequest(client *http.Client, url string, payload interface{}) error {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.IsError || apiResp.Status != 200 {
		return fmt.Errorf("%s (status: %d)", apiResp.Message, apiResp.Status)
	}

	return nil
}

// InitMinioClient menggunakan AccessKey & SecretKey
func InitMinioClient(cfg MinioConfig) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed init minio client: %w", err)
	}
	return client, nil
}

// Upload file dari io.Reader ke MinIO
func UploadToMinio(ctx context.Context, client *minio.Client, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed upload object %s/%s: %w", bucket, objectName, err)
	}
	return nil
}

// Upload excelize.File langsung ke MinIO
func UploadExcelizeToMinio(ctx context.Context, client *minio.Client, bucket, objectName string, xlsx *excelize.File) error {
	var buf bytes.Buffer
	if err := xlsx.Write(&buf); err != nil {
		return fmt.Errorf("failed to write excel file: %w", err)
	}
	return UploadToMinio(ctx, client, bucket, objectName, &buf, int64(buf.Len()),
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

// Download file dari MinIO sebagai io.ReadCloser
func GetObjectFromMinio(ctx context.Context, client *minio.Client, bucket, objectName string) (io.ReadCloser, error) {
	obj, err := client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed get object %s/%s: %w", bucket, objectName, err)
	}
	return obj, nil
}
