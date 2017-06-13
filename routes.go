package main

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/controller"
	"github.com/gilkor/evoucher/internal/model"
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
	r.GetFunc("/public/:page", viewPublic)
	r.GetFunc("/", login)
	//r.PostFunc("/v1/query", controller.CustomQuery)

	//report
	//r.GetFunc("/v1/api/report", controller.MakeReport)
	//r.GetFunc("/v1/api/report/variant", controller.MakeReportVariant)
	//r.GetFunc("/v1/api/report/voucher/variant", controller.MakeCompleteReportVoucherByUser)
	//r.GetFunc("/v1/api/report/vouchers/variant", controller.MakeReportVoucherByUser)

	//r.PostFunc("/v1/create/variant", controller.CreateVariant)
	//r.GetFunc("/v1/api/get/allVariant", controller.GetAllVariants)
	//r.GetFunc("/v1/api/get/variant", controller.GetVariants)
	//r.GetFunc("/v1/api/get/totalVariant", controller.GetTotalVariant)
	//r.GetFunc("/v1/api/get/variantByDate", controller.GetVariantDetailsByDate)
	//r.GetFunc("/v1/api/get/variantByDate", controller.GetVariantDetailsByDate)
	//r.GetFunc("/v1/api/get/variantDetails/custom", controller.GetVariantDetailsCustom)
	//r.GetFunc("/v1/api/get/variant/:id", controller.GetVariantDetailsById)
	//r.GetFunc("/v1/ui/variant", controller.GetVariantDetailsById)
	//r.PostFunc("/v1/update/variant/:id", controller.UpdateVariant)
	//r.PostFunc("/v1/update/variant/:id/broadcastUser", controller.UpdateVariantBroadcast)
	//r.PostFunc("/v1/update/variant/:id/tenant", controller.UpdateVariantTenant)
	//r.GetFunc("/v1/delete/variant/:id", controller.DeleteVariant)

	//variant
	r.PostFunc("/v1/ui/variant/create", controller.CreateVariant)
	r.GetFunc("/v1/ui/variant/all", controller.GetAllVariants)
	r.GetFunc("/v1/ui/variant", controller.GetVariants)
	r.GetFunc("/v1/ui/variant/detail", controller.GetVariantDetailsCustom)
	r.PostFunc("/v1/ui/variant/update", controller.UpdateVariantRoute)
	r.GetFunc("/v1/ui/variant/delete", controller.DeleteVariant)

	//transaction
	r.PostFunc("/v1/transaction/redeem", controller.MobileCreateTransaction)
	r.GetFunc("/v1/ui/transaction/partner", controller.GetAllTransactionsByPartner)
	r.GetFunc("/v1/ui/transaction", controller.CashoutTransactionDetails)
	r.PostFunc("/v1/ui/transaction/cashout/update", controller.CashoutTransactions)
	r.GetFunc("/v1/ui/transaction/cashout/print", controller.PrintCashoutTransaction)

	//r.PostFunc("/v1/update/transaction/:id", controller.UpdateTransaction)
	//r.PostFunc("/v1/delete/transaction/:id", controller.DeleteTransaction)
	//r.GetFunc("/v1/get/transaction", controller.GetAllTransactions)
	//r.GetFunc("/v1/get/transaction/partner", controller.GetAllTransactionsByPartner)
	//r.GetFunc("/v1/get/transaction/cashout/:id", controller.CashoutTransactionDetails)
	//r.GetFunc("/v1/update/transaction/cashout/:id", controller.CashoutTransaction)
	//r.PostFunc("/v1/update/transaction/cashout", controller.CashoutTransactions)
	//r.GetFunc("/v1/print/transaction/cashout", controller.PrintCashoutTransaction)

	//user
	r.PostFunc("/v1/ui/user/create", controller.RegisterUser)
	r.PostFunc("/v1/ui/user/update", controller.UpdateUserRoute)
	r.GetFunc("/v1/ui/user/all", controller.GetUser)
	r.GetFunc("/v1/ui/user", controller.GetUserDetails)
	r.GetFunc("/v1/ui/user/forgot/mail", controller.SendMailForgotPassword)
	r.PostFunc("/v1/ui/user/forgot/password", controller.ForgotPassword)
	r.PostFunc("/v1/ui/user/create/broadcast", controller.InsertBroadcastUser)

	//r.PostFunc("/v1/create/user", controller.RegisterUser)
	//r.PostFunc("/v1/update/user", controller.UpdateUser)
	//r.PostFunc("/v1/update/user/password", controller.ChangePassword)
	//r.GetFunc("/v1/api/get/userByRole", controller.FindUserByRole)
	//r.GetFunc("/v1/api/get/users", controller.GetUser)
	//r.GetFunc("/v1/api/get/userDetails", controller.GetUserDetails)
	//r.GetFunc("/v1/api/mail", controller.ForgotPassword)
	//r.PostFunc("/v1/password", controller.UpdatePassword)
	//r.PostFunc("/v1/upload/user", controller.InsertBroadcastUser)

	//partner
	r.PostFunc("/v1/ui/partner/create", controller.AddPartner)
	r.GetFunc("/v1/ui/partner/all", controller.GetAllPartners)
	r.GetFunc("/v1/ui/partner", controller.GetPartners)
	r.PostFunc("/v1/ui/partner/update", controller.UpdatePartner)
	r.GetFunc("/v1/ui/partner/delete", controller.DeletePartner)

	//tag
	r.GetFunc("/v1/ui/tag/all", controller.GetAllTags)
	r.PostFunc("/v1/ui/tag/create", controller.AddTag)
	r.GetFunc("/v1/ui/tag/delete", controller.DeleteTag)
	r.PostFunc("/v1/ui/tag/delete", controller.DeleteTagBulk)

	//r.PostFunc("/v1/create/partner", controller.AddPartner)
	//r.GetFunc("/v1/get/partner", controller.GetAllPartners)
	//r.GetFunc("/v1/get/partner/:id", controller.GetPartnerDetails)
	//r.PostFunc("/v1/update/partner/:id", controller.UpdatePartner)
	//r.GetFunc("/v1/delete/partner/:id", controller.DeletePartner)

	//account
	r.GetFunc("/v1/ui/account/all", controller.GetAllAccount)
	r.GetFunc("/v1/ui/account", controller.GetAccountDetailByUser)
	r.GetFunc("/v1/ui/role/all", controller.GetAllAccountRoles)

	//r.GetFunc("/v1/api/get/account", controller.GetAccount)
	//r.GetFunc("/v1/api/get/accountsDetail", controller.GetAccountDetailByUser)
	//r.GetFunc("/v1/api/get/accountsByUser", controller.GetAccountsByUser)
	//r.GetFunc("/v1/api/get/role", controller.GetAllAccountRoles)

	//open API
	r.GetFunc("/v1/variants", controller.ListVariants)
	r.GetFunc("/v1/variants/:id", controller.ListVariantsDetails)
	r.GetFunc("/v1/variant/vouchers", controller.GetVoucherOfVariant)
	r.GetFunc("/v1/variant/vouchers/:id", controller.GetVoucherOfVariantDetails)
	r.GetFunc("/v1/api/get/partner", controller.GetAllPartnersCustomParam)

	//voucher
	r.GetFunc("/v1/ui/vouchers", controller.GetVoucherList)
	r.GetFunc("/v1/vouchers/:id", controller.GetVoucherDetails)
	r.GetFunc("/v1/ui/voucher/generate/bulk", controller.GenerateVoucherBulk)
	r.PostFunc("/v1/ui/voucher/link", controller.GetVoucherlink)
	r.GetFunc("/v1/sample/link", controller.GetCsvSample)

	//mobile API
	r.GetFunc("/v1/variants", controller.ListVariants)
	r.GetFunc("/v1/variants/:id", controller.ListVariantsDetails)
	r.GetFunc("/v1/variant/vouchers", controller.GetVoucherOfVariant)
	r.GetFunc("/v1/variant/vouchers/:id", controller.GetVoucherOfVariantDetails)
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

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	}

	if valid {
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
		} else {
			render.FileInLayout(w, "layout.html", "notfound.html", nil)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil)
	}
}

func viewUser(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	if page == "login" {
		render.FileInLayout(w, "layout.html", "user/login.html", nil)
	} else if page == "forgot-password" {
		render.File(w, "user/forgot.html", nil)
	} else if page == "mail-send" {
		render.File(w, "user/forgot_succ.html", nil)
	} else if page == "recover" {
		render.File(w, "user/recover.html", nil)
	} else {
		valid := false
		a := controller.AuthToken(w, r)
		if a.Valid {
			for _, valueRole := range a.User.Role {
				features := model.UiFeatures[valueRole.RoleDetail]
				for _, valueFeature := range features {
					if r.URL.Path == valueFeature {
						valid = true
					}
				}
			}
		}
		if valid {
			if page == "register" {
				render.FileInLayout(w, "layout.html", "user/create.html", nil)
			} else if page == "search" {
				render.FileInLayout(w, "layout.html", "user/check.html", nil)
			} else if page == "update" {
				render.FileInLayout(w, "layout.html", "user/update.html", nil)
			} else if page == "change-password" {
				render.FileInLayout(w, "layout.html", "user/change_pass.html", nil)
			} else if page == "profile" {
				render.FileInLayout(w, "layout.html", "user/profile.html", nil)
			} else {
				render.FileInLayout(w, "layout.html", "notfound.html", nil)
			}
		} else {
			render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil)
		}
	}

}

func viewPartner(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	}
	if valid {
		if page == "create" {
			render.FileInLayout(w, "layout.html", "partner/create.html", nil)
		} else if page == "search" {
			render.FileInLayout(w, "layout.html", "partner/search.html", nil)
		} else if page == "update" {
			render.FileInLayout(w, "layout.html", "partner/update.html", nil)
		} else {
			render.FileInLayout(w, "layout.html", "notfound.html", nil)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil)
	}
}

func viewTag(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	}
	if valid {
		if page == "search" {
			render.FileInLayout(w, "layout.html", "tag/search.html", nil)
		} else {
			render.FileInLayout(w, "layout.html", "notfound.html", nil)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil)
	}
}

func viewVoucher(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	}
	if valid {
		if page == "search" {
			render.FileInLayout(w, "layout.html", "voucher/search.html", nil)
		} else if page == "check" {
			render.FileInLayout(w, "layout.html", "voucher/check.html", nil)
		} else if page == "cashout" {
			render.FileInLayout(w, "layout.html", "voucher/cashout.html", nil)
		} else if page == "print" {
			render.FileInLayout(w, "layout.html", "voucher/print.html", nil)
		} else {
			render.FileInLayout(w, "layout.html", "notfound.html", nil)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil)
	}
}

func viewReport(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.RoleDetail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	}
	if valid {
		if page == "variant" {
			render.FileInLayout(w, "layout.html", "report/variant.html", nil)
		} else if page == "transaction" {
			render.FileInLayout(w, "layout.html", "report/transaction.html", nil)
		} else if page == "" || page == "index" {
			render.FileInLayout(w, "layout.html", "report/test.html", nil)
		} else {
			render.FileInLayout(w, "layout.html", "notfound.html", nil)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil)
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
		render.FileInLayout(w, "layout.html", "public/check.html", nil)
	} else if page == "" || page == "index" {
		render.File(w, "public/index.html", nil)
	} else {
		render.FileInLayout(w, "layout.html", "notfound.html", nil)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "user/login.html", nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	render.FileInLayout(w, "layout.html", "notfound.html", nil)
}
