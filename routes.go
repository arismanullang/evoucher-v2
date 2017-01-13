package main

import (
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/evoucher/voucher/internal/controller"
)

func setRoutes() http.Handler {
	r := bone.New()
	r.GetFunc("/ping", ping)

	//variant
	r.PostFunc("/variant/createVariant", controller.CreateVariant)
	r.GetFunc("/variant/", controller.GetVariantDetails)
	r.PostFunc("/variant/:id/update", controller.UpdateVariant)
	r.PostFunc("/variant/:id/delete", controller.DeleteVariant)

	//transaction
	r.PostFunc("/transaction/createTransaction", controller.CreateTransaction)

	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}
