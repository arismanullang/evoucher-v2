package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	// "github.com/go-zoo/bone"
)

type (
	VoucherUseRequest struct {
		Reference    string                 `json:'reference'`
		Transactions VoucherUseTransactions `json:"transactions"`
		Vouchers     []string               `json:"vouchers"`
		OutletID     string                 `json:"outlet_id"`
	}
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
	accountToken := r.FormValue("xx-token")
	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)

	token, err := VerifyJWT(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	claims, ok := token.Claims.(*JWTJunoClaims)
	if ok && token.Valid {
		// fmt.Printf("Key:%v", token.Header)
	} else {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
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

	voucher, err := model.GetVoucherByID(req.VoucherID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	datas := make(map[string]string)
	datas["ACCOUNTID"] = account.ID
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

	//Transaction struct {
	//	ID              string     `db:"id",json:"id"`
	//	CompanyId       string     `db:"company_id",json:"company_id"`
	//	TransactionCode string     `db:"transaction_code",json:"transaction_code"`
	//	TotalAmount     string     `db:"total_amount",json:"total_amount"`
	//	Holder          string     `db:"holder",json:"holder"`
	//	PartnerId       string     `db:"partner_id",json:"partner_id"`
	//	CreatedBy       string     `db:"created_by",json:"created_by"`
	//	CreatedAt       *time.Time `db:"created_at",json:"created_at"`
	//	UpdatedBy       string     `db:"updated_by",json:"updated_by"`
	//	UpdatedAt       *time.Time `db:"updated_at",json:"updated_at"`
	//	Status          string     `db:"status",json:"status"`
	//}
	//Transactions      []Transaction
	//TransactionDetail struct {
	//	ID            string     `db:"id",json:"id"`
	//	TransactionId string     `db:"transaction_id",json:"transaction_id"`
	//	ProgramId     string     `db:"program_id",json:"program_id"`
	//	VoucherId     string     `db:"voucher_id",json:"voucher_id"`
	//	CreatedBy     string     `db:"created_by",json:"created_by"`
	//	CreatedAt     *time.Time `db:"created_at",json:"created_at"`
	//	UpdatedBy     string     `db:"updated_by",json:"updated_by"`
	//	UpdatedAt     *time.Time `db:"updated_at",json:"updated_at"`
	//	Status        string     `db:"status",json:"status"`
	//}

	var tx model.Transaction
	var td model.TransactionDetail

	td.ProgramId = program.ID
	td.VoucherId = voucher.ID
	td.CreatedBy = "system"
	td.UpdatedBy = "system"

	tx.TransactionDetails = append(tx.TransactionDetails, td)

	tx.CompanyId = "system"
	tx.TransactionCode = u.RandomizeString(8, u.ALPHANUMERIC)
	tx.TotalAmount = "0"
	tx.Holder = account.ID
	tx.PartnerId = req.OutletID
	tx.CreatedBy = "system"
	tx.UpdatedBy = "system"

	if _, err = tx.Insert(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	voucher.State = model.VoucherStateUsed
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

type (
	VoucherClaimRequest struct {
	}
)

//PostVoucherClaim :
func PostVoucherClaim(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.Voucher
	var rule model.Rules
	var accountID string
	accountToken := r.FormValue("xx-token")
	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)
	req.State = model.VoucherStateClaim

	token, err := VerifyJWT(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	claims, ok := token.Claims.(*JWTJunoClaims)
	if ok && token.Valid {
		// fmt.Printf("Key:%v", token.Header)
	} else {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
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

	datas := make(map[string]string)
	datas["ACCOUNTID"] = account.ID
	datas["PROGRAMID"] = req.ProgramID

	u.DEBUG("Claim.AccoundID:", accountID)

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

	//Checking Amount
	voucherAmount, err := model.GetVoucherCreatedAmountByProgramID(program.ID)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrFatal)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}
	if int64(voucherAmount+req.VoucherAmount) > program.Stock {
		res.SetError(JSONErrExceedAmount)
		res.JSON(w, res, http.StatusOK)
		return
	}

	//Create Voucher
	var vouchers model.Vouchers
	var vf model.VoucherFormat
	program.VoucherFormat.Unmarshal(&vf)

	for i := 0; i >= req.VoucherAmount; i++ {
		voucher := new(model.Voucher)
		voucher.Code = vf.Properties.Prefix + u.RandomizeString(vf.Properties.Length, u.ALPHANUMERIC) + vf.Properties.Postfix
		voucher.Holder = &account.ID
		voucher.ProgramID = program.ID
		voucher.CreatedBy = "system"
		//*voucher.UpdatedBy = "system"
		voucher.Status = model.StatusCreated
		voucher.State = model.VoucherStateCreated
		voucher.ExpiredAt = program.EndDate
		voucher.HolderDetail.Scan(account)

		vouchers = append(vouchers, *voucher)
	}

	if _, err := vouchers.Insert(); err != nil {
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.JSON(w, res, http.StatusCreated)
}
