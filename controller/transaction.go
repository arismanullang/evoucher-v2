package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	// "github.com/go-zoo/bone"
)

type (
	//VoucherClaimRequest : body struct of claim voucher request
	VoucherClaimRequest struct {
		Reference string `json:'reference'`
		ProgramID string `json:"program_id"`
		Quantity  int    `json:"quantity"`
	}

	//VoucherUseRequest : body struct of use voucher request
	VoucherUseRequest struct {
		Reference    string                 `json:'reference'`
		Transactions VoucherUseTransactions `json:"transactions"`
		Vouchers     []string               `json:"vouchers"`
		OutletID     string                 `json:"outlet_id"`
	}

	//VoucherUseTransactions :
	VoucherUseTransactions struct {
		TotalAmount float64 `json:"total_amount"`
		Details     string  `json:"details"`
	}
)

//GetTransactions : GET list of partners
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	partners, next, err := model.GetTransactions(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(partners)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetTransactionByID : GET
func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	partner, _, err := model.GetTransactionByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	res.SetResponse(partner)
	res.JSON(w, res, http.StatusOK)
}

//PostVoucherAssignHolder :
func PostVoucherAssignHolder(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateUsed
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if err := req.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

type UseTransaction struct {
	OutletID    string            `db:"outlet_id" json:"outlet_id,omitempty"`
	VoucherID   string            `db:"voucher_id" json:"voucher_id,omitempty"`
	Transaction model.Transaction `json:"transaction,omitempty"`
}

//PostVoucherUse :
func PostVoucherUse(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req UseTransaction
	var rule model.Rules
	var accountID string
	accountID = r.FormValue("xx-token")
	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)

	voucher, err := model.GetVoucherByID(req.VoucherID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	datas := make(map[string]string)
	datas["ACCOUNTID"] = *voucher.Holder
	datas["PROGRAMID"] = voucher.ProgramID

	//Validate Rule Program
	program, err := model.GetProgramByID(voucher.ProgramID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	err = rule.Unmarshal(program.Rule)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	var rules model.RulesExpression
	program.Rule.Unmarshal(&rules)

	result, err := rules.ValidateUse(datas)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	if !result {
		u.DEBUG(err)
		res.SetError(JSONErrInvalidRule)
		res.JSON(w, res, JSONErrInvalidRule.Status)
		return
	}

	voucher.State = model.VoucherStateClaim
	voucher.Holder = &accountID
	if err = voucher.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

//PostVoucherUset :
func PostVoucherUset(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	var rule model.Rules

	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateUsed

	//if key == "" {
	//	res.SetError(JSONErrUnauthorized)
	//	res.JSON(w, res, JSONErrUnauthorized.Status)
	//	return
	//}
	//
	//token, err := VerifyJWT(key)
	//if err != nil {
	//	res.SetError(JSONErrUnauthorized)
	//	res.JSON(w, res, JSONErrUnauthorized.Status)
	//	return
	//}
	//
	//claims, ok := token.Claims.(*JWTJunoClaims)
	//if ok && token.Valid {
	//	// fmt.Printf("Key:%v", token.Header)
	//} else {
	//	res.SetError(JSONErrUnauthorized)
	//	res.JSON(w, res, JSONErrUnauthorized.Status)
	//	return
	//}

	//Validate Rule Program
	program, err := model.GetProgramByID(req.ProgramID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	err = rule.Unmarshal(program.Rule)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	result, err := rule.Validate()
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if !result {
		//expected error
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	//parse data transaction
	//validate user data
	//get program id, user data
	//check voucher
	//
	//

	if err = req.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	res.JSON(w, res, http.StatusCreated)
}

//PostVoucherClaim :
func PostVoucherClaim(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req VoucherClaimRequest
	var rule model.Rules
	var accountID string

	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)

	token, err := VerifyJWT(r.FormValue("token"))
	if err != nil {
		if err == model.ErrorForbidden {
			u.NewResponse().SetError(JSONErrForbidden)
			return
		} else if err == model.ErrorInternalServer {
			u.NewResponse().SetError(JSONErrFatal)
			return
		}
	}

	claims, ok := token.Claims.(*JWTJunoClaims)
	if ok {
		accountID = claims.AccountID
	} else {
		res.SetError(JSONErrForbidden)
		res.JSON(w, res, JSONErrForbidden.Status)
		return
	}

	datas := make(map[string]string)
	datas["ACCOUNTID"] = accountID
	datas["PROGRAMID"] = req.ProgramID

	//Get Holder Detail
	accounts, _, err := model.GetAccountByID(qp, accountID)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}
	account := &(*accounts)[0]
	tmp := model.HolderDetail{
		Name:  account.Name,
		Phone: account.MobileNo,
		Email: account.Email,
	}

	holderDetail, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Validate Rule Program
	program, err := model.GetProgramByID(req.ProgramID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	err = rule.Unmarshal(program.Rule)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	var rules model.RulesExpression
	program.Rule.Unmarshal(&rules)
	voucherValidAt := time.Now()
	voucherExpiredAt := time.Now()
	fmt.Println("rules = ", rules)

	if ruleUseUsagePeriod, ok := rules.And["rule_use_usage_period"]; ok {

		fmt.Println("ruleUseActiveVoucherPeriod = ", ruleUseUsagePeriod)
		validTime, err := model.StringToTime(fmt.Sprint(ruleUseUsagePeriod.Gte))
		if err != nil {
			res.SetError(JSONErrBadRequest)
			res.Error.SetMessage("failed to parse active voucher period")
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		expiredTime, err := model.StringToTime(fmt.Sprint(ruleUseUsagePeriod.Lte))
		if err != nil {
			res.SetError(JSONErrBadRequest)
			res.Error.SetMessage("failed to parse active voucher period")
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		voucherValidAt = validTime
		voucherExpiredAt = expiredTime
		fmt.Println("voucherValidAt = ", voucherValidAt)
		fmt.Println("voucherExpiredAt = ", voucherExpiredAt)
	}

	if ruleUseActiveVoucherPeriod, ok := rules.And["rule_use_active_voucher_period"]; ok && !ruleUseActiveVoucherPeriod.IsEmpty() {
		voucherExpiredAt = voucherValidAt.AddDate(0, 0, int(ruleUseActiveVoucherPeriod.Eq.(float64)))
		fmt.Println("voucherValidAt = ", voucherValidAt)
		fmt.Println("voucherExpiredAt = ", voucherExpiredAt)
	}

	result, err := rules.ValidateClaim(datas)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	if !result {
		u.DEBUG(err)
		res.SetError(JSONErrInvalidRule)
		res.JSON(w, res, JSONErrInvalidRule.Status)
		return
	}

	//Checking Amount : get current total voucher by programID
	currentVoucherAmount, err := model.GetVoucherCreatedAmountByProgramID(program.ID)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if int64(currentVoucherAmount+req.Quantity) > program.Stock {
		res.SetError(JSONErrExceedAmount)
		res.JSON(w, res, http.StatusOK)
		return
	}

	//Create Voucher
	var vouchers model.Vouchers
	var vf model.VoucherFormat
	program.VoucherFormat.Unmarshal(&vf)

	for i := 0; i < req.Quantity; i++ {
		voucher := new(model.Voucher)
		if vf.Type == "fix" {
			voucher.Code = vf.Code
		} else if vf.Type == "random" {
			voucher.Code = vf.Prefix + u.RandomizeString(u.LENGTH, vf.Random) + vf.Postfix
		}

		voucher.ReferenceNo = req.Reference
		voucher.Holder = &accountID
		voucher.HolderDetail = holderDetail
		voucher.ProgramID = program.ID
		voucher.CreatedBy = "system"
		voucher.UpdatedBy = "system"
		voucher.Status = model.StatusCreated
		voucher.State = model.VoucherStateCreated
		voucher.ValidAt = &voucherValidAt
		voucher.ExpiredAt = &voucherExpiredAt

		vouchers = append(vouchers, *voucher)
	}

	res.SetResponse(vouchers)

	// if _, err := vouchers.Insert(); err != nil {
	// 	fmt.Println(err)
	// 	res.SetError(JSONErrFatal.SetArgs(err.Error()))
	// 	res.JSON(w, res, JSONErrFatal.Status)
	// 	return
	// }

	res.JSON(w, res, http.StatusCreated)
}
