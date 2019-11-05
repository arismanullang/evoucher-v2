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
	v2.GetFunc("/:company/vouchers", ping)
	v2.GetFunc("/:company/vouchers/:id", ping)
	v2.PutFunc("/:company/vouchers/:id", ping)
	v2.DeleteFunc("/:company/vouchers/:id", ping)

	//programs
	v2.PostFunc("/:company/programs", c.PostProgram)
	v2.GetFunc("/:company/programs", c.GetProgram)
	v2.GetFunc("/:company/programs/:id", c.GetProgramByID)
	// v2.PutFunc("/:company/programs/:id", ping)
	v2.DeleteFunc("/:company/programs/:id", c.DeleteProgram)

	//partners == outlets
	v2.PostFunc("/:company/partners", c.PostPartner)
	v2.GetFunc("/:company/partners", c.GetPartners)
	v2.GetFunc("/:company/partners/:id", c.GetPartnerByID)
	v2.PutFunc("/:company/partners/:id", c.UpdatePartner)
	v2.DeleteFunc("/:company/partners/:id", c.DeletePartner)

	v2.GetFunc("/:company/partners/tags/:tag_id", c.GetPartnerByTags)
	v2.PostFunc("/:company/partners/tags/:holder", c.PostPartnerTags)

	v2.PostFunc("/:company/outlets", c.PostPartner)
	v2.GetFunc("/:company/outlets", c.GetPartners)
	v2.GetFunc("/:company/outlets/:id", c.GetPartnerByID)
	v2.PutFunc("/:company/outlets/:id", c.UpdatePartner)
	v2.DeleteFunc("/:company/outlets/:id", c.DeletePartner)

	v2.GetFunc("/:company/outlets/tags/:tag_id", c.GetPartnerByTags)
	v2.PostFunc("/:company/outlets/tags/:holder", c.PostPartnerTags)

	//users
	// v2.GetFunc("/:company/login", ping)

	//tags
	v2.PostFunc("/:company/tags", c.PostTag)
	v2.GetFunc("/:company/tags", c.GetTags)
	v2.GetFunc("/:company/tags/:id", c.GetTagByID)
	v2.PutFunc("/:company/tags/:id", c.UpdateTag)
	v2.DeleteFunc("/:company/tags/:id", c.DeleteTag)

	//customers
	v2.PostFunc("/:company/customers", c.PostCustomer)
	v2.GetFunc("/:company/customers", c.GetCustomer)
	v2.GetFunc("/:company/customers/:id", c.GetCustomerByID)
	v2.PutFunc("/:company/customers/:id", c.UpdateCustomer)
	v2.DeleteFunc("/:company/customers/:id", c.DeleteCustomer)

	v2.PostFunc("/:company/customers/tags/:id", c.PostCustomerTags)

	//transaction voucher
	v2.PostFunc("/:company/transaction/voucher/assign", c.PostVoucherAssignHolder)
	v2.PostFunc("/:company/transaction/voucher/assignholder", c.PostVoucherAssignHolder)
	// v2.PostFunc("/:company/transaction/voucher/redeem", c.PostVoucherRedeem)

	// v2.GetFunc("/:company/debug/pprof/", pprof.Index)
	// v2.GetFunc("/:company/debug/pprof/cmdline", pprof.Cmdline)
	// v2.GetFunc("/:company/debug/pprof/profile", pprof.Profile)
	// v2.GetFunc("/:company/debug/pprof/symbol", pprof.Symbol)
	// v2.GetFunc("/:company/debug/pprof/trace", pprof.Trace)
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
