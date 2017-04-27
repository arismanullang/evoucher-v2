package controller

import (
	"net/http"

	"github.com/ruizu/render"
)

func UploadFormTest(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "testform.html", nil)
}
