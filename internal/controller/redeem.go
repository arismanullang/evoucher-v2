package controller

import (
	"net/http"

	"github.com/ruizu/render"
)

func RedeemPage(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "redeem.html", nil)
}

func GetRedeemData(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := NewResponse(nil)
	// k := r.FormValue("key")

	res = NewResponse(nil)
	render.JSON(w, res, status)
}
