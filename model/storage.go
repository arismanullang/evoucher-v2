package model

import (
	"fmt"
	"os"

	"cloud.google.com/go/storage"

	"golang.org/x/net/context"
)

var StorageBucket *storage.BucketHandle

func GcsInit() error {
	var err error
	gcsBucket := os.Getenv("GCS_BUCKET")
	StorageBucket, err = configureStorage(gcsBucket)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}
