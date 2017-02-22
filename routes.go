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
	r.GetFunc("/", login)

	//variant
	r.GetFunc("/get/allVariant", controller.GetAllVariant)
	r.GetFunc("/get/variant", controller.GetVariantDetails)
	r.GetFunc("/get/variantByDate", controller.GetVariantDetailsByDate)
	r.PostFunc("/create/variant", controller.CreateVariant)
	r.GetFunc("/get/variant/:id", controller.GetVariantDetailsById)
	r.GetFunc("/get/session", controller.CheckSession)
	r.PostFunc("/update/variant/:id", controller.UpdateVariant)
	r.PostFunc("/update/variant/:id/broadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/update/variant/:id/tenant", controller.UpdateVariantTenant)
	r.PostFunc("/delete/variant/:id", controller.DeleteVariant)

	//transaction
	r.PostFunc("/transaction/createTransaction", controller.CreateTransaction)
	r.GetFunc("/transaction/:id/", controller.GetTransactionDetails)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)

	//user
	r.PostFunc("/user/createUser/", controller.RegisterUser)
	r.PostFunc("/user/getUserByRole/", controller.FindUserByRole)
	r.PostFunc("/user/getUser/", controller.GetUser)
	r.PostFunc("/login", controller.DoLogin)

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

	if page == "create" {
		render.FileInLayout(w, "layout.html", "variant/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "variant/check.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "variant/update.html", nil)
	} else if page == "" || page == "index.html" {
		render.FileInLayout(w, "layout.html", "variant/index.html", nil)
	}
}

func viewUser(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "register" {
		render.FileInLayout(w, "layout.html", "user/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "user/check.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "user/update.html", nil)
	} else if page == "login" {
		render.FileInLayout(w, "layout.html", "user/login.html", nil)
	} else if page == "" || page == "index.html" {
		render.FileInLayout(w, "layout.html", "user/index.html", nil)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "user/login.html", nil)
}
