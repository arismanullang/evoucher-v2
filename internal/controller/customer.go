package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	u "github.com/gilkor/evoucher/internal/util"
	"github.com/go-zoo/bone"
	"github.com/ruizu/render"
)

//PostCustomer : POST Customer data
func PostCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqCustomer model.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqCustomer); err != nil {
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}
	if err := reqCustomer.Insert(); err != nil {
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}

	render.JSON(w, res, http.StatusCreated)
}

//GetCustomer : GET list of Customers
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	f := u.NewFilter(r)
	customers, next, err := model.GetCustomers(f)
	if err != nil {
		res.SetError(ErrFatal.SetArgs(err.Error()))
		render.JSON(w, res, ErrFatal.Status)
		return
	}

	res.SetResponse(customers)
	res.SetPagination(r, f.Page, next)
	render.JSON(w, res, http.StatusOK)
}

//GetCustomerByID : GET
func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	f := u.NewFilter(r)
	id := bone.GetValue(r, "id")
	customer, _, err := model.GetCustomerByID(id, f)
	if err != nil {
		res.SetError(ErrResourceNotFound)
		render.JSON(w, res, ErrResourceNotFound.Status)
		return
	}

	res.SetResponse(customer)
	render.JSON(w, res, http.StatusOK)
}

// UpdateCustomer :
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqCustomer model.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqCustomer); err != nil {
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}
	if err := reqCustomer.Update(); err != nil {
		res.SetError(ErrFatal)
		render.JSON(w, res, ErrFatal.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}

//DeleteCustomer : remove Customer
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Customer{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(ErrResourceNotFound)
		render.JSON(w, res, ErrResourceNotFound.Status)
		return
	}
	render.JSON(w, res, http.StatusCreated)
}