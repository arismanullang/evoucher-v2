package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
)

type (
	EmailRequest struct {
		ProgramId string
	}
	CreateImageCampaignRequest struct {
		ProgramID    string `json:"program_id"`
		ImageHeader  string `json:"image_header"`
		ImageVoucher string `json:"image_voucher"`
		ImageFooter  string `json:"image_footer"`
	}
	CreateImageCampaignRequestV2 struct {
		ProgramID    string `json:"program_id"`
		Template     string `json:"template"`
		EmailSubject string `json:"email_subject"`
		EmailSender  string `json:"email_sender"`
		EmailContent string `json:"email_content"`
		ImageHeader  string `json:"image_header"`
		ImageVoucher string `json:"image_voucher"`
		ImageFooter  string `json:"image_footer"`
	}
	SendEmailCampaignRequest struct {
		ProgramID  string   `json:"program_id"`
		EmailUsers []string `json:"email_user"`
	}
)

func CreateEmailCampaign(w http.ResponseWriter, r *http.Request) {
	apiName := "update_campaign"
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse("")
	status := http.StatusCreated

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	var rd CreateImageCampaignRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusBadRequest).Panic("param :", rd, "response :", err.Error())
	}

	campaign := model.ProgramCampaign{
		ProgramID:    rd.ProgramID,
		ImageHeader:  rd.ImageHeader,
		ImageVoucher: rd.ImageVoucher,
		ImageFooter:  rd.ImageFooter,
	}

	id, err := model.InsertCampaign(campaign, a.User.ID)
	res = NewResponse(id)
	logger.SetStatus(status).Info("param :", rd, "response :", id)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func CreateEmailCampaignV2(w http.ResponseWriter, r *http.Request) {
	apiName := "update_campaign"
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse("")
	status := http.StatusCreated

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	var rd CreateImageCampaignRequestV2
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusBadRequest).Panic("param :", rd, "response :", err.Error())
	}

	campaign := model.ProgramCampaignV2{
		ProgramID:    rd.ProgramID,
		EmailSubject: rd.EmailSubject,
		EmailSender:  rd.EmailSender,
		EmailContent: rd.EmailContent,
		Template:     "email_campaign_demo",
		ImageHeader:  rd.ImageHeader,
		ImageVoucher: rd.ImageVoucher,
		ImageFooter:  rd.ImageFooter,
	}

	id, err := model.InsertCampaignV2(campaign, a.User.ID)
	res = NewResponse(id)
	logger.SetStatus(status).Info("param :", rd, "response :", id)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd, "response :", res.Errors)
	}

	render.JSON(w, res, status)
}

func GetEmailCampaignByProgram(w http.ResponseWriter, r *http.Request) {
	apiName := "get_campaign"
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	res := NewResponse("")
	status := http.StatusOK

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	id := r.FormValue("id")

	campaign, err := model.GetCampaignV2(id)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", id, "response :", res.Errors)
	}

	res = NewResponse(campaign)
	logger.SetStatus(status).Info("param :", id, "response :", campaign)

	render.JSON(w, res, status)
}

func SendEmailCampaign(w http.ResponseWriter, r *http.Request) {
	apiName := "send_campaign"
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag(apiName)

	var rd SendEmailCampaignRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rd); err != nil {
		logger.SetStatus(http.StatusBadRequest).Panic("param :", rd, "response :", err.Error())
	}

	res := NewResponse("")
	status := http.StatusOK

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		res = a.res
		status = http.StatusUnauthorized
		render.JSON(w, res, status)
		return
	}

	campaign, err := model.GetCampaignV2(rd.ProgramID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", rd.ProgramID, "response :", res.Errors)
	}

	user, err := model.GetEmailUserByIDs(rd.EmailUsers)
	if err != nil && err != model.ErrResourceNotFound {
		log.Panic(err)
	}

	program, err := model.FindProgramDetailsById(rd.ProgramID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd.ProgramID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	gvr := GenerateVoucherRequest{
		AccountID:   a.User.Account.Id,
		ProgramID:   rd.ProgramID,
		Quantity:    1,
		CreatedBy:   a.User.ID,
		ReferenceNo: rd.ProgramID,
	}

	target := []model.TargetEmail{}
	for _, v := range user {
		gvr.Holder.Key = v.Email
		gvr.Holder.Email = v.Email
		gvr.Holder.Description = v.Name

		vo, err := gvr.generateVoucher(&program)
		if err != nil {
			fmt.Println(err)
			model.RollbackVoucher(vo[0].ID, a.User.Account.Id)

			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", rd.ProgramID, "response :", res.Errors.ToString())
			render.JSON(w, res, status)
			return
		}

		tempTarget := model.TargetEmail{
			HolderEmail: v.Email,
			HolderName:  v.Name,
			VoucherUrl:  generateLink(vo[0].ID),
		}
		target = append(target, tempTarget)
	}

	if err := model.SendVoucherMailV2(model.Domain, model.ApiKey, model.PublicApiKey, campaign, target); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", campaign, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	insertCampaign := model.InsertCampaignUserRequest{
		CampaignID: campaign.ID,
		EmailUsers: rd.EmailUsers,
		CreatedBy:  a.User.ID,
	}
	if err := model.InsertCampaignUser(insertCampaign); err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", rd.ProgramID, "response :", res.Errors.ToString())
		render.JSON(w, res, status)
		return
	}

	res = NewResponse("")
	render.JSON(w, res, status)
}

func SendForgotPasswordMail(w http.ResponseWriter, r *http.Request) {
	var username = r.FormValue("username")
	var accountId = r.FormValue("accountId")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Forgot Password : " + username)

	if err := model.SendMailForgotPassword(model.Domain, model.ApiKey, model.PublicApiKey, username, accountId); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", username, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}
	render.JSON(w, http.StatusOK)
}

func SendCustomMailRoute(w http.ResponseWriter, r *http.Request) {
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("send_mail")

	route := r.FormValue("route")
	if route == "sone" {
		fmt.Println("sone")
		//SendSedayuOneEmailTest(w, r)
	} else {
		res := NewResponse(nil)
		res.AddError(its(http.StatusNotFound), model.ErrCodeResourceNotFound, model.ErrResourceNotFound.Error(), logger.TraceID)
		render.JSON(w, res, http.StatusOK)
	}
}

func SendSedayuOneEmail(w http.ResponseWriter, r *http.Request) {
	apiName := "voucher_generate-bulk"
	var gvd GenerateVoucherRequest
	var status int
	res := NewResponse(nil)
	vrID := r.FormValue("program")

	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("send_mail")

	a := AuthTokenWithLogger(w, r, logger)
	if !a.Valid {
		render.JSON(w, a.res, http.StatusUnauthorized)
		return
	}

	if CheckAPIRole(a, apiName) {
		logger.SetStatus(status).Info("param :", a.User.ID, "response :", "Invalid Role")

		status = http.StatusUnauthorized
		res.AddError(its(status), model.ErrCodeInvalidRole, model.ErrInvalidRole.Error(), logger.TraceID)
		render.JSON(w, res, status)
		return
	}

	if getCountVoucher(vrID) > 0 {
		status = http.StatusBadRequest
		res.AddError(its(status), model.ErrCodeInvalidProgram, model.ErrMessageProgramHasBeenUsed, logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	program, err := model.FindProgramDetailsById(vrID)
	if err == model.ErrResourceNotFound {
		status = http.StatusNotFound
		res.AddError(its(status), model.ErrCodeResourceNotFound, model.ErrMessageResourceNotFound, logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	} else if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	var listBroadcast []model.BroadcastUser
	listBroadcast, err = model.FindBroadcastUser(map[string]string{"program_id": vrID})
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}

	gvd.AccountID = a.User.Account.Id
	gvd.ProgramID = vrID
	gvd.Quantity = 1
	gvd.CreatedBy = a.User.ID

	totalVoucher := []model.Voucher{}
	tempListVoucher := []model.Voucher{}

	for _, v := range listBroadcast {
		gvd.ReferenceNo = its(v.ID)
		gvd.Holder.Key = v.Target
		gvd.Holder.Email = v.Target
		gvd.Holder.Description = v.Description

		tempListVoucher, err = gvd.generateVoucher(&program)
		if err != nil {
			fmt.Println(err)
			rollback(vrID)

			status = http.StatusInternalServerError
			res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
			logger.SetStatus(status).Log("param :", vrID, "response :", err.Error())
			render.JSON(w, res, status)
			return
		}

		for _, vv := range tempListVoucher {
			totalVoucher = append(totalVoucher, vv)
		}
	}

	campaign, err := model.GetCampaign(vrID)
	if err != nil {
		status = http.StatusInternalServerError
		res.AddError(its(status), model.ErrCodeInternalError, model.ErrMessageInternalError+"("+err.Error()+")", logger.TraceID)
		logger.SetStatus(status).Log("param :", vrID, "response :", res.Errors)
		render.JSON(w, res, status)
		return
	}
	campaign.AccountID = a.User.Account.Id
	campaign.CreatedBy = a.User.ID
	listEmail := []model.TargetEmail{}

	for _, v := range totalVoucher {
		listEmail = append(listEmail, model.TargetEmail{HolderName: v.HolderDescription.String, VoucherUrl: generateLink(v.ID), HolderEmail: v.Holder.String})
	}

	if err := model.SendMailSedayuOne(model.Domain, model.ApiKey, model.PublicApiKey, "Sedayu One Voucher", listEmail, campaign); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", listEmail, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	status = http.StatusCreated
	res = NewResponse("success")
	logger.SetStatus(status).Log("param :", vrID, "response :", res.Data)
	render.JSON(w, res, status)
	return
}
