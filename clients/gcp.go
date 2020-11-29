package clients

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// InitiateGCPClient Init instance
func InitiateGCPClient(ctx context.Context) (*storage.Client, error) {
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("lotus-fitness.json"))
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}
