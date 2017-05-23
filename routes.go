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
	r.GetFunc("/partner/:page", viewPartner)
	r.GetFunc("/tag/:page", viewTag)
	r.GetFunc("/voucher/:page", viewVoucher)
	r.GetFunc("/report/:page", viewReport)
	r.GetFunc("/", login)
	r.PostFunc("/v1/query", controller.CustomQuery)

	//report
	r.GetFunc("/v1/report", controller.MakeReport)
	r.GetFunc("/v1/report/variant", controller.MakeReportVariant)
	r.GetFunc("/v1/report/voucher/variant", controller.MakeCompleteReportVoucherByUser)
	r.GetFunc("/v1/report/vouchers/variant", controller.MakeReportVoucherByUser)

	//variant
	r.PostFunc("/v1/create/variant", controller.CreateVariant)
	r.GetFunc("/v1/api/get/allVariant", controller.GetAllVariants)
	r.GetFunc("/v1/api/get/variant", controller.GetVariants)
	r.GetFunc("/v1/api/get/totalVariant", controller.GetTotalVariant)
	r.GetFunc("/v1/api/get/variantByDate", controller.GetVariantDetailsByDate)
	r.GetFunc("/v1/api/get/variantDetails/custom", controller.GetVariantDetailsCustom)

	r.GetFunc("/v1/api/get/variant/:id", controller.GetVariantDetailsById)
	r.PostFunc("/v1/update/variant/:id", controller.UpdateVariant)
	r.PostFunc("/v1/update/variant/:id/broadcastUser", controller.UpdateVariantBroadcast)
	r.PostFunc("/v1/update/variant/:id/tenant", controller.UpdateVariantTenant)
	r.GetFunc("/v1/delete/variant/:id", controller.DeleteVariant)

	//transaction
	r.PostFunc("/v1/transaction/redeem", controller.MobileCreateTransaction)
	r.GetFunc("/transaction/:id/", controller.GetTransaction)
	r.PostFunc("/transaction/:id/update", controller.UpdateTransaction)
	r.PostFunc("/transaction/:id/delete", controller.DeleteTransaction)
	r.GetFunc("/v1/get/transaction", controller.GetAllTransactions)
	r.GetFunc("/v1/get/transaction/partner", controller.GetAllTransactionsByPartner)

	//user
	r.PostFunc("/v1/create/user", controller.RegisterUser)
	r.PostFunc("/v1/update/user", controller.UpdateUser)
	r.PostFunc("/v1/update/user/password", controller.ChangePassword)
	r.GetFunc("/v1/api/get/userByRole", controller.FindUserByRole)
	r.GetFunc("/v1/api/get/users", controller.GetUser)
	r.GetFunc("/v1/api/get/userDetails", controller.GetUserDetails)
	r.GetFunc("/v1/api/mail", controller.ForgotPassword)
	r.PostFunc("/v1/password", controller.UpdatePassword)
	r.PostFunc("/v1/upload/user", controller.InsertBroadcastUser)

	//partner
	r.GetFunc("/v1/get/partner", controller.GetAllPartners)
	r.GetFunc("/v1/get/partner/:id", controller.GetPartnerDetails)
	r.PostFunc("/v1/update/partner/:id", controller.UpdatePartner)
	r.GetFunc("/v1/delete/partner/:id", controller.DeletePartner)
	r.GetFunc("/v1/api/get/partner", controller.GetAllPartnersCustomParam)
	r.PostFunc("/v1/create/partner", controller.AddPartner)

	r.GetFunc("/v1/get/tag", controller.GetAllTags)
	r.PostFunc("/v1/create/tag", controller.AddTag)
	r.GetFunc("/v1/delete/tag/:id", controller.DeleteTag)
	r.PostFunc("/v1/delete/tag", controller.DeleteTagBulk)

	//account
	r.GetFunc("/v1/api/get/account", controller.GetAccount)
	r.GetFunc("/v1/api/get/accountsDetail", controller.GetAccountDetailByUser)
	r.GetFunc("/v1/api/get/accountsByUser", controller.GetAccountsByUser)
	r.GetFunc("/v1/api/get/role", controller.GetAllAccountRoles)

	//open API
	r.GetFunc("/v1/variants", controller.ListVariants)
	r.GetFunc("/v1/variants/:id", controller.ListVariantsDetails)
	r.GetFunc("/v1/variant/vouchers", controller.GetVoucherOfVariant)
	r.GetFunc("/v1/variant/vouchers/:id", controller.GetVoucherOfVariantDetails)

	//voucher
	r.GetFunc("/v1/vouchers", controller.GetVoucherList)
	r.GetFunc("/v1/vouchers/:id", controller.GetVoucherDetails)
	// r.PostFunc("/v1/voucher/delete", controller.DeleteVoucher)
	// r.PostFunc("/v1/voucher/pay", controller.PayVoucher)
	r.GetFunc("/v1/voucher/generate/bulk", controller.GenerateVoucherBulk)
	r.PostFunc("/v1/voucher/generate/single", controller.GenerateVoucherOnDemand)
	r.GetFunc("/v1/voucher/link", controller.GetVoucherlink)

	//public
	r.GetFunc("/v1/public/challenge", controller.GetChallenge)
	r.GetFunc("/v1/public/redeem/profile", controller.GetRedeemData)
	r.PostFunc("/v1/public/transaction", controller.WebCreateTransaction)

	//
	r.GetFunc("/v1/token", controller.GetToken)
	r.GetFunc("/v1/token/check", controller.CheckToken)

	//custom
	r.GetFunc("/view/", viewHandler)
	r.GetFunc("/viewNoLayout", viewNoLayoutHandler)

	// r.NotFoundFunc(errorHandler)

	// r.GetFunc("/test", controller.UploadFormTest)
	r.PostFunc("/file/upload", controller.UploadFile)
	r.GetFunc("/file/delete", controller.DeleteFile)
	// r.GetFunc("/listfile/", controller.GetListFile)

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
		render.FileInLayout(w, "layout.html", "variant/search.html", nil)
	} else if page == "check" {
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
	} else if page == "change-password" {
		render.FileInLayout(w, "layout.html", "user/change_pass.html", nil)
	} else if page == "login" {
		render.FileInLayout(w, "layout.html", "user/login.html", nil)
	} else if page == "profile" {
		render.FileInLayout(w, "layout.html", "user/profile.html", nil)
	} else if page == "forgot-password" {
		render.File(w, "user/forgot.html", nil)
	} else if page == "mail-send" {
		render.File(w, "user/forgot_succ.html", nil)
		//render.FileInLayout(w, "layout.html", "user/forgot.html", nil)
	} else if page == "recover" {
		render.File(w, "user/recover.html", nil)
		//render.FileInLayout(w, "layout.html", "user/recover.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "user/index.html", nil)
	}
}

func viewPartner(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "create" {
		render.FileInLayout(w, "layout.html", "partner/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "partner/search.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "partner/update.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "partner/index.html", nil)
	}
}

func viewTag(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "create" {
		render.FileInLayout(w, "layout.html", "tag/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "tag/search.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "tag/index.html", nil)
	}
}

func viewVoucher(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "create" {
		render.FileInLayout(w, "layout.html", "voucher/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "voucher/search.html", nil)
	} else if page == "check" {
		render.FileInLayout(w, "layout.html", "voucher/check.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "voucher/update.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "voucher/index.html", nil)
	}
}

func viewReport(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "variant" {
		render.FileInLayout(w, "layout.html", "report/variant.html", nil)
	} else if page == "transaction" {
		render.FileInLayout(w, "layout.html", "report/transaction.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "report/test.html", nil)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "user/login.html", nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "notfound.html", nil)
}
