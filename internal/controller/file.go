package controller

import (
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gilkor/evoucher/internal/model"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
)

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

func uploadFile(req *http.Request, mpf multipart.File, hdr *multipart.FileHeader) (string, error) {

	ext, err := fileFilter(req, hdr)
	if err != nil {
		return "", err
	}
	name := getSha(mpf) + `.` + ext
	mpf.Seek(0, 0)

	ctx := appengine.NewContext(req)
	return name, model.PutFile(ctx, name, mpf)
}

func fileFilter(req *http.Request, hdr *multipart.FileHeader) (string, error) {

	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]
	ctx := appengine.NewContext(req)
	log.Infof(ctx, "FILE EXTENSION: %s", ext)

	switch ext {
	case "jpg", "jpeg", "txt", "md":
		return ext, nil
	}
	return ext, fmt.Errorf("We do not allow files of type %s. We only allow jpg, jpeg, txt, md extensions.", ext)
}

func getSha(src multipart.File) string {
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}
