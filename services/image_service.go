package services

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/lequocduyquang/cdn-service/db"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

const (
	errNameFile = "No file to upload"
	errSizeFile = 0
)

var (
	// ImageService exported image service
	ImageService ImageServiceInterface = &imageService{}
)

// ImageServiceInterface interface
type ImageServiceInterface interface {
	Upload(multipart.File, multipart.FileHeader) (string, int, error)
}

type imageService struct{}

func (s *imageService) Upload(file multipart.File, header multipart.FileHeader) (string, int, error) {
	bucket, err := gridfs.NewBucket(db.Client.Database(os.Getenv("MONGO_DB_NAME")))
	if err != nil {
		return errNameFile, errSizeFile, err
	}
	randomID := uuid.New()
	fileName := randomID.String() + "-" + header.Filename
	uploadStream, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		return errNameFile, errSizeFile, err
	}
	defer uploadStream.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errNameFile, errSizeFile, errors.New("Cannot read file upload")
	}
	fileSize, errWrite := uploadStream.Write(data)
	if errWrite != nil {
		return errNameFile, errSizeFile, errors.New("Cannot upload file data to db")
	}
	return fileName, fileSize, nil
}
