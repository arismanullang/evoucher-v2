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
<<<<<<< HEAD
	r.PostFunc("/create/variant", controller.CreateVariant)
	r.GetFunc("/api/get/allVariant", controller.GetAllVariants)
	r.GetFunc("/v1/api/get/allVariant", controller.GetAllVariants)
	r.GetFunc("/api/get/variant", controller.GetVariants)
	r.GetFunc("/api/get/variantByDate", controller.GetVariantDetailsByDate)
	r.GetFunc("/api/get/variant/:id", controller.GetVariantDetailsById)
	r.PostFunc("/update/variant/:id", controller.UpdateVariant)
	r.PostFunc("/update/variant/:id/broadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/update/variant/:id/tenant", controller.UpdateVariantTenant)

	r.GetFunc("/get/allVariant", controller.DashboardGetAllVariants)
	r.GetFunc("/get/variant/:id", controller.DashboardGetVariantDetailsById)
	r.GetFunc("/get/role", controller.GetAllAccountRoles)
	r.GetFunc("/delete/variant/:id", controller.DeleteVariant)

	//transaction
	// r.PostFunc("/transaction/redeem", controller.CreateTransaction)
	r.PostFunc("/v1/transaction/redeem", controller.CreateTransaction)
=======
	r.PostFunc("/v1/create/variant", controller.CreateVariant)
	r.GetFunc("/v1/api/get/allVariant", controller.GetAllVariants)
	r.GetFunc("/v1/api/get/variant", controller.GetVariants)
	r.GetFunc("/v1/api/get/variantByDate", controller.GetVariantDetailsByDate)
	r.GetFunc("/v1/api/get/variantDetails/custom", controller.GetVariantDetailsCustom)

	r.GetFunc("/v1/api/get/variant/:id", controller.GetVariantDetailsById)
	r.PostFunc("/v1/update/variant/:id", controller.UpdateVariant)
	r.PostFunc("/v1/update/variant/:id/broadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/v1/update/variant/:id/tenant", controller.UpdateVariantTenant)
	r.GetFunc("/v1/delete/variant/:id", controller.DeleteVariant)

	//transaction
	r.PostFunc("/api/create/transaction", controller.CreateTransaction)
	//r.Post("/v1/transaction/redeem", controller.CheckTokenAuth(controller.CreateTransaction))
>>>>>>> 51aa9b843f00fdf42197190b0c9affc4be794ec6
	r.GetFunc("/transaction/:id/", controller.GetTransactionDetails)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)

	//user
	r.PostFunc("/v1/create/user", controller.RegisterUser)
	r.GetFunc("/v1/api/get/session", controller.CheckSession)
	r.GetFunc("/v1/api/get/userByRole", controller.FindUserByRole)
	r.GetFunc("/v1/api/get/user", controller.GetUser)
	r.PostFunc("/v1/api/login", controller.DoLogin)

	//partner
	r.GetFunc("/v1/api/get/partner", controller.GetAllPartner)
	r.PostFunc("/create/partner", controller.AddPartner)

	//account
	r.GetFunc("/v1/api/get/accountId", controller.GetAccountId)
	r.GetFunc("/v1/api/get/role", controller.GetAllAccountRoles)

	//voucher
	r.GetFunc("/v1/voucher", controller.GetVoucherDetail)
	// r.PostFunc("/v1/voucher/delete", controller.DeleteVoucher)
	// r.PostFunc("/v1/voucher/pay", controller.PayVoucher)
	r.PostFunc("/v1/voucher/generate/bulk", controller.GenerateVoucher)
	r.PostFunc("/v1/voucher/generate/single", controller.GenerateVoucherOnDemand)

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
	} else if page == "" || page == "index" {
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
