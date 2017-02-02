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

	//variant ui
	r.GetFunc("/variant/:page", viewVariant)

	//variant
	r.PostFunc("/variant/getAllVariant", controller.GetAllVariant)
	r.PostFunc("/variant/getVariant", controller.GetVariantDetails)
	r.PostFunc("/variant/getVariantByUser", controller.GetVariantDetailsByUser)
	r.PostFunc("/variant/getVariantByDate", controller.GetVariantDetailsByDate)
	r.PostFunc("/variant/createVariant", controller.CreateVariant)
	r.PostFunc("/variant/:id", controller.GetVariantDetailsByID)
	r.PostFunc("/variant/:id/update", controller.UpdateVariant)
	r.PostFunc("/variant/:id/updateBroadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/variant/:id/updateTenant", controller.UpdateVariantTenant)
	r.PostFunc("/variant/:id/delete", controller.DeleteVariant)

	//transaction
	r.PostFunc("/transaction/createTransaction", controller.CreateTransaction)
	r.GetFunc("/transaction/:id", controller.GetTransactionDetails)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)

	//user
	r.PostFunc("/user/getUserByRole", controller.GetUserByRole)
	r.GetFunc("/login", viewHandlers)
	r.GetFunc("/:id/", controller.GetToken)

	//custom
	r.GetFunc("/view/", viewHandler)
	r.GetFunc("/viewNoLayout/", viewHandlers)

	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "view/index.html", nil)
}

func viewVariant(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	switch page {
	case "create":
		render.FileInLayout(w, "layout.html", "variant/create.html", nil)
	case "search":
		render.FileInLayout(w, "layout.html", "variant/check.html", nil)
	case "update":
		render.FileInLayout(w, "layout.html", "variant/update.html", nil)
	default:
		render.FileInLayout(w, "layout.html", "variant/index.html", nil)
	}
}

func viewHandlers(w http.ResponseWriter, r *http.Request) {
	//url := "http://juno-staging.elys.id/v1/signin?redirect_url=http://127.0.0.1:8080/ping"
	url := "http://juno-staging.elys.id/v1/signin?redirect_url=http://juno-staging.elys.id/v1/signin"
	http.Redirect(w, r, url, http.StatusFound)
}
