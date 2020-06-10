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

	for idx, trx := range *result {
		partner, _, err := model.GetPartnerByID(qp, trx.PartnerId)
		if err != nil {
			res.SetError(JSONErrFatal.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrFatal.Status)
			return
		}

		selectedPartner := model.Partner{
			ID:          partner.ID,
			Name:        partner.Name,
			Description: partner.Description,
		}

		(*result)[idx].Partner = selectedPartner
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
	transaction, _, err := model.GetTransactionByID(qp, id)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	// partner
	resPartner, _, err := model.GetPartnerByID(qp, transaction.PartnerId)
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	partner := model.Partner{
		ID:          resPartner.ID,
		Name:        resPartner.Name,
		Description: resPartner.Description,
		Tags:        resPartner.Tags,
	}

	transaction.Partner = partner

	listProgramVouchersMap := make(map[string]model.Vouchers)
	listPrograms := model.Programs{}
	listVouchers := model.Vouchers{}
	for _, trxDetail := range transaction.TransactionDetails {
		program, err := model.GetProgramByID(trxDetail.ProgramId, qp)
		if err != nil {
			res.SetError(JSONErrFatal.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrFatal.Status)
			return
		}
		tmpProgram := model.Program{}

		tmpProgram.ID = program.ID
		tmpProgram.Name = program.Name
		tmpProgram.Type = program.Type
		tmpProgram.Value = program.Value
		tmpProgram.MaxValue = program.MaxValue
		tmpProgram.StartDate = program.StartDate
		tmpProgram.EndDate = program.EndDate
		tmpProgram.Description = program.Description
		tmpProgram.ImageURL = program.ImageURL
		tmpProgram.Price = program.Price
		tmpProgram.ProgramChannels = program.ProgramChannels
		tmpProgram.State = program.State
		tmpProgram.Status = program.Status

		listPrograms = append(listPrograms, tmpProgram)

		voucher, err := model.GetVoucherByID(trxDetail.VoucherId, qp)
		if err != nil {
			res.SetError(JSONErrFatal.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrFatal.Status)
			return
		}

		listVouchers = append(listVouchers, *voucher)

		if trxDetail.ProgramId == voucher.ProgramID {
			listProgramVouchersMap[trxDetail.ProgramId] = append(listProgramVouchersMap[voucher.ProgramID], *voucher)
		}

	}
	// program

	for idx, program := range listPrograms {
		listPrograms[idx].Vouchers = listProgramVouchersMap[program.ID]
	}

	transaction.Programs = listPrograms

	res.SetResponse(transaction)
	res.JSON(w, res, http.StatusOK)
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

	var tx model.Transaction

	var f model.VoucherFilter
	f.ID = req.VoucherID

	voucherQP := u.NewQueryParam(r)
	voucherQP.SetFilterModel(f)
	voucherQP.Count = -1

	fmt.Println("voucher filter f value = ", f)

	voucherQP.SetFilterModel(f)

	listVoucherByID, _, err := model.GetVouchers(voucherQP)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	// listProgramID := []string{}

	listProgramVouchersMap := make(map[string][]string)

	// =================================== Loop trough voucherID ===================================
	for _, voucher := range listVoucherByID {
		err := voucher.ValidateVoucher()
		if err != nil {
			res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		// listProgramID = append(listProgramID, voucher.ProgramID)
		listProgramVouchersMap[voucher.ProgramID] = append(listProgramVouchersMap[voucher.ProgramID], voucher.ID)

		var td model.TransactionDetail
		td.ProgramId = voucher.ProgramID
		td.VoucherId = voucher.ID
		td.CreatedBy = account.ID
		td.UpdatedBy = account.ID
		tx.TransactionDetails = append(tx.TransactionDetails, td)
	}

	// =================================== Loop through programID ===================================

	// uniqueProgramID := u.UniqueString(listProgramID)
	var totalDiscount = float64(0)

	listPrograms := model.Programs{}
	for programID, vouchers := range listProgramVouchersMap {
		datas := make(map[string]string)
		datas["ACCOUNTID"] = account.ID
		datas["PROGRAMID"] = programID
		datas["OUTLETID"] = req.OutletID
		datas["TIMEZONE"] = fmt.Sprint(configs["timezone"])

		//Validate Rule Program
		program, err := model.GetProgramByID(programID, qp)
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

		totalDiscount += (program.MaxValue * float64(len(vouchers)))

		tmpProgram := model.Program{}

		tmpProgram.ID = program.ID
		tmpProgram.Name = program.Name
		tmpProgram.Type = program.Type
		tmpProgram.Value = program.Value
		tmpProgram.MaxValue = program.MaxValue
		tmpProgram.StartDate = program.StartDate
		tmpProgram.EndDate = program.EndDate
		tmpProgram.Description = program.Description
		tmpProgram.ImageURL = program.ImageURL
		tmpProgram.Price = program.Price
		tmpProgram.ProgramChannels = program.ProgramChannels
		tmpProgram.State = program.State
		tmpProgram.Status = program.Status

		vouchersByProgram := model.Vouchers{}
		for _, voucher := range listVoucherByID {
			if voucher.ProgramID == programID {
				tempVoucher := model.Voucher{
					ID:        voucher.ID,
					Code:      voucher.Code,
					ExpiredAt: voucher.ExpiredAt,
					ValidAt:   voucher.ValidAt,
					State:     voucher.State,
				}
				vouchersByProgram = append(vouchersByProgram, tempVoucher)
			}
		}
		tmpProgram.Vouchers = vouchersByProgram

		listPrograms = append(listPrograms, tmpProgram)
	}

	partner, _, err := model.GetPartnerByID(qp, req.OutletID)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	//send email confirmation

	selectedPartner := model.Partner{
		ID:          partner.ID,
		Name:        partner.Name,
		Description: partner.Description,
	}

	tx.CompanyId = companyID
	tx.TransactionCode = u.RandomizeString(u.TRANSACTION_CODE_LENGTH, u.NUMERALS)
	tx.TotalAmount = fmt.Sprint(totalDiscount)
	tx.Holder = account.ID
	tx.PartnerId = req.OutletID
	tx.CreatedBy = account.ID
	tx.UpdatedBy = account.ID

	resTrx, err := tx.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	for _, voucher := range listVoucherByID {
		voucher.State = model.VoucherStateUsed
		voucher.UpdatedBy = account.ID
		if err = voucher.Update(); err != nil {
			u.DEBUG(err)
			res.SetErrorWithDetail(JSONErrFatal, err)
			res.JSON(w, res, JSONErrFatal.Status)
			return
		}
	}

	finalResponse := *resTrx
	finalResponse[0].Vouchers = listVoucherByID
	finalResponse[0].Programs = listPrograms
	finalResponse[0].Partner = selectedPartner

	//send email confirmation
	err = finalResponse[0].SendEmailConfirmation()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
	}
	// remove `Vouchers` field on transaction to use combine voucher -> use tx.Programs

	res.SetResponse(finalResponse)
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

//PostPublicVoucherUse : Post public voucher use without token validation
func PostPublicVoucherUse(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req UseTransaction
	var rule model.Rules
	companyID := bone.GetValue(r, "company")
	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)

	//get config TimeZone
	configs, err := model.GetConfigs(companyID, "company")
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.Error.SetMessage("timezone config not found, please add timezone config")
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	//Validate VoucherID
	voucher, err := model.GetVoucherByID(req.VoucherID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	err = voucher.ValidateVoucher()
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	datas := make(map[string]string)
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
	td.CreatedBy = "web"
	td.UpdatedBy = "web"

	tx.TransactionDetails = append(tx.TransactionDetails, td)

	tx.CompanyId = companyID
	tx.TransactionCode = u.RandomizeString(u.TRANSACTION_CODE_LENGTH, u.NUMERALS)
	// Multiply with total voucher used
	tx.TotalAmount = fmt.Sprint(program.MaxValue)
	tx.Holder = *voucher.Holder
	tx.PartnerId = req.OutletID
	tx.CreatedBy = "web"
	tx.UpdatedBy = "web"

	partner, _, err := model.GetPartnerByID(qp, req.OutletID)
	if err != nil {
		res.SetError(JSONErrResourceNotFound)
		res.JSON(w, res, JSONErrResourceNotFound.Status)
		return
	}

	resTrx, err := tx.Insert()
	if err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	voucher.State = model.VoucherStateUsed
	voucher.UpdatedBy = "web"
	if err = voucher.Update(); err != nil {
		u.DEBUG(err)
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	//send email confirmation

	selectedPartner := model.Partner{
		ID:          partner.ID,
		Name:        partner.Name,
		Description: partner.Description,
	}

	tmpProgram := model.Program{}

	tmpProgram.ID = program.ID
	tmpProgram.Name = program.Name
	tmpProgram.Type = program.Type
	tmpProgram.Value = program.Value
	tmpProgram.MaxValue = program.MaxValue
	tmpProgram.StartDate = program.StartDate
	tmpProgram.EndDate = program.EndDate
	tmpProgram.Description = program.Description
	tmpProgram.ImageURL = program.ImageURL
	tmpProgram.Price = program.Price
	tmpProgram.ProgramChannels = program.ProgramChannels
	tmpProgram.State = program.State
	tmpProgram.Status = program.Status

	tmpProgram.Vouchers = model.Vouchers{*voucher}

	finalResponse := *resTrx
	finalResponse[0].Vouchers = model.Vouchers{*voucher}
	finalResponse[0].Programs = model.Programs{tmpProgram}
	finalResponse[0].Partner = selectedPartner

	//send email confirmation
	err = finalResponse[0].SendEmailConfirmation()
	if err != nil {
		res.SetErrorWithDetail(JSONErrFatal, err)
		res.JSON(w, res, JSONErrFatal.Status)
	}

	res.SetResponse(finalResponse[0])
	res.JSON(w, res, http.StatusCreated)
}
