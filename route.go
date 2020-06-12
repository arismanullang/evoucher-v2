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
	v2.PostFunc("/:company/vouchers", ping)
	v2.Get("/:company/vouchers", c.CheckFuncJWT(c.GetVoucherByHolder, "voucher-view"))
	v2.Get("/:company/vouchers/program/:id", c.CheckFuncJWT(c.GetVoucherByProgramID, "voucher-view"))
	v2.GetFunc("/:company/vouchers/:id", ping)
	v2.PutFunc("/:company/vouchers/:id", ping)
	v2.Delete("/:company/vouchers/:id", c.CheckFuncJWT(c.DeleteVoucher, "voucher-delete"))

	//programs
	v2.Post("/:company/programs", c.CheckFuncJWT(c.PostProgram, "program-create"))
	v2.Get("/:company/programs", c.CheckFuncJWT(c.GetProgram, "program-view"))
	v2.Get("/:company/programs/:id", c.CheckFuncJWT(c.GetProgramByID, "program-view"))
	v2.Post("/:company/programs/image/:id", c.CheckFuncJWT(c.UploadProgramImage, "program-edit"))
	v2.Put("/:company/programs/:id", c.CheckFuncJWT(c.UpdateProgram, "program-edit"))
	v2.Delete("/:company/programs/:id", c.CheckFuncJWT(c.DeleteProgram, "program_delete"))

	//partners == outlets
	v2.Post("/:company/partners", c.CheckFuncJWT(c.PostPartner, "outlet-create"))
	v2.Get("/:company/partners", c.CheckFuncJWT(c.GetPartners, "outlet-view"))
	// v2.GetFunc("/:company/partners", c.GetPartners)
	v2.Get("/:company/partners/:id", c.CheckFuncJWT(c.GetPartnerByID, "outlet-view"))
	v2.Put("/:company/partners/:id", c.CheckFuncJWT(c.UpdatePartner, "outlet-edit"))
	v2.Delete("/:company/partners/:id", c.CheckFuncJWT(c.DeletePartner, "outlet-delete"))

	v2.Get("/:company/partners/tags/:tag_id", c.CheckFuncJWT(c.GetPartnerByTags, "tag-view"))
	v2.Post("/:company/partners/tags/:holder", c.CheckFuncJWT(c.PostPartnerTags, "tag-edit"))

	v2.PostFunc("/:company/outlets", c.PostPartner)
	v2.GetFunc("/:company/outlets", c.GetPartners)
	v2.GetFunc("/:company/outlets/:id", c.GetPartnerByID)
	v2.PutFunc("/:company/outlets/:id", c.UpdatePartner)
	v2.DeleteFunc("/:company/outlets/:id", c.DeletePartner)

	v2.GetFunc("/:company/outlets/tags/:tag_id", c.GetPartnerByTags)
	v2.PostFunc("/:company/outlets/tags/:holder", c.PostPartnerTags)

	// partner/outlet bank
	v2.Post("/:company/banks/:pid", c.CheckFuncJWT(c.PostBank, "bank-create"))
	v2.Get("/:company/banks", c.CheckFuncJWT(c.GetBanks, "bank-view"))
	v2.Get("/:company/banks/:pid", c.CheckFuncJWT(c.GetBankByPartnerID, "bank-view"))
	v2.Put("/:company/banks/:pid", c.CheckFuncJWT(c.UpdateBank, "bank-edit"))
	v2.Delete("/:company/banks/:pid", c.CheckFuncJWT(c.DeleteBank, "bank-delete"))

	//channel
	v2.Post("/:company/channels", c.CheckFuncJWT(c.PostChannel, "channel-create"))
	v2.Get("/:company/channels", c.CheckFuncJWT(c.GetChannels, "channel-view"))
	v2.Get("/:company/channels/:id", c.CheckFuncJWT(c.GetChannelByID, "channel-view"))
	v2.Put("/:company/channels/:id", c.CheckFuncJWT(c.UpdateChannel, "channel-edit"))
	v2.Delete("/:company/channels/:id", c.CheckFuncJWT(c.DeleteChannel, "channel-delete"))

	//users
	// v2.GetFunc("/:company/login", ping)

	//tags
	v2.Post("/:company/tags", c.CheckFuncJWT(c.PostTag, "tag-create"))
	v2.Get("/:company/tags", c.CheckFuncJWT(c.GetTags, "tag-view"))
	v2.Get("/:company/tags/:id", c.CheckFuncJWT(c.GetTagByID, "tag-view"))
	v2.Get("/:company/tags/key/:key", c.CheckFuncJWT(c.GetTagByKey, "tag-view"))
	v2.Get("/:company/tags/category/:category", c.CheckFuncJWT(c.GetTagByKey, "tag-view"))
	// v2.GetFunc("/:company/tags/category/:key", c.GetTagByKey)
	v2.Put("/:company/tags/:id", c.CheckFuncJWT(c.UpdateTag, "tag-edit"))
	v2.Delete("/:company/tags/:id", c.CheckFuncJWT(c.DeleteTag, "tag-delete"))

	v2.Post("/:company/tags/assign/:id", c.CheckFuncJWT(c.PostObjectTags, "tag-edit"))
	v2.Post("/:company/tags/assign", c.CheckFuncJWT(c.PostObjectTags, "tag-edit"))

	//customers
	v2.PostFunc("/:company/customers", c.PostCustomer)
	v2.GetFunc("/:company/customers", c.GetCustomer)
	v2.GetFunc("/:company/customers/:id", c.GetCustomerByID)
	v2.PutFunc("/:company/customers/:id", c.UpdateCustomer)
	v2.DeleteFunc("/:company/customers/:id", c.DeleteCustomer)

	v2.PostFunc("/:company/customers/tags/:id", c.PostCustomerTags)

	//transaction voucher
	v2.Post("/:company/transaction/voucher/assign", c.CheckFuncJWT(c.PostVoucherAssignHolder, "transaction-create"))
	v2.Post("/:company/transaction/voucher/inject/holder", c.CheckFuncJWT(c.PostVoucherInjectByHolder, "transaction-create"))
	v2.Post("/:company/transaction/voucher/claim", c.CheckFuncJWT(c.PostVoucherClaim, "transaction-edit"))
	v2.Post("/:company/transaction/voucher/use", c.CheckFuncJWT(c.PostVoucherUse, "transaction-edit"))
	// v2.PostFunc("/:company/transaction/voucher/redeem", c.PostVoucherRedeem)

	v2.Get("/:company/transaction", c.CheckFuncJWT(c.GetTransactions, "transaction-view"))
	v2.Get("/:company/transaction/outlet/:id", c.CheckFuncJWT(c.GetTransactionsByOutlet, "transaction-view"))
	v2.Get("/:company/transaction/holder/:id", c.CheckFuncJWT(c.GetTransactionsByHolder, "transaction-view"))
	v2.Get("/:company/transaction/program/:id", c.CheckFuncJWT(c.GetTransactionsByProgram, "transaction-view"))
	v2.Get("/:company/transaction/:id", c.CheckFuncJWT(c.GetTransactionByID, "transaction-view"))

	v2.Get("/:company/reimburse/summary", c.CheckFuncJWT(c.GetCashoutSummary, "reimburse-view"))
	v2.Get("/:company/reimburse/list", c.CheckFuncJWT(c.GetCashouts, "reimburse-view"))
	v2.Post("/:company/reimburse/partner", c.CheckFuncJWT(c.PostCashoutByPartner, "reimburse-create"))
	v2.Get("/:company/reimburse/voucher/:program_id", c.CheckFuncJWT(c.GetCashoutUsedVoucher, "reimburse-create"))
	v2.Get("/:company/reimburse/unpaid/", c.CheckFuncJWT(c.GetCashoutsUnpaid, "reimburse-create"))

	// v2.GetFunc("/:company/debug/pprof/", pprof.Index)
	// v2.GetFunc("/:company/debug/pprof/cmdline", pprof.Cmdline)
	// v2.GetFunc("/:company/debug/pprof/profile", pprof.Profile)
	// v2.GetFunc("/:company/debug/pprof/symbol", pprof.Symbol)
	// v2.GetFunc("/:company/debug/pprof/trace", pprof.Trace)

	// blast
	v2.Get("/:company/blasts", c.CheckFuncJWT(c.GetBlasts, "blast-view"))
	v2.Get("/:company/blasts/:id", c.CheckFuncJWT(c.GetBlastByID, "blast-view"))
	v2.Post("/:company/blasts/create", c.CheckFuncJWT(c.CreateEmailBlast, "blast-create"))
	v2.Post("/:company/blasts/send/:id", c.CheckFuncJWT(c.SendEmailBlast, "blast-edit"))
	v2.Get("/:company/blasts/template/:id", c.CheckFuncJWT(c.GetBlastsTemplate, "blast-view"))
	v2.Get("/:company/template/:name", c.CheckFuncJWT(c.GetTemplateByName, "blast-view"))

	// config
	v2.Get("/:company/config", c.CheckFuncJWT(c.GetConfigs, "setting-view"))
	v2.Post("/:company/config/:category", c.CheckFuncJWT(c.SetConfig, "setting-create"))
	v2.Put("/:company/config", c.CheckFuncJWT(c.UpdateConfig, "setting-edit"))

	//public
	v2.Get("/:company/public/voucher", c.CheckFuncJWT(c.GetVoucherByID, "voucher-view"))
	v2.Post("/:company/public/voucher/use", c.CheckFuncJWT(c.PostPublicVoucherUse, "voucher-edit"))

	// GCS
	v2.Post("/:company/file/upload", c.CheckFuncJWT(c.UploadFile, "file-create"))
	v2.Get("/:company/file/delete", c.CheckFuncJWT(c.DeleteFile, "file-delete"))

	v2.Get("/:company/accounts", c.CheckFuncJWT(c.GetAccounts, "member-view"))
	v2.Get("/:company/accounts/:id", c.CheckFuncJWT(c.GetAccountByID, "member-view"))

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
