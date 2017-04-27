package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"

	"cloud.google.com/go/storage"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)

	imgURL, err := UploadFileFromForm(r)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), "Upload File")
		fmt.Println(err)
	}

	fmt.Println(imgURL)
	res = NewResponse(imgURL)
	render.JSON(w, res, status)
}

func UploadFileFromForm(r *http.Request) (url string, err error) {
	r.ParseMultipartForm(32 << 20)
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

	ext := path.Ext(fh.Filename)
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return "We do not allow files of type " + ext + ". We only allow jpg, jpeg, png extensions.", nil
	}

	// random filename, retaining existing extension.
	name := randStr(32, "Alphanumeric") + ext

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

	return fmt.Sprintf(model.PublicURL, model.GCS_BUCKET, name), nil
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	objname := r.FormValue("obj")
	if !deleteFile(w, r, objname) {
		return
	}

	render.JSON(w, nil, http.StatusOK)
}

func deleteFile(w http.ResponseWriter, r *http.Request, objname string) bool {
	res := NewResponse(nil)
	status := http.StatusOK

	err := model.GcsInit()
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, "("+err.Error()+")", "file")
		return false
	}

	if model.StorageBucket == nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, "storage bucket is missing "+"("+err.Error()+")", "file")
		return false
	}

	ctx := context.Background()
	o := model.StorageBucket.Object(objname)
	if err := o.Delete(ctx); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", "file")
		return false
	}

	render.JSON(w, nil, status)
	return true
}
