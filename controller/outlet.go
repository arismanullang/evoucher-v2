package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
)

type (
	//OutletFilter : an interface that used as an outlet filter
	OutletFilter struct {
		ID          string `schema:"id" filter:"array"`
		Name        string `schema:"name" filter:"string"`
		Description string `schema:"description" filter:"record"`
		CompanyID   string `schema:"company_id" filter:"string"`
		CreatedAt   string `schema:"created_at" filter:"date"`
		CreatedBy   string `schema:"created_by" filter:"string"`
		UpdatedAt   string `schema:"updated_at" filter:"date"`
		UpdatedBy   string `schema:"updated_by" filter:"string"`
		Status      string `schema:"status" filter:"enum"`
		Banks       string `schema:"outlet_banks" filter:"json"`
		Tags        string `schema:"outlet_tags" filter:"json"`
	}

	//BankFilter : an interface that used as an bank filter
	BankFilter struct {
		ID              string `schema:"id" filter:"array"`
		OutletID        string `schema:"outlet_id" filter:"array"`
		BankName        string `schema:"bank_name" filter:"string"`
		BankBranch      string `schema:"bank_branch" filter:"string"`
		BankAccountName string `schema:"bank_account_name" filter:"string"`
		CompanyName     string `schema:"company_name" filter:"string"`
		Name            string `schema:"name" filter:"string"`
		Phone           string `schema:"phone" filter:"string"`
		Email           string `schema:"email" filter:"string"`
	}
)

//PostOutlet : POST outlet data
func PostOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var reqOutlet model.Outlet
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqOutlet); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqOutlet.CreatedBy = accData.AccountID
	reqOutlet.UpdatedBy = accData.AccountID
	response, err := reqOutlet.Insert()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetOutlets : GET list of outlets
func GetOutlets(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(bone.GetValue(r, "company"))

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f OutletFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	outlets, next, err := model.GetOutlets(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(outlets)
	if len(*outlets) > 0 {
		res.SetNewPagination(r, qp.Page, next, (*outlets)[0].Count)
	}
	res.JSON(w, res, http.StatusOK)
}

//GetOutletByID : GET
func GetOutletByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	outlet, _, err := model.GetOutletByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(model.Outlets{*outlet})
	res.JSON(w, res, http.StatusOK)
}

// UpdateOutlet :
func UpdateOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "id")
	var reqOutlet model.Outlet
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqOutlet); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqOutlet.ID = id
	reqOutlet.UpdatedBy = accData.AccountID
	err = reqOutlet.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Outlets{reqOutlet})
	res.JSON(w, res, http.StatusOK)
}

//DeleteOutlet : remove outlet
func DeleteOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "id")
	p := model.Outlet{ID: id}
	p.UpdatedBy = accData.AccountID
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//PostOutletTags : POST tags of outlet
func PostOutletTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var req model.ObjectTag
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	// reqOutlet.ID = bone.GetValue(r, "holder")
	req.CreatedBy = accData.AccountID
	response, err := req.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetOutletByTags : GET
func GetOutletByTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "tag_id")
	outlets, next, err := model.GetOutletsByTags(qp, id)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(outlets)
	res.SetNewPagination(r, qp.Page, next, (*outlets)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//PostBank : POST bank data
func PostBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var reqBank model.Bank
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBank); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqBank.OutletID = bone.GetValue(r, "pid")
	reqBank.CreatedBy = accData.AccountID
	reqBank.UpdatedBy = accData.AccountID
	response, err := reqBank.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

//GetBanks : GET list of banks
func GetBanks(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f BankFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	banks, next, err := model.GetBanks(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(banks)
	res.SetNewPagination(r, qp.Page, next, (*banks)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetBankByOutletID : GET bank by outlet id
func GetBankByOutletID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	outletID := bone.GetValue(r, "pid")
	bank, _, err := model.GetBankByOutletID(qp, outletID)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(bank)
	res.JSON(w, res, http.StatusOK)
}

// UpdateBank :
func UpdateBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "pid")
	var reqBank model.Bank
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBank); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqBank.OutletID = id
	reqBank.UpdatedBy = accData.AccountID
	err = reqBank.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Banks{reqBank})
	res.JSON(w, res, http.StatusOK)
}

//DeleteBank : remove bank
func DeleteBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	companyID := bone.GetValue(r, "company")

	accData, err := model.GetSessionDataJWT(token, companyID)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}
	id := bone.GetValue(r, "id")
	p := model.Bank{ID: u.StringToInt(id)}
	p.UpdatedBy = accData.AccountID
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//GetOutletBanks : get outlet banks
func GetOutletBanks(r *http.Request, outletID string) ([]model.Bank, error) {

	qp := u.NewQueryParam(r)
	outlet, _, err := model.GetOutletByID(qp, outletID)
	if err != nil {
		return []model.Bank{}, err
	}

	outletBank := []model.Bank{}
	err = json.Unmarshal([]byte(outlet.Banks), &outletBank)
	if err != nil {
		fmt.Println("err = ", err)
		return []model.Bank{}, err
	}

	if len(outletBank) < 1 {
		return []model.Bank{}, model.ErrorBankNotFound
	}

	return outletBank, nil
}
