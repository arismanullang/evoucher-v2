package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/controller"
)

func setRoutes() http.Handler {
	r := bone.New()
	// http.ListenAndServe(":8888", nil)
	r.GetFunc("/ping", ping)

	//ui
	r.GetFunc("/variant/:page", viewVariant)
	r.GetFunc("/user/:page", viewUser)

	//variant
	r.GetFunc("/variant/getAllVariant", controller.GetAllVariant)
	r.PostFunc("/variant/getVariant", controller.GetVariantDetails)
	r.PostFunc("/variant/getVariantByUser", controller.GetVariantDetailsByUser)
	r.PostFunc("/variant/getVariantByDate", controller.GetVariantDetailsByDate)
	r.PostFunc("/variant/createVariant", controller.CreateVariant)
	r.PostFunc("/variant/getVariant/:id", controller.GetVariantDetailsById)
	r.PostFunc("/variant/:id/update", controller.UpdateVariant)
	r.PostFunc("/variant/:id/updateBroadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/variant/:id/updateTenant", controller.UpdateVariantTenant)
	r.PostFunc("/variant/:id/delete", controller.DeleteVariant)

	//transaction
	r.PostFunc("/transaction/createTransaction", controller.CreateTransaction)
	r.GetFunc("/transaction/:id/", controller.GetTransactionDetails)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)

	//user
	r.PostFunc("/user/register/", controller.RegisterUser)
	r.PostFunc("/user/getUserByRole/", controller.FindUserByRole)
	r.PostFunc("/user/getUser/", controller.GetUser)

	//Voucher
	r.GetFunc("/voucher/:id/get", controller.GetVoucherDetail)
	r.PostFunc("/voucher/redeem", controller.RedeemVoucher)
	r.PostFunc("/voucher/delete", controller.DeleteVoucher)
	r.PostFunc("/voucher/pay", controller.PayVoucher)
	r.PostFunc("/voucher/generateondemand", controller.GenerateVoucherOnDemand)
	r.PostFunc("/voucher/generate", controller.GenerateVoucher)

	//custom
	r.GetFunc("/view/", viewHandler)
	r.GetFunc("/viewNoLayout", viewNoLayoutHandler)

	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "view/index.html", nil)
}

func viewNoLayoutHandler(w http.ResponseWriter, r *http.Request) {
	render.File(w, "view/noLayout.html", nil)
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

func viewUser(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	switch page {
	case "create":
		render.FileInLayout(w, "layout.html", "user/create.html", nil)
	case "search":
		render.FileInLayout(w, "layout.html", "user/check.html", nil)
	case "update":
		render.FileInLayout(w, "layout.html", "user/update.html", nil)
	case "login":
		http.Redirect(w, r, "http://juno-staging.elys.id/v1/signin?redirect_url=http://evoucher.elys.id:8080/variant/", http.StatusFound)
	default:
		http.Redirect(w, r, "http://juno-staging.elys.id/v1/signin?redirect_url=http://evoucher.elys.id:8080/variant/", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//url := "http://juno-staging.elys.id/v1/signin?redirect_url=http://127.0.0.1:8080/ping"
	url := "http://juno-staging.elys.id/v1/signin?redirect_url=http://evoucher.elys.id:8080/variant/"
	http.Redirect(w, r, url, http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request) {
	//url := "http://juno-staging.elys.id/v1/signin?redirect_url=http://127.0.0.1:8080/ping"
	url := "http://juno-staging.elys.id/v1/register?redirect_url=http://evoucher.elys.id:8080/variant/"
	http.Redirect(w, r, url, http.StatusFound)
}
