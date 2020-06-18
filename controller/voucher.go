package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx/types"
)

//GetVoucherByID : get voucher by id
func GetVoucherByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	id := bone.GetValue(r, "id")

	vouchers, err := model.GetVoucherByID(id, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	td, err := model.GetTransactionDetailByVoucherID(qp, id)
	if err == nil {
		fmt.Println("vouchers = ", td)
		transactionDetail := *td
		vouchers.TransactionDetail = &transactionDetail[0]
	}

	res.SetResponse(vouchers)
	res.JSON(w, res, http.StatusOK)
}

//PostVoucherInjectByHolder : Inject voucher by holder
func PostVoucherInjectByHolder(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.InjectVoucherByHolderRequest
	decoder := json.NewDecoder(r.Body)
	qp := u.NewQueryParam(r)
	err := decoder.Decode(&req)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	companyID := bone.GetValue(r, "company")
	accountToken := r.FormValue("token")

	auth, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	accountID := auth.AccountID
	req.UpdatedBy = accountID

	//get config TimeZone
	configs, err := model.GetConfigs(companyID, "company")
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.Error.SetMessage("timezone config not found, please add timezone config")
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	totalVouchers := model.Vouchers{}

	for _, reqData := range req.Data {
		//Validate Rule Program
		program, err := model.GetProgramByID(reqData.ProgramID, qp)
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		var rules model.RulesExpression
		program.Rule.Unmarshal(&rules)

		datas := make(map[string]interface{})
		datas["ACCOUNTID"] = req.HolderID
		datas["PROGRAMID"] = reqData.ProgramID
		datas["QUANTITY"] = reqData.Quantity
		datas["TIMEZONE"] = fmt.Sprint(configs["timezone"])

		result, err := rules.ValidateClaim(datas)
		if err != nil {
			u.DEBUG(err)
			res.SetErrorWithDetail(JSONErrBadRequest, err)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		currentClaimedVoucher, err := model.GetVoucherCreatedAmountByProgramID(reqData.ProgramID)
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}
		if int64(currentClaimedVoucher+reqData.Quantity) > program.Stock {
			u.DEBUG(err)
			res.SetError(JSONErrExceedAmount)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		if !result {
			u.DEBUG(err)
			res.SetError(JSONErrInvalidRule)
			res.JSON(w, err, JSONErrInvalidRule.Status)
			return
		}

		gvr := GenerateVoucherRequest{
			ProgramID:    program.ID,
			Quantity:     reqData.Quantity,
			ReferenceNo:  req.ReferenceNo,
			HolderID:     req.HolderID,
			HolderDetail: req.HolderDetail,
			UpdatedBy:    req.UpdatedBy,
		}

		vouchers, err := gvr.GenerateVoucher(fmt.Sprint(configs["timezone"]), *program)
		if err != nil {
			res.SetError(JSONErrFatal.SetArgs(err.Error()))
			res.JSON(w, res, JSONErrFatal.Status)
		}

		totalVouchers = append(totalVouchers, *vouchers...)
	}

	response, err := totalVouchers.Insert()
	if err != nil {
		fmt.Println(err)
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)

	res.JSON(w, res, http.StatusCreated)
}

//GetPublicVoucherByID : GET list of program and vouchers by voucher_id
func GetPublicVoucherByID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)
	qp.Count = -1
	encodedVoucherID := r.FormValue("x")
	encodedCompanyID := r.FormValue("y")

	voucherID := u.StrDecode(encodedVoucherID)
	companyID := u.StrDecode(encodedCompanyID)

	qp.SetCompanyID(companyID)

	voucher, err := model.GetVoucherByID(voucherID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	vouchers := model.Vouchers{}

	r.Form.Set("fields", model.MProgramFields)
	qp2 := u.NewQueryParam(r)
	qp2.SetCompanyID(companyID)
	program, err := model.GetProgramByID(voucher.ProgramID, qp2)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	r.Form.Set("fields", model.MOutletFields)
	qp3 := u.NewQueryParam(r)
	qp3.SetCompanyID(companyID)
	outlets, _, err := model.GetOutletByProgramID(voucher.ProgramID, qp3)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	program.Outlets = *outlets
	vouchers = append(vouchers, *voucher)
	program.Vouchers = vouchers

	res.SetResponse(program)
	res.JSON(w, res, http.StatusOK)
}

//GetVoucherByToken : GET list of program and vouchers by holder juno token
func GetVoucherByToken(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	// r.Form.Set("fields", model.MVoucherFields)
	qp := u.NewQueryParam(r)
	qp.Count = -1
	qp.Sort = "expired_at-"
	token := r.FormValue("token")

	qp.SetCompanyID(bone.GetValue(r, "company"))
	var f model.VoucherFilter
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(&f, r.Form); err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	qp.SetFilterModel(f)

	accData, err := model.GetSessionDataJWT(token)
	if err != nil {
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	// vouchers, _, err := model.GetVouchers(qp)
	vouchers, err := model.GetVouchersByHolder(accData.AccountID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(vouchers)
	res.JSON(w, res, http.StatusOK)
}

type (
	//GenerateVoucherRequest : generate voucher based on spesific program and holder
	GenerateVoucherRequest struct {
		ProgramID    string         `json:"program_id"`
		Quantity     int            `json:"quantity"`
		ReferenceNo  string         `json:"reference_no"`
		UpdatedBy    string         `json:"updated_by"`
		HolderID     string         `json:"holder_id"`
		HolderDetail types.JSONText `json:"holder_detail"`
	}
)

func (req *GenerateVoucherRequest) GenerateVoucher(timezone string, program model.Program) (*model.Vouchers, error) {

	var rules model.RulesExpression
	program.Rule.Unmarshal(&rules)

	loc, _ := time.LoadLocation(timezone)

	voucherValidAt := time.Now().In(loc)
	voucherExpiredAt := time.Date(voucherValidAt.Year(), voucherValidAt.Month(), voucherValidAt.Day(), 23, 59, 59, 59, loc)

	if ruleUseUsagePeriod, ok := rules.And["rule_use_usage_period"]; ok {

		validTime, err := model.StringToTime(fmt.Sprint(ruleUseUsagePeriod.Gte))
		if err != nil {
			return &model.Vouchers{}, err
		}

		expiredTime, err := model.StringToTime(fmt.Sprint(ruleUseUsagePeriod.Lte))
		if err != nil {
			return &model.Vouchers{}, err
		}

		voucherValidAt = validTime
		voucherExpiredAt = expiredTime
	}

	if ruleUseActiveVoucherPeriod, ok := rules.And["rule_use_active_voucher_period"]; ok && !ruleUseActiveVoucherPeriod.IsEmpty() {
		voucherExpiredAt = voucherExpiredAt.AddDate(0, 0, int(ruleUseActiveVoucherPeriod.Eq.(float64)))
	}

	var vouchers model.Vouchers
	var vf model.VoucherFormat
	program.VoucherFormat.Unmarshal(&vf)

	for i := 0; i < req.Quantity; i++ {
		voucher := new(model.Voucher)
		if vf.Type == "fix" {
			voucher.Code = vf.Code
		} else if vf.Type == "random" {
			voucher.Code = vf.Prefix + u.RandomizeString(u.DEFAULT_LENGTH, vf.Random) + vf.Postfix
		}

		voucher.ReferenceNo = req.ReferenceNo
		voucher.Holder = &req.HolderID
		voucher.HolderDetail = req.HolderDetail
		voucher.ProgramID = program.ID
		voucher.CreatedBy = req.UpdatedBy
		voucher.UpdatedBy = req.UpdatedBy
		voucher.Status = model.StatusCreated
		voucher.State = model.VoucherStateCreated
		voucher.ValidAt = &voucherValidAt
		voucher.ExpiredAt = &voucherExpiredAt

		vouchers = append(vouchers, *voucher)
	}

	return &vouchers, nil
}

// DeleteVoucher : Delete voucher by id
func DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	qp := u.NewQueryParam(r)
	id := bone.GetValue(r, "id")

	accountToken := r.FormValue("token")

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	voucher, err := model.GetVoucherByID(id, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	voucher.UpdatedBy = claims.AccountID
	// Delete voucher
	fmt.Println("delete voucher ", voucher)
	response, err := voucher.Delete()
	if err != nil {
		res.SetError(JSONErrFatal.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrFatal.Status)
		return
	}

	res.SetResponse(response)
	res.JSON(w, res, http.StatusOK)
}

// GetVoucherByProgramID : get list of vouchers by programID
func GetVoucherByProgramID(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)

	// need to check program company_id
	// qp.SetCompanyID(bone.GetValue(r, "company"))
	programID := bone.GetValue(r, "id")

	vouchers, next, err := model.GetVouchersByProgramID(programID, qp)
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	if len(vouchers) > 0 {
		res.SetNewPagination(r, qp.Page, next, (vouchers)[0].Count)
	}

	res.SetResponse(vouchers)
	res.JSON(w, res, http.StatusOK)
}
