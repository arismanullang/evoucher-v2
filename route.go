package main

import (
	"fmt"
	"net/http"

	c "github.com/gilkor/evoucher/controller"
	"github.com/go-zoo/bone"
)

var router http.Handler

func init() {
	//main router
	r := bone.New()
	// r.NotFoundFunc(notFound)
	r.GetFunc("/ping", ping)

	//define sub router
	v2 := bone.New()
	r.SubRoute("/api/v2.0", v2)

	//voucher
	v2.PostFunc("/:company/vouchers", ping)
	v2.GetFunc("/:company/vouchers", ping)
	v2.GetFunc("/:company/vouchers/:id", ping)
	v2.PutFunc("/:company/vouchers", ping)
	v2.DeleteFunc("/:company/vouchers", ping)

	//programs
	v2.PostFunc("/:company/programs", c.PostProgram)
	v2.GetFunc("/:company/programs", c.GetProgram)
	v2.GetFunc("/:company/programs/:id", c.GetProgramByID)
	v2.PutFunc("/:company/programs", ping)
	v2.DeleteFunc("/:company/programs/:id", c.DeleteProgram)

	//partners
	v2.PostFunc("/:company/partners", c.PostPartner)
	v2.GetFunc("/:company/partners", c.GetPartners)
	v2.GetFunc("/:company/partners/:id", c.GetPartnerByID)
	v2.PutFunc("/:company/partners", c.UpdatePartner)
	v2.DeleteFunc("/:company/partners/:id", c.DeletePartner)

	//tags
	v2.PostFunc("/:company/tags", c.PostTag)
	v2.GetFunc("/:company/tags", c.GetTags)
	v2.GetFunc("/:company/tags/:id", c.GetTagByID)
	v2.PutFunc("/:company/tags", c.UpdateTag)
	v2.DeleteFunc("/:company/tags", c.DeleteTag)

	//users
	v2.GetFunc("/:company/login", ping)

	//Roles
	v2.PostFunc("/:company/tags", c.PostTag)
	v2.GetFunc("/:company/tags", c.GetTags)
	v2.GetFunc("/:company/tags/:id", c.GetTagByID)
	v2.PutFunc("/:company/tags", c.UpdateTag)
	v2.DeleteFunc("/:company/tags", c.DeleteTag)

	//customers
	v2.PostFunc("/:company/customers", c.PostCustomer)
	v2.GetFunc("/:company/customers", c.GetCustomer)
	v2.GetFunc("/:company/customers/:id", c.GetCustomerByID)
	v2.PutFunc("/:company/customers", c.UpdateCustomer)
	v2.DeleteFunc("/:company/customers", c.DeleteCustomer)
	v2.PostFunc("/:company/customers/tags/:id", c.PostCustomerTags)

	// v2.GetFunc("/:company/debug/pprof/", pprof.Index)
	// v2.GetFunc("/:company/debug/pprof/cmdline", pprof.Cmdline)
	// v2.GetFunc("/:company/debug/pprof/profile", pprof.Profile)
	// v2.GetFunc("/:company/debug/pprof/symbol", pprof.Symbol)
	// v2.GetFunc("/:company/debug/pprof/trace", pprof.Trace)

	router = r
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ping")
	w.Write([]byte("pong"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not found"))
}
