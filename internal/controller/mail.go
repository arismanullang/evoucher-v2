package controller

import (
	"github.com/gilkor/evoucher/internal/model"
	"github.com/ruizu/render"
	"net/http"
)

type (
	EmailRequest struct {
		ProgramId string
	}
)

func SendForgotPasswordMail(w http.ResponseWriter, r *http.Request) {
	var username = r.FormValue("username")
	logger := model.NewLog()
	logger.SetService("API").
		SetMethod(r.Method).
		SetTag("Forgot Password : " + username)

	if err := model.SendMailForgotPassword(model.Domain, model.ApiKey, model.PublicApiKey, username); err != nil {
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
		SendSedayuOneEmail(w, r, logger)
	} else {
		res := NewResponse(nil)
		res.AddError(its(http.StatusNotFound), model.ErrCodeResourceNotFound, model.ErrResourceNotFound.Error(), logger.TraceID)
		render.JSON(w, res, http.StatusOK)
	}
}

func SendSedayuOneEmail(w http.ResponseWriter, r *http.Request, logger *model.LogField) {
	listEmail := []string{"richard@gilkor.com", "andrie@gilkor.com", "oscar@gilkor.com", "novariandy@gilkor.com "}
	listParam := []model.SedayuOneEmail{}

	listParam = append(listParam, model.SedayuOneEmail{Name: "Richard", VoucherUrl: "voucher.elys.id"})
	listParam = append(listParam, model.SedayuOneEmail{Name: "Andrie", VoucherUrl: "voucher.elys.id"})
	listParam = append(listParam, model.SedayuOneEmail{Name: "Oscar", VoucherUrl: "voucher.elys.id"})
	listParam = append(listParam, model.SedayuOneEmail{Name: "Nova", VoucherUrl: "voucher.elys.id"})

	if err := model.SendMailSedayuOne(model.Domain, model.ApiKey, model.PublicApiKey, "Sedayu One Voucher Test", listEmail, listParam); err != nil {
		res := NewResponse(nil)
		status := http.StatusInternalServerError
		errTitle := model.ErrCodeInternalError
		if err == model.ErrResourceNotFound {
			status = http.StatusNotFound
			errTitle = model.ErrCodeResourceNotFound
		}

		res.AddError(its(status), errTitle, err.Error(), logger.TraceID)
		logger.SetStatus(status).Info("param :", listParam, "response :", err.Error())
		render.JSON(w, res, status)
		return
	}

	render.JSON(w, http.StatusOK)
}
