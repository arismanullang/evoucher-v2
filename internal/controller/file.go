package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/gilkor/evoucher/internal/model"
	uuid "github.com/satori/go.uuid"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
)

func uploadFileFromForm(r *http.Request) (url string, err error) {
	f, fh, err := r.FormFile("image-url")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	err = model.GcsInit()
	if err != nil {
		return "", err
	}

	if model.StorageBucket == nil {
		return "", errors.New("storage bucket is missing")
	}

	// random filename, retaining existing extension.
	name := uuid.NewV4().String() + path.Ext(fh.Filename)

	ctx := context.Background()
	w := model.StorageBucket.Object(name).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, model.GCS_BUCKET, name), nil
}

func GetListFile(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	client, err := storage.NewClient(ctx)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	defer client.Close()

	bucket := client.Bucket(model.GCS_BUCKET)
	objs, err := bucket.Objects(ctx, nil).Next()
	if err != nil {
		w.Write([]byte("error"))
		return
	}

	w.Write([]byte(objs.Name))

	// status = http.StatusOK
	// res = NewResponse(d)
	// render.JSON(w, res, status)
}
