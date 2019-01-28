package main

import (
	"net/http"

	c "github.com/gilkor/evoucher/internal/controller"
	"github.com/go-zoo/bone"
)

var router http.Handler

func init() {
	r := bone.New()
	//use company prefix

	//define router
	r.Prefix("/api/v1.0/:company")
	// r.NotFoundFunc(notFound)

	//voucher
	r.PostFunc("/vouchers", ping)
	r.GetFunc("/vouchers", ping)
	r.GetFunc("/vouchers/:id", ping)
	r.PutFunc("/vouchers", ping)
	r.DeleteFunc("/vouchers", ping)

	//programs
	r.PostFunc("/programs", ping)
	r.GetFunc("/programs", ping)
	r.GetFunc("/programs/:id", ping)
	r.PutFunc("/programs", ping)
	r.DeleteFunc("/programs", ping)

	//partners
	r.PostFunc("/partners", c.PostPartner)
	r.GetFunc("/partners", c.GetPartner)
	r.GetFunc("/partners/:id", c.GetPartnerByID)
	r.PutFunc("/partners", c.UpdatePartner)
	r.DeleteFunc("/partners/:id", c.DeletePartner)

	//tags
	r.PostFunc("/tags", c.PostTag)
	r.GetFunc("/tags", c.GetTag)
	r.GetFunc("/tags/:id", c.GetTagByID)
	r.PutFunc("/tags", c.UpdateTag)
	r.DeleteFunc("/tags", c.DeleteTag)

	//users
	r.GetFunc("/login", ping)

	//Roles
	r.PostFunc("/tags", c.PostTag)
	r.GetFunc("/tags", c.GetTag)
	r.GetFunc("/tags/:id", c.GetTagByID)
	r.PutFunc("/tags", c.UpdateTag)
	r.DeleteFunc("/tags", c.DeleteTag)

	//customers
	r.PostFunc("/customers", c.PostCustomer)
	r.GetFunc("/customers", c.GetCustomer)
	r.GetFunc("/customers/:id", c.GetCustomerByID)
	r.PutFunc("/customers", c.UpdateCustomer)
	r.DeleteFunc("/customers", c.DeleteCustomer)

	// r.GetFunc("/ping", ping)

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

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not found"))
}
