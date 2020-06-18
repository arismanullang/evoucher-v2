package main

import (
	"fmt"
	"net/http"

	c "github.com/gilkor/evoucher-v2/controller"
	"github.com/go-zoo/bone"
)

var router http.Handler

func init() {
	//main router
	r := bone.New()
	// AutoCORS(r)
	// r.NotFoundFunc(notFound)
	r.GetFunc("/", healthCheck)
	r.GetFunc("/ping", ping)

	//define sub router
	v2 := bone.New()
	r.SubRoute("/v2/api", v2)
	// AutoCORS(v2)

	//voucher
	v2.PostFunc("/:company/web/vouchers", ping)

	v2.Get("/:company/web/vouchers/program/:id", c.CheckFuncJWT(c.GetVoucherByProgramID, "voucher-view"))
	v2.Get("/:company/web/vouchers/:id", c.CheckFuncJWT(c.GetVoucherByID, "voucher-view"))

	v2.PutFunc("/:company/web/vouchers/:id", ping)
	v2.Delete("/:company/web/vouchers/:id", c.CheckFuncJWT(c.DeleteVoucher, "voucher-delete"))

	//programs
	v2.Post("/:company/web/programs", c.CheckFuncJWT(c.PostProgram, "program-create"))
	v2.Get("/:company/web/programs", c.CheckFuncJWT(c.GetProgram, "program-view"))
	v2.Get("/:company/web/programs/:id", c.CheckFuncJWT(c.GetProgramByID, "program-view"))
	v2.Post("/:company/web/programs/image/:id", c.CheckFuncJWT(c.UploadProgramImage, "program-edit"))
	v2.Put("/:company/web/programs/:id", c.CheckFuncJWT(c.UpdateProgram, "program-edit"))
	v2.Delete("/:company/web/programs/:id", c.CheckFuncJWT(c.DeleteProgram, "program_delete"))

	//outlets
	v2.Post("/:company/web/outlets", c.CheckFuncJWT(c.PostOutlet, "outlet-create"))
	v2.Get("/:company/web/outlets", c.CheckFuncJWT(c.GetOutlets, "outlet-view"))
	v2.Get("/:company/web/outlets/:id", c.CheckFuncJWT(c.GetOutletByID, "outlet-view"))
	v2.Put("/:company/web/outlets/:id", c.CheckFuncJWT(c.UpdateOutlet, "outlet-edit"))
	v2.Delete("/:company/web/outlets/:id", c.CheckFuncJWT(c.DeleteOutlet, "outlet-delete"))

	v2.Get("/:company/web/outlets/tags/:tag_id", c.CheckFuncJWT(c.GetOutletByTags, "tag-view"))
	v2.Post("/:company/web/outlets/tags/:holder", c.CheckFuncJWT(c.PostOutletTags, "tag-edit"))

	v2.GetFunc("/:company/web/outlets/tags/:tag_id", c.GetOutletByTags)
	v2.PostFunc("/:company/web/outlets/tags/:holder", c.PostOutletTags)

	// outlet bank
	v2.Post("/:company/web/banks/:pid", c.CheckFuncJWT(c.PostBank, "outlet-create"))
	v2.Get("/:company/web/banks", c.CheckFuncJWT(c.GetBanks, "outlet-view"))
	v2.Get("/:company/web/banks/:pid", c.CheckFuncJWT(c.GetBankByOutletID, "outlet-view"))
	v2.Put("/:company/web/banks/:pid", c.CheckFuncJWT(c.UpdateBank, "outlet-edit"))
	v2.Delete("/:company/web/banks/:pid", c.CheckFuncJWT(c.DeleteBank, "outlet-delete"))

	//channel
	v2.Post("/:company/web/channels", c.CheckFuncJWT(c.PostChannel, "channel-create"))
	v2.Get("/:company/web/channels", c.CheckFuncJWT(c.GetChannels, "channel-view"))
	v2.Get("/:company/web/channels/:id", c.CheckFuncJWT(c.GetChannelByID, "channel-view"))
	v2.Put("/:company/web/channels/:id", c.CheckFuncJWT(c.UpdateChannel, "channel-edit"))
	v2.Delete("/:company/web/channels/:id", c.CheckFuncJWT(c.DeleteChannel, "channel-delete"))

	//users
	// v2.GetFunc("/:company/web/login", ping)

	//tags
	v2.Post("/:company/web/tags", c.CheckFuncJWT(c.PostTag, "tag-create"))
	v2.Get("/:company/web/tags", c.CheckFuncJWT(c.GetTags, "tag-view"))
	v2.Get("/:company/web/tags/:id", c.CheckFuncJWT(c.GetTagByID, "tag-view"))
	v2.Get("/:company/web/tags/key/:key", c.CheckFuncJWT(c.GetTagByKey, "tag-view"))
	v2.Get("/:company/web/tags/category/:category", c.CheckFuncJWT(c.GetTagByKey, "tag-view"))
	// v2.GetFunc("/:company/web/tags/category/:key", c.GetTagByKey)
	v2.Put("/:company/web/tags/:id", c.CheckFuncJWT(c.UpdateTag, "tag-edit"))
	v2.Delete("/:company/web/tags/:id", c.CheckFuncJWT(c.DeleteTag, "tag-delete"))

	v2.Post("/:company/web/tags/assign/:id", c.CheckFuncJWT(c.PostObjectTags, "tag-edit"))
	v2.Post("/:company/web/tags/assign", c.CheckFuncJWT(c.PostObjectTags, "tag-edit"))

	//customers
	v2.PostFunc("/:company/web/customers", c.PostCustomer)
	v2.GetFunc("/:company/web/customers", c.GetCustomer)
	v2.GetFunc("/:company/web/customers/:id", c.GetCustomerByID)
	v2.PutFunc("/:company/web/customers/:id", c.UpdateCustomer)
	v2.DeleteFunc("/:company/web/customers/:id", c.DeleteCustomer)

	v2.PostFunc("/:company/web/customers/tags/:id", c.PostCustomerTags)

	//transaction voucher
	v2.Post("/:company/web/transaction/voucher/inject/holder", c.CheckFuncJWT(c.PostVoucherInjectByHolder, "transaction-create"))

	v2.Get("/:company/web/transaction", c.CheckFuncJWT(c.GetTransactions, "transaction-view")) // cumicumi
	v2.Get("/:company/web/transaction/outlet/:id", c.CheckFuncJWT(c.GetTransactionsByOutlet, "transaction-view"))
	v2.Get("/:company/web/transaction/program/:id", c.CheckFuncJWT(c.GetTransactionsByProgram, "transaction-view"))
	v2.Get("/:company/web/transaction/:id", c.CheckFuncJWT(c.GetTransactionByID, "transaction-view")) // cumicumi

	// v2.Get("/:company/web/reimburse/summary", c.CheckFuncJWT(c.GetCashoutSummary, "reimburse-view"))
	// v2.Get("/:company/web/reimburse/list", c.CheckFuncJWT(c.GetCashouts, "reimburse-view"))
	// v2.Post("/:company/web/reimburse/partner", c.CheckFuncJWT(c.PostCashoutByPartner, "reimburse-create"))
	// v2.Get("/:company/web/reimburse/voucher/:program_id", c.CheckFuncJWT(c.GetCashoutUsedVoucher, "reimburse-create"))
	// v2.Get("/:company/web/reimburse/unpaid/", c.CheckFuncJWT(c.GetCashoutsUnpaid, "reimburse-create")) -> still grouped by date

	v2.Get("/:company/web/reimburse/unpaid/list", c.CheckFuncJWT(c.GetUnpaidReimburse, "transaction-view"))
	v2.Get("/:company/web/reimburse/unpaid/vouchers/:outlet_id", c.CheckFuncJWT(c.GetUnpaidVouchersByOutlet, "transaction-view"))
	v2.Post("/:company/web/reimburse/create", c.CheckFuncJWT(c.PostCashout, "transaction-view"))
	v2.Get("/:company/web/reimburse/vouchers/:id", c.CheckFuncJWT(c.GetCashoutVouchers, "transaction-view"))
	v2.Get("/:company/web/reimburse/:id", c.CheckFuncJWT(c.GetCashoutByID, "transaction-view"))
	v2.Get("/:company/web/reimburse/paid/list", c.CheckFuncJWT(c.GetCashouts, "transaction-view"))
	// v2.Post("/:company/web/reimburse/void", c.CheckFuncJWT(c.VoidCashout, "reimburse-view"))
	// v2.Put("/:company/web/reimburse/approval", c.CheckFuncJWT(c.PutApproveReimburse, "reimburse-approval"))

	// v2.GetFunc("/:company/web/vouchers/report", c.GetVoucherReport)

	// v2.GetFunc("/:company/web/debug/pprof/", pprof.Index)
	// v2.GetFunc("/:company/web/debug/pprof/cmdline", pprof.Cmdline)
	// v2.GetFunc("/:company/web/debug/pprof/profile", pprof.Profile)
	// v2.GetFunc("/:company/web/debug/pprof/symbol", pprof.Symbol)
	// v2.GetFunc("/:company/web/debug/pprof/trace", pprof.Trace)

	//Dashboard
	v2.GetFunc("/:company/dashboard/topoutlet", c.GetDashboardTopOutlet)
	v2.GetFunc("/:company/dashboard/topprogram", c.GetDashboardTopProgram)
	v2.GetFunc("/:company/dashboard/voucherusage", c.GetDashboardVoucherUsage)

	//REPORT
	v2.GetFunc("/:company/report/outlet-transaction", c.GetReportDailyVoucherTransactionWithOutlet)
	v2.GetFunc("/:company/report/topoutlet", c.GetReportDailyVoucherTransaction)
	v2.GetFunc("/:company/report/", c.GetReportDailyOutletTransaction)
	v2.GetFunc("/:company/report/topoutlet", c.GetReportMonthlyOutletTransaction)
	v2.GetFunc("/:company/report/topoutlet", c.GetReportYearlyOutletTransaction)

	// blast
	v2.Get("/:company/web/blasts", c.CheckFuncJWT(c.GetBlasts, "blast-view"))
	v2.Get("/:company/web/blasts/:id", c.CheckFuncJWT(c.GetBlastByID, "blast-view"))
	v2.Post("/:company/web/blasts/create", c.CheckFuncJWT(c.CreateEmailBlast, "blast-create"))
	v2.Post("/:company/web/blasts/send/:id", c.CheckFuncJWT(c.SendEmailBlast, "blast-edit"))
	v2.Get("/:company/web/blasts/template/:id", c.CheckFuncJWT(c.GetBlastsTemplate, "blast-view"))
	v2.Get("/:company/web/template/:name", c.CheckFuncJWT(c.GetTemplateByName, "blast-view"))

	// config
	v2.Get("/:company/web/config", c.CheckFuncJWT(c.GetConfigs, "setting-view"))
	v2.Post("/:company/web/config/:category", c.CheckFuncJWT(c.SetConfig, "setting-create"))
	v2.Put("/:company/web/config", c.CheckFuncJWT(c.UpdateConfig, "setting-edit"))

	// public
	v2.GetFunc("/:company/web/public/voucher", c.GetPublicVoucherByID)      // cumicumi
	v2.PostFunc("/:company/web/public/voucher/use", c.PostPublicVoucherUse) // cumicumi

	// Mobile / 3rd party
	v2.Get("/:company/programs/channel/:channel_id", c.CheckFuncJWT(c.GetProgramsByChannel, "program-view"))
	v2.Get("/:company/programs/:id", c.CheckFuncJWT(c.GetMProgramByID, "program-view"))
	v2.Post("/:company/transaction/voucher/claim", c.CheckFuncJWT(c.PostVoucherClaim, "program-view"))
	v2.Post("/:company/transaction/voucher/use", c.CheckFuncJWT(c.PostVoucherUse, "program-view"))
	v2.Get("/:company/vouchers", c.CheckFuncJWT(c.GetVoucherByToken, "program-view"))
	v2.Get("/:company/transaction/history/:holder", c.CheckFuncJWT(c.GetHolderTrxHistory, "transaction-view"))
	v2.Get("/:company/transaction/history/:holder/:trx_id", c.CheckFuncJWT(c.GetHolderTrxHistoryDetail, "transaction-view"))

	// GCS
	v2.Post("/:company/web/file/upload", c.CheckFuncJWT(c.UploadFile, "file-create"))
	v2.Get("/:company/web/file/delete", c.CheckFuncJWT(c.DeleteFile, "file-delete"))

	v2.Get("/:company/web/accounts", c.CheckFuncJWT(c.GetAccounts, "member-view"))
	v2.Get("/:company/web/accounts/:id", c.CheckFuncJWT(c.GetAccountByID, "member-view"))

	router = r
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ping")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not found"))
}

func AutoCORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Accept, Content-Type")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkJunoToken() {
	//token, err := VerifyJWT(key)
	//if err != nil {
	//	res.SetStatus(http.StatusUnauthorized)
	//	res.AddErrors(err)
	//	res.Write(w)
	//	return
	//}
	//
	//claims, ok := token.Claims.(*JWTJunoClaims)
	//if ok && token.Valid {
	//	// fmt.Printf("Key:%v", token.Header)
	//} else {
	//	res.SetStatus(http.StatusUnauthorized)
	//	res.AddErrors(errors.New("key is invalid"))
	//	res.Write(w)
	//	return
	//}
}
