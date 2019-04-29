package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gilkor/evoucher/model"
	u "github.com/gilkor/evoucher/util"
	"github.com/go-zoo/bone"
)

//PostCustomer : POST Customer data
func PostCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqCustomer model.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqCustomer); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqCustomer.Insert(); err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.JSON(w, res, http.StatusCreated)
}

//GetCustomer : GET list of Customers
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	customers, next, err := model.GetCustomers(qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(customers)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetCustomerByID : GET
func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	customer, _, err := model.GetCustomerByID(id, qp)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(customer)
	res.JSON(w, res, http.StatusOK)
}

// UpdateCustomer :
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var reqCustomer model.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqCustomer); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if err := reqCustomer.Update(); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

//DeleteCustomer : remove Customer
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	p := model.Customer{ID: id}
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

//Customer Tag

//PostCustomerTags :
func PostCustomerTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	id := bone.GetValue(r, "id")
	model := model.TagHolder{Holder: id}
	if err := model.Insert(); err != nil {
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}
