package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	imgURL, err := UploadFileFromForm(r)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		fmt.Println(err)
		return
	}

	res.SetResponse(imgURL)
	res.JSON(w, res, http.StatusOK)
}

func UploadFileFromForm(r *http.Request) (url string, err error) {
	r.ParseMultipartForm(32 << 20)
	fmt.Println("upload file from form data")
	f, fh, err := r.FormFile("image-url")
	if err == http.ErrMissingFile {
		return "", err
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

	ext := strings.ToLower(path.Ext(fh.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		return "We do not allow files of type " + ext + " , We only allow jpg, jpeg, png extensions.", nil
	}

	// random filename, retaining existing extension. -> v2 used to add folder
	name := "v2/" + u.RandomizeString(32, "Alphanumeric") + ext

	fmt.Println("filename = ", name)
	b, err := json.Marshal(model.StorageBucket)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("storageBucket = ", string(b))
	fmt.Println("storageBucket = ", model.StorageBucket)

	ctx := context.Background()
	acls, err := model.StorageBucket.ACL().List(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	for _, rule := range acls {
		fmt.Printf("%s has role %s\n", rule.Entity, rule.Role)
	}

	fmt.Println("acl = ", acls)

	w := model.StorageBucket.Object(name).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		fmt.Println("error copy", err.Error())
		return "", err
	}
	if err := w.Close(); err != nil {
		fmt.Println("error w.close", err.Error())
		return "", err
	}

	return fmt.Sprintf(os.Getenv("GCS_PUBLIC_URL"), os.Getenv("GCS_BUCKET"), name), nil
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	objname := r.FormValue("obj")
	if !deleteFile(w, r, objname) {
		return
	}

	res.JSON(w, nil, http.StatusOK)

	// 	ctx := context.Background()

	// ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// defer cancel()
	// o := client.Bucket(bucket).Object(object)
	// if err := o.Delete(ctx); err != nil {
	//         return err
	// }
}

func deleteFile(w http.ResponseWriter, r *http.Request, objname string) bool {
	res := u.NewResponse()
	status := http.StatusOK

	err := model.GcsInit()
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		return false
	}

	if model.StorageBucket == nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		return false
	}

	ctx := context.Background()
	o := model.StorageBucket.Object(objname)
	if err := o.Delete(ctx); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		return false
	}

	res.JSON(w, nil, status)
	return true
}
