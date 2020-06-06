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
func GetTransactionsByOutlet(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	result, next, err := model.GetTransactionByPartner(qp, id)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(result)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetTransactions : GET list of partners
func GetTransactionsByProgram(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")
	result, next, err := model.GetTransactionByProgram(qp, id)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(result)
	res.SetPagination(r, qp.Page, next)
	res.JSON(w, res, http.StatusOK)
}

//GetTransactions : GET list of partners
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)

	result, next, err := model.GetTransactions(qp)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(result)
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
// func PostVoucherAssignHolder(w http.ResponseWriter, r *http.Request) {
// 	res := u.NewResponse()

// 	var req model.Voucher
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&req)
// 	req.State = model.VoucherStateUsed
// 	if err != nil {
// 		u.DEBUG(err)
// 		res.SetError(JSONErrBadRequest)
// 		res.JSON(w, res, JSONErrBadRequest.Status)
// 		return
// 	}
// 	if err := req.Update(); err != nil {
// 		u.DEBUG(err)
// 		res.SetErrorWithDetail(JSONErrFatal, err)
// 		res.JSON(w, res, JSONErrFatal.Status)
// 		return
// 	}
// 	res.JSON(w, res, http.StatusCreated)
// }

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
	companyID := bone.GetValue(r, "company")
	accountToken := r.FormValue("token")
	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	//get config TimeZone
	configs, err := model.GetConfigs(companyID, "company")
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.Error.SetMessage("timezone config not found, please add timezone config")
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	//getUser
	account, err := model.GetAccountByID(qp, claims.AccountID)
	if err != nil {
		u.DEBUG(err, " User Not Found")
		res.SetError(JSONErrForbidden)
		res.JSON(w, res, JSONErrForbidden.Status)
		return
	}

	//Validate VoucherID
	voucher, err := model.GetVoucherByID(req.VoucherID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	} else if voucher.State == model.VoucherStateUsed {
		res.SetError(JSONErrBadRequest)
		res.Error.Message = "Voucher has been used"
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	} else if voucher.State == model.VoucherStatePaid {
		res.SetError(JSONErrBadRequest)
		res.Error.Message = "Voucher has been paid"
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	} else if !voucher.ExpiredAt.After(time.Now()) {
		res.SetError(JSONErrBadRequest)
		res.Error.Message = "Voucher has expired at " + voucher.ExpiredAt.String()
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	} else if !voucher.ValidAt.Before(time.Now()) {
		res.SetError(JSONErrBadRequest)
		res.Error.Message = "Voucher can be used at " + voucher.ValidAt.String()
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	datas := make(map[string]string)
	datas["ACCOUNTID"] = account.ID
	datas["PROGRAMID"] = voucher.ProgramID
	datas["OUTLETID"] = req.OutletID
	datas["TIMEZONE"] = fmt.Sprint(configs["timezone"])

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
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	if !result {
		u.DEBUG(err)
		res.SetError(JSONErrInvalidRule)
		res.JSON(w, res, JSONErrInvalidRule.Status)
		return
	}

	var tx model.Transaction
	var td model.TransactionDetail

	td.ProgramId = program.ID
	td.VoucherId = voucher.ID
	td.CreatedBy = account.ID
	td.UpdatedBy = account.ID

	tx.TransactionDetails = append(tx.TransactionDetails, td)

	tx.CompanyId = companyID
	tx.TransactionCode = u.RandomizeString(u.TRANSACTION_CODE_LENGTH, u.NUMERALS)
	// Multiply with total voucher used
	tx.TotalAmount = fmt.Sprint(program.MaxValue)
	tx.Holder = account.ID
	tx.PartnerId = req.OutletID
	tx.CreatedBy = account.ID
	tx.UpdatedBy = account.ID

	if _, err = tx.Insert(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	voucher.State = model.VoucherStateUsed
	voucher.UpdatedBy = account.ID
	if err = voucher.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	//send email confirmation

	res.SetResponse(tx)
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

	var req model.VoucherClaimRequest
	var accountID string

	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)

	companyID := bone.GetValue(r, "company")
	accountToken := r.FormValue("token")

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	accountID = claims.AccountID

	//get config TimeZone
	configs, err := model.GetConfigs(companyID, "company")
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.Error.SetMessage("timezone config not found, please add timezone config")
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	//Get Holder Detail
	account, err := model.GetAccountByID(qp, accountID)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

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

	var rules model.RulesExpression
	program.Rule.Unmarshal(&rules)

	datas := make(map[string]interface{})
	datas["ACCOUNTID"] = accountID
	datas["PROGRAMID"] = req.ProgramID
	datas["QUANTITY"] = req.Quantity
	datas["TIMEZONE"] = fmt.Sprint(configs["timezone"])

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
		res.JSON(w, err, JSONErrInvalidRule.Status)
		return
	}

	currentClaimedVoucher, err := model.GetVoucherCreatedAmountByProgramID(req.ProgramID)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}
	if int64(currentClaimedVoucher+req.Quantity) > program.Stock {
		u.DEBUG(err)
		res.SetError(JSONErrExceedAmount)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	//Create Voucher
	gvr := GenerateVoucherRequest{
		ProgramID:    req.ProgramID,
		Quantity:     req.Quantity,
		ReferenceNo:  req.Reference,
		HolderID:     accountID,
		HolderDetail: holderDetail,
		UpdatedBy:    accountID,
	}

	vouchers, err := gvr.GenerateVoucher(fmt.Sprint(configs["timezone"]), *program)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	response, err := vouchers.Insert()
	if err != nil {
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)

	res.JSON(w, res, http.StatusCreated)
}
