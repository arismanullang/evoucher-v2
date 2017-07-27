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
	r.GetFunc("/program/:page", viewProgram)
	r.GetFunc("/user/:page", viewUser)
	r.GetFunc("/partner/:page", viewPartner)
	r.GetFunc("/tag/:page", viewTag)
	r.GetFunc("/voucher/:page", viewVoucher)
	r.GetFunc("/report/:page", viewReport)
	r.GetFunc("/public/:page", viewPublic)
	r.GetFunc("/", login)

	//report
	//r.GetFunc("/v1/api/report", controller.MakeReport)
	//r.GetFunc("/v1/api/report/Program", controller.MakeReportProgram)
	//r.GetFunc("/v1/api/report/voucher/Program", controller.MakeCompleteReportVoucherByUser)
	//r.GetFunc("/v1/api/report/vouchers/Program", controller.MakeReportVoucherByUser)

	//program
	r.PostFunc("/v1/ui/program/create", controller.CreateProgram)
	r.GetFunc("/v1/ui/program/all", controller.GetAllPrograms)
	r.GetFunc("/v1/ui/program", controller.GetPrograms)
	r.GetFunc("/v1/ui/program/detail", controller.GetProgramDetailsCustom)
	r.PostFunc("/v1/ui/program/update", controller.UpdateProgramRoute)
	r.GetFunc("/v1/ui/program/delete", controller.DeleteProgram)

	//transaction
	r.GetFunc("/v1/ui/transaction/partner", controller.GetAllTransactionsByPartner)
	r.GetFunc("/v1/ui/transaction", controller.CashoutTransactionDetails)
	r.PostFunc("/v1/ui/transaction/cashout/update", controller.CashoutTransactions)
	r.GetFunc("/v1/ui/transaction/cashout/print", controller.PrintCashoutTransaction)

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

	//partner
	r.PostFunc("/v1/ui/partner/create", controller.AddPartner)
	r.GetFunc("/v1/ui/partner/all", controller.GetAllPartners)
	r.GetFunc("/v1/ui/partner/program", controller.GetProgramPartners)
	r.GetFunc("/v1/ui/partner", controller.GetPartners)
	r.GetFunc("/v1/ui/partner/performance", controller.GetPerformancePartner)
	r.PostFunc("/v1/ui/partner/update", controller.UpdatePartner)
	r.GetFunc("/v1/ui/partner/delete", controller.DeletePartner)

	//tag
	r.GetFunc("/v1/ui/tag/all", controller.GetAllTags)
	r.PostFunc("/v1/ui/tag/create", controller.AddTag)
	r.GetFunc("/v1/ui/tag/delete", controller.DeleteTag)
	r.PostFunc("/v1/ui/tag/delete", controller.DeleteTagBulk)

	//account
	r.GetFunc("/v1/ui/account/all", controller.GetAllAccount)
	r.GetFunc("/v1/ui/account", controller.GetAccountDetailByUser)
	r.GetFunc("/v1/ui/role/all", controller.GetAllAccountRoles)

	//open API
	r.GetFunc("/v1/api/get/partner", controller.GetAllPartnersCustomParam)

	//voucher
	r.GetFunc("/v1/ui/voucher", controller.GetVoucherList)
	r.GetFunc("/v1/ui/voucher/:id", controller.GetVoucherDetails)
	r.PostFunc("/v1/ui/voucher/generate/bulk", controller.GenerateVoucherBulk)
	r.PostFunc("/v1/ui/voucher/link", controller.GetVoucherlink)
	r.PostFunc("/v1/ui/voucher/send-voucher", controller.SendVoucherBulk)
	r.GetFunc("/v1/sample/link", controller.GetCsvSample)

	//mobile API
	r.GetFunc("/v1/program", controller.ListPrograms)
	r.GetFunc("/v1/program/:id", controller.ListProgramsDetails)
	r.GetFunc("/v1/voucher", controller.GetVoucherOfProgram)
	r.GetFunc("/v1/voucher/:id", controller.GetVoucherOfProgramDetails)
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
	r.PostFunc("/v1/send/mail", controller.SendCustomMailRoute)

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
	render.FileInLayout(w, "layout.html", "partner/check.html", nil)
}

func viewNoLayoutHandler(w http.ResponseWriter, r *http.Request) {
	render.File(w, "view/noLayout.html", nil)
}

func viewProgram(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.Detail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	} else {
		render.File(w, "notfound.html", nil, 404)
		return
	}

	if valid {
		if page == "create" {
			render.FileInLayout(w, "layout.html", "program/create.html", nil)
		} else if page == "search" {
			render.FileInLayout(w, "layout.html", "program/search.html", nil)
		} else if page == "check" {
			render.FileInLayout(w, "layout.html", "program/check.html", nil)
		} else if page == "update" {
			render.FileInLayout(w, "layout.html", "program/update.html", nil)
		} else if page == "" || page == "index" {
			render.FileInLayout(w, "layout.html", "program/index.html", nil)
		} else {
			render.File(w, "notfound.html", nil, 404)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil, 401)
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
	} else {
		valid := false
		a := controller.AuthToken(w, r)
		if a.Valid {
			for _, valueRole := range a.User.Role {
				features := model.UiFeatures[valueRole.Detail]
				for _, valueFeature := range features {
					if r.URL.Path == valueFeature {
						valid = true
					}
				}
			}
		} else {
			render.File(w, "notfound.html", nil, 404)
			return
		}

		if valid {
			if page == "register" {
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
		} else {
			render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil, 401)
		}
	}

}

func viewPartner(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.Detail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	} else {
		render.File(w, "notfound.html", nil, 404)
		return
	}

	if valid {
		if page == "create" {
			render.FileInLayout(w, "layout.html", "partner/create.html", nil)
		} else if page == "search" {
			render.FileInLayout(w, "layout.html", "partner/search.html", nil)
		} else if page == "check" {
			render.FileInLayout(w, "layout.html", "partner/check.html", nil)
		} else if page == "update" {
			render.FileInLayout(w, "layout.html", "partner/update.html", nil)
		} else {
			render.File(w, "notfound.html", nil, 404)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil, 401)
	}
}

func viewTag(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.Detail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	} else {
		render.File(w, "notfound.html", nil, 404)
		return
	}

	if valid {
		if page == "search" {
			render.FileInLayout(w, "layout.html", "tag/search.html", nil)
		} else {
			render.File(w, "notfound.html", nil, 404)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil, 401)
	}
}

func viewVoucher(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.Detail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	} else {
		render.File(w, "notfound.html", nil, 404)
		return
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
			render.File(w, "notfound.html", nil, 404)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil, 401)
	}
}

func viewReport(w http.ResponseWriter, r *http.Request) {
	page := bone.GetValue(r, "page")

	valid := false
	a := controller.AuthToken(w, r)
	if a.Valid {
		for _, valueRole := range a.User.Role {
			features := model.UiFeatures[valueRole.Detail]
			for _, valueFeature := range features {
				if r.URL.Path == valueFeature {
					valid = true
				}
			}
		}
	} else {
		render.File(w, "notfound.html", nil, 404)
		return
	}

	if valid {
		if page == "program" {
			render.FileInLayout(w, "layout.html", "report/program.html", nil)
		} else if page == "transaction" {
			render.FileInLayout(w, "layout.html", "report/transaction.html", nil)
		} else if page == "" || page == "index" {
			render.FileInLayout(w, "layout.html", "report/test.html", nil)
		} else {
			render.File(w, "notfound.html", nil, 404)
		}
	} else {
		render.FileInLayout(w, "layout.html", "user/unauthorize.html", nil, 401)
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

func login(w http.ResponseWriter, r *http.Request) {
	render.FileInLayout(w, "layout.html", "user/login.html", nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	render.File(w, "notfound.html", nil, 404)
}
