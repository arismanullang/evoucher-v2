package main

import (
	"fmt"
	"net/http"

	c "github.com/gilkor/evoucher/internal/controller"
	"github.com/go-zoo/bone"
)

var router http.Handler

func init() {
	//main router
	r := bone.New()
	// r.NotFoundFunc(notFound)
	r.GetFunc("/ping", ping)

	//define sub router
	v1 := bone.New()
	r.SubRoute("/api/v1.0", v1)

	//voucher
	v1.PostFunc("/:company/vouchers", ping)
	v1.GetFunc("/:company/vouchers", ping)
	v1.GetFunc("/:company/vouchers/:id", ping)
	v1.PutFunc("/:company/vouchers", ping)
	v1.DeleteFunc("/:company/vouchers", ping)

	//programs
	v1.PostFunc("/:company/programs", ping)
	v1.GetFunc("/:company/programs", ping)
	v1.GetFunc("/:company/programs/:id", ping)
	v1.PutFunc("/:company/programs", ping)
	v1.DeleteFunc("/:company/programs", ping)

	//partners
	v1.PostFunc("/:company/partners", c.PostPartner)
	v1.GetFunc("/:company/partners", c.GetPartner)
	v1.GetFunc("/:company/partners/:id", c.GetPartnerByID)
	v1.PutFunc("/:company/partners", c.UpdatePartner)
	v1.DeleteFunc("/:company/partners/:id", c.DeletePartner)

	//tags
	v1.PostFunc("/:company/tags", c.PostTag)
	v1.GetFunc("/:company/tags", c.GetTag)
	v1.GetFunc("/:company/tags/:id", c.GetTagByID)
	v1.PutFunc("/:company/tags", c.UpdateTag)
	v1.DeleteFunc("/:company/tags", c.DeleteTag)

	//users
	v1.GetFunc("/:company/login", ping)

	//Roles
	v1.PostFunc("/:company/tags", c.PostTag)
	v1.GetFunc("/:company/tags", c.GetTag)
	v1.GetFunc("/:company/tags/:id", c.GetTagByID)
	v1.PutFunc("/:company/tags", c.UpdateTag)
	v1.DeleteFunc("/:company/tags", c.DeleteTag)

	//customers
	v1.PostFunc("/:company/customers", c.PostCustomer)
	v1.GetFunc("/:company/customers", c.GetCustomer)
	v1.GetFunc("/:company/customers/:id", c.GetCustomerByID)
	v1.PutFunc("/:company/customers", c.UpdateCustomer)
	v1.DeleteFunc("/:company/customers", c.DeleteCustomer)

	// v1.GetFunc("/:company/debug/pprof/", pprof.Index)
	// v1.GetFunc("/:company/debug/pprof/cmdline", pprof.Cmdline)
	// v1.GetFunc("/:company/debug/pprof/profile", pprof.Profile)
	// v1.GetFunc("/:company/debug/pprof/symbol", pprof.Symbol)
	// v1.GetFunc("/:company/debug/pprof/trace", pprof.Trace)

	router = r
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ping")
	w.Write([]byte("pong"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not found"))
}
