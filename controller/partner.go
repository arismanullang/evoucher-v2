package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
)

//PostPartner : POST partner data
func PostPartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	token := r.FormValue("token")

	accData, err := model.GetSessionDataJWT(token)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqPartner.CreatedBy = accData.AccountID
	reqPartner.UpdatedBy = accData.AccountID
	response, err := reqPartner.Insert()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusCreated)
}

type PartnerFilter struct {
	ID          string `schema:"id" filter:"array"`
	Name        string `schema:"name" filter:"string"`
	Description string `schema:"description" filter:"record"`
	CompanyID   string `schema:"company_id" filter:"string"`
	CreatedAt   string `schema:"created_at" filter:"date"`
	CreatedBy   string `schema:"created_by" filter:"string"`
	UpdatedAt   string `schema:"updated_at" filter:"date"`
	UpdatedBy   string `schema:"updated_by" filter:"string"`
	Status      string `schema:"status" filter:"enum"`
	Banks       string `schema:"partner_banks" filter:"json"`
	Tags        string `schema:"partner_tags" filter:"json"`
}

//GetPartners : GET list of partners
func GetPartners(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	qp.SetCompanyID(bone.GetValue(r, "company"))

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	var f PartnerFilter
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	partners, next, err := model.GetPartners(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(partners)
	res.SetNewPagination(r, qp.Page, next, (*partners)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//GetPartnerByID : GET
func GetPartnerByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	partner, _, err := model.GetPartnerByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(model.Partners{*partner})
	res.JSON(w, res, http.StatusOK)
}

// UpdatePartner :
func UpdatePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	token := r.FormValue("token")

	accData, err := model.GetSessionDataJWT(token)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "id")
	var reqPartner model.Partner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPartner); err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	reqPartner.ID = id
	reqPartner.UpdatedBy = accData.AccountID
	err = reqPartner.Update()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.SetResponse(model.Partners{reqPartner})
	res.JSON(w, res, http.StatusOK)
}

//DeletePartner : remove partner
func DeletePartner(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	accData, err := model.GetSessionDataJWT(token)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	id := bone.GetValue(r, "id")
	p := model.Partner{ID: id}
	p.UpdatedBy = accData.AccountID
	if err := p.Delete(); err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	res.JSON(w, res, http.StatusOK)
}

//PostPartnerTags : POST tags of partner
func PostPartnerTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	accData, err := model.GetSessionDataJWT(token)
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
	// reqPartner.ID = bone.GetValue(r, "holder")
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

//GetPartnerByTags : GET
func GetPartnerByTags(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "tag_id")
	partners, next, err := model.GetPartnersByTags(qp, id)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(partners)
	res.SetNewPagination(r, qp.Page, next, (*partners)[0].Count)
	res.JSON(w, res, http.StatusOK)
}

//PostBank : POST bank data
func PostBank(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	token := r.FormValue("token")

	accData, err := model.GetSessionDataJWT(token)
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
	reqBank.PartnerID = bone.GetValue(r, "pid")
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

type BankFilter struct {
	ID              string `schema:"id" filter:"array"`
	PartnerID       string `schema:"partner_id" filter:"array"`
	BankName        string `schema:"bank_name" filter:"string"`
	BankBranch      string `schema:"bank_branch" filter:"string"`
	BankAccountName string `schema:"bank_account_name" filter:"string"`
	CompanyName     string `schema:"company_name" filter:"string"`
	Name            string `schema:"name" filter:"string"`
	Phone           string `schema:"phone" filter:"string"`
	Email           string `schema:"email" filter:"string"`
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

//GetBankByPartnerID : GET bank by partner id
func GetBankByPartnerID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	partnerID := bone.GetValue(r, "pid")
	bank, _, err := model.GetBankByPartnerID(qp, partnerID)
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

	accData, err := model.GetSessionDataJWT(token)
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
	reqBank.PartnerID = id
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

	accData, err := model.GetSessionDataJWT(token)
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

func GetPartnerBanks(r *http.Request, partnerID string) ([]model.Bank, error) {

	qp := u.NewQueryParam(r)
	partner, _, err := model.GetPartnerByID(qp, partnerID)
	if err != nil {
		return []model.Bank{}, err
	}

	partnerBank := []model.Bank{}
	err = json.Unmarshal([]byte(partner.Banks), &partnerBank)
	if err != nil {
		return []model.Bank{}, err
	}

	if len(partnerBank) < 1 {
		return []model.Bank{}, model.ErrorBankNotFound
	}

	return partnerBank, nil
}
