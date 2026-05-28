package uploader

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const maxImageSize = 4 << 20

type Minio struct {
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	Address       string
	PublicAddress string
	Port          int
	Bucket        string
	FileLocation  string
	Region        string
}

func (m *Minio) SetRegion(data string) {
	m.Region = data
}

func (m *Minio) SetAccessKey(data string) {
	m.AccessKey = data
}

func (m *Minio) SetSecretKey(data string) {
	m.SecretKey = data
}

func (m *Minio) SetBucket(data string) {
	m.Bucket = data
}

func (m *Minio) SetFileLocation(data string) {
	m.FileLocation = data
}

func (m Minio) GetBucket() string {
	return m.Bucket
}

func (m Minio) GetFileLocation() string {
	return m.FileLocation
}

func NewMinio(accessKey, secretKey, address, publicAddress, bucket, region string, useSSL bool, port int) *Minio {
	return &Minio{
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Address:       address,
		Bucket:        bucket,
		Region:        region,
		UseSSL:        useSSL,
		Port:          port,
		PublicAddress: publicAddress,
	}
}

func (m *Minio) newClient() (*minio.Client, error) {
	miniClient, err := minio.New(
		fmt.Sprintf("%s:%d", m.Address, m.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(m.AccessKey, m.SecretKey, ""),
			Secure: m.UseSSL,
			Region: m.Region,
		},
	)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if ok, _ := miniClient.BucketExists(ctx, m.Bucket); !ok {
		if err := miniClient.MakeBucket(ctx, m.Bucket, minio.MakeBucketOptions{
			Region: m.Region,
		}); err != nil {
			return nil, err
		}
	}

	return miniClient, nil
}

func (m *Minio) UploadBase64(image []byte, fileLocation string) (string, string, error) {
	client, err := m.newClient()
	if err != nil {
		return "", "", err
	}

	// convert byte slice to io.Reader
	reader := bytes.NewReader(image)

	// read only 4 byte from our io.Reader
	buf := make([]byte, maxImageSize)
	n, err := reader.Read(buf)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	contentType := ""
	if counter := strings.Split(http.DetectContentType(image), "/"); len(counter) > 1 {
		contentType = counter[1]
	}

	fileName := m.randomString() + "." + contentType
	fileLocationPath := fmt.Sprintf("spectogram/%v", fileName)
	if fileLocation != "" {
		fileLocationPath = fmt.Sprintf("%v/%v", fileLocation, fileName)
	}
	_, err = client.PutObject(
		ctx,
		m.Bucket,
		fileLocationPath,
		strings.NewReader(string(image[:])),
		int64(n),
		minio.PutObjectOptions{
			ContentType: http.DetectContentType(image),
		},
	)

	if err != nil {
		return "", "", err
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s.%s\"", fileName, http.DetectContentType(image)))

	// presignedUrl, err := client.PresignedGetObject(
	// 	ctx,
	// 	m.Bucket,
	// 	fileName,
	// 	24*time.Hour,
	// 	reqParams,
	// )
	// if err != nil {
	// 	return "", "", err
	// }
	// publicUrl := fmt.Sprintf("%s://%s%s", presignedUrl.Scheme, presignedUrl.Host, presignedUrl.Path)
	return fmt.Sprintf("%v/%v/%v", m.PublicAddress, m.Bucket, fileLocationPath), http.DetectContentType(image), nil
	// return publicUrl, http.DetectContentType(image), nil
}

func (m Minio) randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func (m *Minio) Upload(file *multipart.FileHeader) (fileDestination, fileType string, err error) {
	client, err := m.newClient()
	if err != nil {
		return "", "", err
	}

	read, errOpen := file.Open()
	defer read.Close()

	if errOpen != nil {
		return fileDestination, fileType, errOpen
	}

	objectName := file.Filename
	fileBuffer := read
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size
	ctx := context.Background()

	_, err = client.PutObject(ctx, m.Bucket, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return "", "", err
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s.%s\"", objectName, contentType))

	presignedUrl, err := client.PresignedGetObject(
		ctx,
		m.Bucket,
		objectName,
		24*time.Hour,
		reqParams,
	)
	if err != nil {
		return "", "", err
	}
	publicUrl := fmt.Sprintf("%s://%s%s", presignedUrl.Scheme, presignedUrl.Host, presignedUrl.Path)
	return publicUrl, contentType, nil
}

func (m Minio) Delete(bucket, obj string) error {
	client, err := m.newClient()
	if err != nil {
		return err
	}

	if err := client.RemoveObject(context.Background(), bucket, obj, minio.RemoveObjectOptions{
		ForceDelete: true,
	}); err != nil {
		return err
	}
	return nil

}

// */
