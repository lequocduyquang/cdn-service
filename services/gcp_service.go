package services

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"

	"cloud.google.com/go/storage"
)

var (
	// GCPService exported gcp service
	GCPService GCPServiceInterface = &gcpService{}
)

// GCPServiceInterface interface
type GCPServiceInterface interface {
	Delete(context.Context, *storage.ObjectHandle) error
	Get(context.Context, *storage.ObjectHandle) ([]byte, error)
	Upload(context.Context, *storage.ObjectHandle, multipart.File) (*url.URL, error)
}

type gcpService struct{}

// Delete function
func (s *gcpService) Delete(ctx context.Context, obj *storage.ObjectHandle) error {
	if err := obj.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// Get function
func (s *gcpService) Get(ctx context.Context, obj *storage.ObjectHandle) ([]byte, error) {
	rc, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Upload function
func (s *gcpService) Upload(ctx context.Context, obj *storage.ObjectHandle, file multipart.File) (*url.URL, error) {
	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return nil, err
	}

	if err := wc.Close(); err != nil {
		return nil, err
	}

	resultURL, err := url.Parse("/" + os.Getenv("GCP_BUCKET") + "/" + wc.Attrs().Name)
	if err != nil {
		return nil, err
	}
	return resultURL, nil
}
