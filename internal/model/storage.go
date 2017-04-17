package model

import (
	"fmt"

	"cloud.google.com/go/storage"

	"golang.org/x/net/context"
)

var StorageBucket *storage.BucketHandle

func GcsInit() error {
	var err error
	StorageBucket, err = configureStorage(GCS_BUCKET)
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
