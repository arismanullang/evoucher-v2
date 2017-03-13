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
	r.PostFunc("/create/variant", controller.CreateVariant)
	r.GetFunc("/api/get/allVariant", controller.GetAllVariants)
	r.Get("/v1/api/get/allVariant", controller.CheckTokenAuth(controller.GetAllVariants))
	r.GetFunc("/api/get/variant", controller.GetVariants)
	r.GetFunc("/api/get/variantByDate", controller.GetVariantDetailsByDate)
	r.GetFunc("/api/get/variant/:id", controller.GetVariantDetailsById)
	r.PostFunc("/api/update/variant/:id", controller.UpdateVariant)
	r.PostFunc("/update/variant/:id/broadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/update/variant/:id/tenant", controller.UpdateVariantTenant)
	r.PostFunc("/delete/variant/:id", controller.DeleteVariant)

	//transaction
	r.PostFunc("/transaction/redeem", controller.CreateTransaction)
	r.PostFunc("/transaction/redeem", controller.CreateTransaction)
	r.Post("/v1/transaction/redeem", controller.CheckTokenAuth(controller.CreateTransaction))
	r.GetFunc("/transaction/:id/", controller.GetTransactionDetails)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)

	//user
	r.GetFunc("/get/session", controller.CheckSession)
	r.PostFunc("/create/user", controller.RegisterUser)
	r.GetFunc("/get/userByRole", controller.FindUserByRole)
	r.GetFunc("/get/user", controller.GetUser)
	r.PostFunc("/login", controller.DoLogin)

	//partner
	r.GetFunc("/api/get/partner", controller.GetAllPartner)
	r.PostFunc("/create/partner", controller.AddPartner)
	r.GetFunc("/get/partner", controller.DashboardGetAllPartner)

	//account
	r.GetFunc("/get/accountId", controller.GetAccountId)

	//Voucher
	r.GetFunc("/voucher", controller.GetVoucherDetail)
	r.PostFunc("/voucher/delete", controller.DeleteVoucher)
	r.PostFunc("/voucher/pay", controller.PayVoucher)
	r.PostFunc("/voucher/generate/bulk", controller.GenerateVoucher)
	r.PostFunc("/voucher/generate/single", controller.GenerateVoucherOnDemand)

	r.Get("/v1/voucher", controller.CheckTokenAuth(controller.GetVoucherDetail))
	r.Post("/v1/voucher/delete", controller.CheckTokenAuth(controller.DeleteVoucher))
	r.Post("/v1/voucher/pay", controller.CheckTokenAuth(controller.PayVoucher))
	r.Post("/v1/voucher/generate/bulk", controller.CheckTokenAuth(controller.GenerateVoucher))
	r.Post("/v1/voucher/generate/single", controller.CheckTokenAuth(controller.GenerateVoucherOnDemand))

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
