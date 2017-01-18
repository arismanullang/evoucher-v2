package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/evoucher/voucher/internal/controller"
)

func setRoutes() http.Handler {
	r := bone.New()
	r.GetFunc("/ping", ping)

	//variant
	r.PostFunc("/variant/createVariant", controller.CreateVariant)
	r.GetFunc("/variant/:id", controller.GetVariantDetailsByID)
	r.PostFunc("/variant/:id/search", controller.SearchVariant)
	r.PostFunc("/variant/:id/update", controller.UpdateVariant)
	r.PostFunc("/variant/:id/delete", controller.DeleteVariant)

	//transaction
	r.PostFunc("/transaction/createTransaction", controller.CreateTransaction)
	r.GetFunc("/transaction/:id", controller.GetTransactionDetails)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)

	//custom
	r.GetFunc("/view/:url", viewHandler)

	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "variant/index.html", nil)
}
