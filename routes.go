package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"fmt"
	"github.com/gilkor/evoucher/internal/controller"
)

func setRoutes() http.Handler {
	r := bone.New()
	// http.ListenAndServe(":8888", nil)
	r.NotFoundFunc(errorHandler)
	r.GetFunc("/ping", ping)

	//ui
	r.GetFunc("/program/:page", viewProgram)
	r.GetFunc("/account/:page", viewAccount)
	r.GetFunc("/user/:page", viewUser)
	r.GetFunc("/partner/:page", viewPartner)
	r.GetFunc("/tag/:page", viewTag)
	r.GetFunc("/voucher/:page", viewVoucher)
	r.GetFunc("/report/:page", viewReport)
	r.GetFunc("/public/:page", viewPublic)
	r.GetFunc("/sa/:page", viewSuperAdmin)
	r.GetFunc("/role/:page", viewRole)
	r.GetFunc("/unauthorize", viewUnauthorize)
	r.GetFunc("/notfound", viewNotFound)
	r.GetFunc("/", login)

	//program
	r.PostFunc("/v1/ui/program/create", controller.CreateProgram)
	r.GetFunc("/v1/ui/program/all", controller.GetAllPrograms)
	r.GetFunc("/v1/ui/program/ongoing", controller.GetOnGoingPrograms)
	r.GetFunc("/v1/ui/program", controller.GetPrograms)
	r.GetFunc("/v1/ui/program/detail", controller.GetProgramDetailsCustom)
	r.PostFunc("/v1/ui/program/update", controller.UpdateProgramRoute)
	r.GetFunc("/v1/ui/program/delete", controller.DeleteProgram)
	r.GetFunc("/v1/ui/program/partner", controller.GetProgramPartnerSummary)
	r.GetFunc("/v1/ui/program/visibility", controller.VisibilityProgram)

	//campaign
	r.PostFunc("/v1/ui/campaign/create", controller.CreateEmailCampaign)

	//transaction
	r.GetFunc("/v1/ui/transaction/partner", controller.GetTransactionsByPartner)
	r.GetFunc("/v1/ui/transaction/date", controller.GetTransactionsByDate)
	r.GetFunc("/v1/ui/transaction/cashout/partner", controller.GetTransactionsByPartner)
	r.GetFunc("/v1/ui/transaction/voucher", controller.GetVoucherTransactionDetails)
	r.GetFunc("/v1/ui/transaction/cashout", controller.CashoutTransactionDetails)

	//cashout
	r.PostFunc("/v1/ui/cashout", controller.CashoutTransactions)
	r.GetFunc("/v1/ui/cashout/print", controller.PrintCashoutTransaction)
	r.GetFunc("/v1/ui/voucher/partner", controller.GetVouchersByPartner)
	r.GetFunc("/v1/ui/voucher/daily/partner", controller.GetTodayVouchersByPartner)

	//user
	r.GetFunc("/v1/ui/user/login", controller.Login)
	r.PostFunc("/v1/ui/user/create", controller.RegisterUser)
	r.PostFunc("/v1/ui/user/update", controller.UpdateUserRoute)
	r.PostFunc("/v1/ui/user/block", controller.BlockUser)
	r.PostFunc("/v1/ui/user/activate", controller.ActivateUser)
	r.GetFunc("/v1/ui/user/all", controller.GetUser)
	r.GetFunc("/v1/ui/user", controller.GetUserDetails)
	r.GetFunc("/v1/ui/user/other", controller.GetOtherUserDetails)
	r.GetFunc("/v1/ui/user/forgot/mail", controller.SendForgotPasswordMail)
	r.PostFunc("/v1/ui/user/forgot/password", controller.ForgotPassword)
	r.PostFunc("/v1/ui/user/create/broadcast", controller.InsertBroadcastUser)

	//sa
	r.PostFunc("/v1/ui/sa/create", controller.SuperadminRegisterUser)
	r.GetFunc("/v1/ui/sa/all", controller.SuperadminGetUser)
	r.PostFunc("/v1/ui/sa/a-create", controller.RegisterAccount)
	r.PostFunc("/v1/ui/sa/a-update", controller.UpdateAccount)
	r.GetFunc("/v1/ui/sa/account", controller.GetAllAccountsDetail)
	r.PostFunc("/v1/ui/sa/a-block", controller.BlockAccount)
	r.PostFunc("/v1/ui/sa/a-activate", controller.ActivateAccount)
	//r.PostFunc("/v1/ui/sa/update", controller.UpdateUserRoute)
	//r.GetFunc("/v1/ui/sa/forgot/mail", controller.SendForgotPasswordMail)
	//r.PostFunc("/v1/ui/sa/forgot/password", controller.ForgotPassword)

	//partner
	r.PostFunc("/v1/ui/partner/create", controller.AddPartner)
	r.GetFunc("/v1/ui/partner/all", controller.GetAllPartners)
	r.GetFunc("/v1/ui/partner/program", controller.GetProgramPartners)
	r.GetFunc("/v1/ui/partner", controller.GetPartners)
	r.GetFunc("/v1/ui/partner/programs", controller.GetProgramsPartner)
	r.GetFunc("/v1/ui/partner/performance", controller.GetPerformancePartner)
	r.GetFunc("/v1/ui/partner/daily/performance", controller.GetDailyPerformancePartner)
	r.PostFunc("/v1/ui/partner/update", controller.UpdatePartner)
	r.GetFunc("/v1/ui/partner/delete", controller.DeletePartner)
	r.PostFunc("/v1/ui/bank_account/create", controller.RegisterBankAccount)
	r.GetFunc("/v1/ui/bank_account/all", controller.GetAllBankAccounts)

	//tag
	r.GetFunc("/v1/ui/tag/all", controller.GetAllTags)
	r.PostFunc("/v1/ui/tag/create", controller.AddTag)
	r.GetFunc("/v1/ui/tag/delete", controller.DeleteTag)
	r.PostFunc("/v1/ui/tag/delete", controller.DeleteTagBulk)

	//account
	r.GetFunc("/v1/ui/account/all", controller.GetAllAccount)
	r.GetFunc("/v1/ui/account", controller.GetAccountDetailByUser)
	r.GetFunc("/v1/ui/account/other", controller.GetAccountDetailByOtherUser)

	//role
	r.GetFunc("/v1/ui/role/all", controller.GetAllAccountRoles)
	r.GetFunc("/v1/ui/role/account", controller.GetAccountRoles)
	r.GetFunc("/v1/ui/feature/all", controller.GetAllFeatures)
	r.GetFunc("/v1/ui/role/detail", controller.GetFeaturesDetail)
	r.PostFunc("/v1/ui/role/create", controller.AddRole)
	r.PostFunc("/v1/ui/role/update", controller.UpdateRole)

	//open API
	r.GetFunc("/v1/api/get/partner", controller.GetAllPartnersCustomParam)

	//voucher
	r.GetFunc("/v1/ui/voucher", controller.GetVoucherList)
	r.GetFunc("/v1/ui/voucher/:id", controller.GetVoucherDetails)
	r.PostFunc("/v1/ui/voucher/generate/bulk", controller.GenerateVoucherBulk)
	r.PostFunc("/v1/ui/voucher/link", controller.GetVoucherlink)
	r.PostFunc("/v1/ui/voucher/send-voucher", controller.SendSedayuOneEmail)
	r.GetFunc("/v1/ui/sample/link", controller.GetCsvSample)

	//mobile API
	r.GetFunc("/v1/program", controller.ListMobilePrograms)
	r.GetFunc("/v1/mall-program", controller.ListMallPrograms)
	r.GetFunc("/v1/program/:id", controller.ListProgramsDetails)
	r.GetFunc("/v1/voucher", controller.GetVoucherOfProgram)
	r.GetFunc("/v1/voucher/:id", controller.GetVoucherOfProgramDetails)
	r.GetFunc("/v1/voucher/generate/single/:id/rollback", controller.RollbackVoucher)
	r.PostFunc("/v1/voucher/generate/single", controller.GenerateVoucherOnDemand)
	r.PostFunc("/v1/transaction/redeem", controller.MobileCreateTransaction)

	//public API
	r.GetFunc("/v1/public/challenge", controller.GetChallenge)
	r.GetFunc("/v1/public/redeem/profile", controller.GetRedeemData)
	r.PostFunc("/v1/public/transaction", controller.WebCreateTransaction)
	r.GetFunc("/v1/public/transaction/:id", controller.PublicCashoutTransactionDetails)

	//auth
	r.GetFunc("/v1/token", controller.GetToken)
	r.GetFunc("/v1/token/check", controller.CheckToken)
	r.GetFunc("/v1/ui/token/check", controller.UICheckToken)

	//custom
	r.GetFunc("/view/", viewHandler)
	r.GetFunc("/viewNoLayout", viewNoLayoutHandler)
	r.PostFunc("/v1/send/mail", controller.SendCustomMailRoute)

	// r.GetFunc("/test", controller.UploadFormTest)
	r.PostFunc("/file/upload", controller.UploadFile)
	r.GetFunc("/file/delete", controller.DeleteFile)

	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ping")
	w.Write([]byte("ping"))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "partner/check.html", nil)
}

func viewNoLayoutHandler(w http.ResponseWriter, r *http.Request) {
	render.File(w, "view/noLayout.html", nil)
}

func viewProgram(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "create" {
		render.FileInLayout(w, "layout.html", "program/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "program/search.html", nil)
	} else if page == "check" {
		render.FileInLayout(w, "layout.html", "program/check.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "program/update.html", nil)
	} else if page == "campaign" {
		render.FileInLayout(w, "layout.html", "program/campaign.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "program/index.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}

}

func viewUser(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "login" {
		render.FileInLayout(w, "layout.html", "user/login.html", nil)
	} else if page == "forgot-password" {
		render.File(w, "user/forgot/forgot.html", nil)
	} else if page == "mail-send" {
		render.File(w, "user/forgot/forgot_succ.html", nil)
	} else if page == "recover" {
		render.File(w, "user/forgot/recover.html", nil)
	} else if page == "register" {
		render.FileInLayout(w, "layout.html", "user/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "user/search.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "user/update.html", nil)
	} else if page == "change-password" {
		render.FileInLayout(w, "layout.html", "user/change_pass.html", nil)
	} else if page == "profile" {
		render.FileInLayout(w, "layout.html", "user/profile.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}

}

func viewSuperAdmin(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "register" {
		render.FileInLayout(w, "layout.html", "superadmin/create_user.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "superadmin/search_user.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "user/update.html", nil)
	} else if page == "change-password" {
		render.FileInLayout(w, "layout.html", "user/change_pass.html", nil)
	} else if page == "a-create" {
		render.FileInLayout(w, "layout.html", "superadmin/create_account.html", nil)
	} else if page == "a-search" {
		render.FileInLayout(w, "layout.html", "superadmin/search_account.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewPartner(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "create" {
		render.FileInLayout(w, "layout.html", "partner/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "partner/search.html", nil)
	} else if page == "check" {
		render.FileInLayout(w, "layout.html", "partner/check.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "partner/update.html", nil)
	} else if page == "bank" {
		render.FileInLayout(w, "layout.html", "bank_account/create.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewAccount(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "create" {
		render.FileInLayout(w, "layout.html", "account/create.html", nil)
	} else if page == "search" {
		render.FileInLayout(w, "layout.html", "account/search.html", nil)
	} else if page == "check" {
		render.FileInLayout(w, "layout.html", "account/check.html", nil)
	} else if page == "update" {
		render.FileInLayout(w, "layout.html", "account/update.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewTag(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "search" {
		render.FileInLayout(w, "layout.html", "tag/search.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewRole(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "search" {
		render.FileInLayout(w, "layout.html", "role/search.html", nil)
	} else if page == "create" {
		render.FileInLayout(w, "layout.html", "role/create.html", nil)
	} else if page == "edit" {
		render.FileInLayout(w, "layout.html", "role/edit.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewVoucher(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "search" {
		render.FileInLayout(w, "layout.html", "voucher/search.html", nil)
	} else if page == "check" {
		render.FileInLayout(w, "layout.html", "voucher/check.html", nil)
	} else if page == "cashout" {
		render.FileInLayout(w, "layout.html", "voucher/cashout.html", nil)
	} else if page == "cashout_detail" {
		render.FileInLayout(w, "layout.html", "voucher/cashout_detail.html", nil)
	} else if page == "print" {
		render.FileInLayout(w, "layout.html", "voucher/print.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewReport(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "program" {
		render.FileInLayout(w, "layout.html", "report/program.html", nil)
	} else if page == "transaction" {
		render.FileInLayout(w, "layout.html", "report/transaction.html", nil)
	} else if page == "" || page == "index" {
		render.FileInLayout(w, "layout.html", "report/test.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewPublic(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "fail" {
		render.File(w, "public/fail.html", nil)
	} else if page == "success" {
		render.File(w, "public/success.html", nil)
	} else if page == "redeem" {
		render.File(w, "public/index.html", nil)
	} else if page == "check" {
		render.File(w, "public/check.html", nil)
	} else if page == "" || page == "index" {
		render.File(w, "public/index.html", nil)
	} else {
		render.File(w, "notfound.html", nil, 404)
	}
}

func viewUnauthorize(w http.ResponseWriter, r *http.Request) {
	render.File(w, "notfound.html", nil, 401)
}

func viewNotFound(w http.ResponseWriter, r *http.Request) {
	render.File(w, "notfound.html", nil, 404)
}

func login(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "user/login.html", nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	render.File(w, "notfound.html", nil, 404)
}
