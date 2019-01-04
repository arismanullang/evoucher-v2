package main

import (
	"net/http"

	c "github.com/gilkor/evoucher/internal/controller"
	"github.com/go-zoo/bone"
)

var router http.Handler

func init() {
	r := bone.New()
	//voucher
	r.PostFunc("/api/v1/vouchers", ping)
	r.GetFunc("/api/v1/vouchers", ping)
	r.GetFunc("/api/v1/vouchers/:id", ping)
	r.PutFunc("/api/v1/vouchers", ping)
	r.DeleteFunc("/api/v1/vouchers", ping)

	//programs
	r.PostFunc("/api/v1/programs", ping)
	r.GetFunc("/api/v1/programs", ping)
	r.GetFunc("/api/v1/programs/:id", ping)
	r.PutFunc("/api/v1/programs", ping)
	r.DeleteFunc("/api/v1/programs", ping)

	//partners
	r.PostFunc("/api/v1/partners", c.PostPartner)
	r.GetFunc("/api/v1/partners", c.GetPartner)
	r.GetFunc("/api/v1/partners/:id", c.GetPartnerByID)
	r.PutFunc("/api/v1/partners", c.UpdatePartner)
	r.DeleteFunc("/api/v1/partners/:id", c.DeleltePartner)

	//tags
	r.PostFunc("/api/v1/tags", c.PostTag)
	r.GetFunc("/api/v1/tags", c.GetTag)
	r.GetFunc("/api/v1/tags/:id", c.GetTagByID)
	r.PutFunc("/api/v1/tags", c.UpdateTag)
	r.DeleteFunc("/api/v1/tags", c.DelelteTag)

	//transactions
	r.PostFunc("/api/v1/transactions", ping)
	r.GetFunc("/api/v1/transactions", ping)
	r.GetFunc("/api/v1/transactions/:id", ping)
	r.PutFunc("/api/v1/transactions", ping)
	r.DeleteFunc("/api/v1/transactions", ping)

	//customers
	r.PostFunc("/api/v1/customers", ping)
	r.GetFunc("/api/v1/customers", ping)
	r.GetFunc("/api/v1/customers/:id", ping)
	r.PutFunc("/api/v1/customers", ping)
	r.DeleteFunc("/api/v1/customers", ping)

	r.GetFunc("/ping", ping)

	// r.GetFunc("/debug/pprof/", pprof.Index)
	// r.GetFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// r.GetFunc("/debug/pprof/profile", pprof.Profile)
	// r.GetFunc("/debug/pprof/symbol", pprof.Symbol)
	// r.GetFunc("/debug/pprof/trace", pprof.Trace)

	router = r
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
