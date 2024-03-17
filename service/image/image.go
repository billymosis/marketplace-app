package image

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/billymosis/marketplace-app/handler/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func isValidFile(filename string) bool {
	allowedExtensions := []string{".jpg", ".jpeg"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range allowedExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func isValidFileSize(fileSize int64) bool {
	return fileSize >= 10*1024 && fileSize <= 2*1024*1024
}

type uploadResponse struct {
	ImageUrl string `json:"imageUrl"`
}

func Upload(client *s3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			http.Error(w, "Failed to retrieve file from form data", http.StatusBadRequest)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		file, handler, err := r.FormFile("file")
		if file == nil {
			render.ErrorCode(w, errors.New("empty"), 400)
			return
		}
		if err != nil {
			http.Error(w, "Failed to retrieve file from form data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if !isValidFile(handler.Filename) {
			http.Error(w, "Invalid file format. Must be *.jpg or *.jpeg", http.StatusBadRequest)
			return
		}

		if !isValidFileSize(handler.Size) {
			http.Error(w, "File size must be between 10KB and 2MB", http.StatusBadRequest)
			return
		}

		filename := uuid.New().String() + filepath.Ext(handler.Filename)
		sess := s3.New(client.Options())

		bucket := os.Getenv("S3_BUCKET_NAME")
		res, err := sess.PutObject(r.Context(),
			&s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(filename),
				ACL:    types.ObjectCannedACL(*aws.String("public-read")),
				Body:   file,
			})
		if err != nil {
			render.InternalError(w, err)
			return
		}
		logrus.Printf("S3: %+v", res.ResultMetadata)

		url := fmt.Sprintf("https://%s.%s/%s", bucket, "s3.amazonaws.com", filename)

		render.JSON(w, uploadResponse{ImageUrl: url}, 200)
	}
}
