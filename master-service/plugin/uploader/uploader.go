package uploader

import (
	"context"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

var IUpload IUploader

func NewUploader() IUploader {
	// switching driver
	switch driver := os.Getenv("BUCKET_DRIVER"); {
	case strings.ToLower(driver) == "minio":
		IUpload = InitMinio()
	}
	// initiate
	return IUpload
}

type IUploader interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (url, fileType string, err error)
	UploadBase64(ctx context.Context, image []byte, fileLocation string) (url, fileType string, err error)
	Delete(ctx context.Context, bucket, obj string) error
}

type MinioUploader struct {
	Driver *Minio
}

func InitMinio() IUploader {
	port, _ := strconv.Atoi(os.Getenv("MINIO_PORT"))
	init := &MinioUploader{
		Driver: &Minio{
			AccessKey:     os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey:     os.Getenv("MINIO_SECRET_KEY"),
			Region:        os.Getenv("MINIO_REGION"),
			Bucket:        os.Getenv("MINIO_BUCKET"),
			FileLocation:  os.Getenv("BUCKET_CONTAINER"),
			Address:       os.Getenv("MINIO_ADDRESS"),
			PublicAddress: os.Getenv("MINIO_PUBLIC_ADDRESS"),
			Port:          port,
		},
	}
	return init
}

func (r MinioUploader) Upload(ctx context.Context, file *multipart.FileHeader) (string, string, error) {
	return r.Driver.Upload(file)
}

func (r MinioUploader) UploadBase64(ctx context.Context, image []byte, fileLocation string) (string, string, error) {
	return r.Driver.UploadBase64(image, fileLocation)
}

func (r MinioUploader) Delete(ctx context.Context, bucket, obj string) error {
	return r.Driver.Delete(bucket, obj)
}

func (r *MinioUploader) Setup() {
	r.Driver.SetAccessKey(os.Getenv("MINIO_ACCESS_KEY"))
	r.Driver.SetSecretKey(os.Getenv("MINIO_SECRET_KEY"))
	r.Driver.SetRegion(os.Getenv("MINIO_REGION"))
	r.Driver.SetBucket(os.Getenv("MINIO_BUCKET"))
}
