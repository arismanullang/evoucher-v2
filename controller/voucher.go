package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gilkor/evoucher-v2/model"
	u "github.com/gilkor/evoucher-v2/util"
	"github.com/go-zoo/bone"
)

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

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	accountID := claims.AccountID
	req.User = accountID

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
		datas["ACCOUNTID"] = accountID
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

		loc, _ := time.LoadLocation(fmt.Sprint(configs["timezone"]))

		voucherValidAt := time.Now().In(loc)
		voucherExpiredAt := time.Date(voucherValidAt.Year(), voucherValidAt.Month(), voucherValidAt.Day(), 23, 59, 59, 59, loc)

		if ruleUseUsagePeriod, ok := rules.And["rule_use_usage_period"]; ok {

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
		}

		if ruleUseActiveVoucherPeriod, ok := rules.And["rule_use_active_voucher_period"]; ok && !ruleUseActiveVoucherPeriod.IsEmpty() {
			voucherExpiredAt = voucherExpiredAt.AddDate(0, 0, int(ruleUseActiveVoucherPeriod.Eq.(float64)))
		}

		//Create Voucher
		var vouchers model.Vouchers
		var vf model.VoucherFormat
		program.VoucherFormat.Unmarshal(&vf)

		for i := 0; i < reqData.Quantity; i++ {
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
			voucher.CreatedBy = req.User
			voucher.UpdatedBy = req.User
			voucher.Status = model.StatusCreated
			voucher.State = model.VoucherStateCreated
			voucher.ValidAt = &voucherValidAt
			voucher.ExpiredAt = &voucherExpiredAt

			vouchers = append(vouchers, *voucher)
		}

		totalVouchers = append(totalVouchers, vouchers...)
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

//PostVoucherAssignHolder :
func PostVoucherAssignHolder(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()

	var req model.AssignVoucherRequest
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

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	accountID := claims.AccountID
	req.User = accountID

	//get config TimeZone
	configs, err := model.GetConfigs(companyID, "company")
	if err != nil {
		res.SetError(JSONErrBadRequest.SetArgs(err.Error()))
		res.Error.SetMessage("timezone config not found, please add timezone config")
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	// validate each data
	for idx, assignData := range req.Data {

		//Validate Rule Program
		program, err := model.GetProgramByID(assignData.ProgramID, qp)
		if err != nil {
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		datas := make(map[string]interface{})
		datas["ACCOUNTID"] = req.HolderID
		datas["PROGRAMID"] = assignData.ProgramID
		datas["TIMEZONE"] = fmt.Sprint(configs["timezone"])
		datas["QUANTITY"] = len(assignData.VoucherIDs)

		var rules model.RulesExpression
		program.Rule.Unmarshal(&rules)

		result, err := rules.ValidateAssign(datas)
		if err != nil {
			res.SetErrorWithDetail(JSONErrBadRequest, err)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		//check voucher id status / availability
		for _, voucherID := range assignData.VoucherIDs {
			voucherDetail, err := model.GetVoucherByID(voucherID, qp)
			if err != nil {
				res.SetError(JSONErrBadRequest)
				res.JSON(w, res, JSONErrBadRequest.Status)
				return
			}

			// emptyString := ""

			if *voucherDetail.Holder != "" {
				res.SetError(JSONErrBadRequest)
				res.Error.SetMessage("voucher with id " + voucherID + " has been assigned to " + *voucherDetail.Holder)
				res.JSON(w, res, JSONErrBadRequest.Status)
				return
			}

		}

		if !result {
			u.DEBUG(err)
			res.SetError(JSONErrInvalidRule)
			res.JSON(w, err, JSONErrInvalidRule.Status)
			return
		}

		loc, _ := time.LoadLocation(fmt.Sprint(configs["timezone"]))

		voucherValidAt := time.Now().In(loc)
		voucherExpiredAt := time.Date(voucherValidAt.Year(), voucherValidAt.Month(), voucherValidAt.Day(), 23, 59, 59, 59, loc)

		if ruleUseUsagePeriod, ok := rules.And["rule_use_usage_period"]; ok {

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
		}

		if ruleUseActiveVoucherPeriod, ok := rules.And["rule_use_active_voucher_period"]; ok && !ruleUseActiveVoucherPeriod.IsEmpty() {
			voucherExpiredAt = voucherExpiredAt.AddDate(0, 0, int(ruleUseActiveVoucherPeriod.Eq.(float64)))
		}

		assignData.ValidAt = voucherValidAt
		assignData.ExpiredAt = voucherExpiredAt

		// update the request Data
		req.Data[idx] = assignData
	}

	msg, err := req.AssignVoucher()
	if err != nil {
		fmt.Println("err = ", err)
		res.SetError(JSONErrBadRequest)
		res.Error.SetMessage("error assign voucher : " + err.Error())
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	res.SetResponse(msg)
	res.JSON(w, res, http.StatusOK)
}

//GetVoucherByHolder : GET list of program and vouchers by holder
func GetVoucherByHolder(w http.ResponseWriter, r *http.Request) {
	res := u.NewResponse()
	qp := u.NewQueryParam(r)
	qp.Count = -1
	accountToken := r.FormValue("token")

	claims, err := model.VerifyAccountToken(accountToken)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrUnauthorized)
		res.JSON(w, res, JSONErrUnauthorized.Status)
		return
	}

	vouchers, err := model.GetVouchersByHolder(claims.AccountID, qp)
	if err != nil {
		u.DEBUG(err)
		res.SetError(JSONErrBadRequest)
		res.JSON(w, res, JSONErrBadRequest.Status)
		return
	}

	distinctProgram := []string{}
	for _, v := range *vouchers {
		if !u.StringInSlice(v.ProgramID, distinctProgram) {
			distinctProgram = append(distinctProgram, v.ProgramID)
		}
	}

	listPrograms := model.Programs{}
	for _, programID := range distinctProgram {

		program := model.Program{}
		vouchersByProgram := model.Vouchers{}
		partnersByProgram := model.Partners{}

		detailProgram, err := model.GetProgramByID(programID, qp)
		if err != nil {
			u.DEBUG(err)
			res.SetError(JSONErrBadRequest)
			res.JSON(w, res, JSONErrBadRequest.Status)
			return
		}

		program.ID = detailProgram.ID
		program.Name = detailProgram.Name
		program.Type = detailProgram.Type
		program.Value = detailProgram.Value
		program.MaxValue = detailProgram.MaxValue
		program.StartDate = detailProgram.StartDate
		program.EndDate = detailProgram.EndDate
		program.Description = detailProgram.Description
		program.ImageURL = detailProgram.ImageURL
		program.Price = detailProgram.Price
		program.ProgramChannels = detailProgram.ProgramChannels
		program.State = detailProgram.State
		program.Status = detailProgram.Status

		for _, voucher := range *vouchers {
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
		program.Vouchers = vouchersByProgram

		for _, outlet := range detailProgram.Partners {
			tempOutlet := model.Partner{
				ID:          outlet.ID,
				Name:        outlet.Name,
				Description: outlet.Description,
				Status:      outlet.Status,
			}
			partnersByProgram = append(partnersByProgram, tempOutlet)
		}

		program.Partners = partnersByProgram

		listPrograms = append(listPrograms, program)
	}

	res.SetResponse(listPrograms)
	res.JSON(w, res, http.StatusOK)
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
