package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
	ResponseReport struct {
		Chart interface{} `json:"chart"`
		Data  interface{} `json:"data"`
	}
	ReportLine struct {
		Label string      `json:"label"`
		Color string      `json:"color"`
		Data  [][2]string `json:"data"`
	}
	ReportFlotBar struct {
		Label string      `json:"label"`
		Bars  FlotBar     `json:"bars"`
		Data  [][2]string `json:"data"`
	}
	FlotBar struct {
		Order     int    `json:"order"`
		FillColor string `json:"fillColor"`
	}
)

func MakeReport(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	result, err := model.MakeReport(id)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportLine{}
	m := [][2]string{}
	label := result[0].Creator
	color := [5]string{"#00BCD4", "#CDDC39", "#FF5722", "#42f44b", "#ff0000"}
	indexColor := 0
	for _, v := range result {
		if v.Creator != label {
			fmt.Println("go " + v.Creator)
			temp := ReportLine{
				Label: label,
				Color: color[indexColor%5],
				Data:  m,
			}

			resultVal = append(resultVal, temp)

			m = [][2]string{}
			label = v.Creator
			indexColor++
		}

		mm := [2]string{}
		mm[0] = v.Name
		mm[1] = v.Total
		m = append(m, mm)
	}

	temp := ReportLine{
		Label: label,
		Color: color[indexColor%5],
		Data:  m,
	}

	resultVal = append(resultVal, temp)
	response := ResponseReport{
		Chart: resultVal,
		Data:  result,
	}
	res := NewResponse(response)
	render.JSON(w, res, http.StatusOK)
}

func MakeReportProgram(w http.ResponseWriter, r *http.Request) {
	result, err := model.MakeReportProgram()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportFlotBar{}
	m := [][2]string{}
	label := result[0].Username
	color := [5]string{"#CEE7FB", "#F4D989", "#FF5722", "#42f44b", "#ff0000"}
	indexColor := 0
	for _, v := range result {
		if v.Username != label {
			fmt.Println("go " + v.Username)
			bars := FlotBar{
				Order:     indexColor,
				FillColor: color[indexColor%5],
			}
			temp := ReportFlotBar{
				Label: label,
				Bars:  bars,
				Data:  m,
			}

			resultVal = append(resultVal, temp)

			m = [][2]string{}
			label = v.Username
			indexColor++
		}

		mm := [2]string{}
		mm[0] = v.Month
		mm[1] = v.Total
		m = append(m, mm)
	}

	bars := FlotBar{
		Order:     indexColor,
		FillColor: color[indexColor%5],
	}
	temp := ReportFlotBar{
		Label: label,
		Bars:  bars,
		Data:  m,
	}

	resultVal = append(resultVal, temp)
	response := ResponseReport{
		Chart: resultVal,
		Data:  result,
	}
	res := NewResponse(response)
	render.JSON(w, res, http.StatusOK)
}

func MakeCompleteReportVoucherByUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	result, err := model.MakeCompleteReportVoucherByUser(id)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportFlotBar{}
	quota := make(map[string]int)

	m := [][2]string{}
	label := result[0].State
	color := [5]string{"#00BCD4", "#CDDC39", "#FF5722", "#42f44b", "#ff0000"}
	indexColor := 0
	for _, v := range result {
		if v.State != label {
			fmt.Println("go " + v.Creator)

			bars := FlotBar{
				Order:     indexColor,
				FillColor: color[indexColor%5],
			}
			temp := ReportFlotBar{
				Label: label,
				Bars:  bars,
				Data:  m,
			}

			resultVal = append(resultVal, temp)

			m = [][2]string{}
			label = v.State
			indexColor++
		}

		mm := [2]string{}
		mm[0] = v.Name
		mm[1] = v.Total
		m = append(m, mm)

		quota[v.Name] = v.Quota
	}

	bars := FlotBar{
		Order:     indexColor,
		FillColor: color[indexColor%5],
	}
	temp := ReportFlotBar{
		Label: label,
		Bars:  bars,
		Data:  m,
	}
	resultVal = append(resultVal, temp)

	m = [][2]string{}
	for k, v := range quota {
		fmt.Println(k)
		fmt.Println(v)

		mm := [2]string{}
		mm[0] = k
		mm[1] = strconv.Itoa(v)
		m = append(m, mm)
	}
	indexColor++
	bars = FlotBar{
		Order:     indexColor,
		FillColor: color[indexColor%5],
	}
	temp = ReportFlotBar{
		Label: "remaining",
		Bars:  bars,
		Data:  m,
	}
	resultVal = append(resultVal, temp)

	response := ResponseReport{
		Chart: resultVal,
		Data:  result,
	}
	res := NewResponse(response)
	render.JSON(w, res, http.StatusOK)
}

func MakeReportVoucherByUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	result, err := model.MakeReportVoucherByUser(id)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportLine{}
	color := [5]string{"#CEE7FB", "#F4D989", "#FF5722", "#42f44b", "#ff0000"}
	indexColor := 0

	m := [][2]string{}
	label := "voucher"
	for _, v := range result {
		mm := [2]string{}
		mm[0] = v.Name
		mm[1] = v.Total
		m = append(m, mm)
	}
	temp := ReportLine{
		Label: label,
		Color: color[indexColor%5],
		Data:  m,
	}
	resultVal = append(resultVal, temp)
	indexColor++

	m = [][2]string{}
	label = "quota"
	for _, v := range result {
		mm := [2]string{}
		mm[0] = v.Name
		mm[1] = v.Quota
		m = append(m, mm)
	}
	temp = ReportLine{
		Label: "remaining",
		Color: color[indexColor%5],
		Data:  m,
	}
	resultVal = append(resultVal, temp)

	response := ResponseReport{
		Chart: resultVal,
		Data:  result,
	}
	res := NewResponse(response)
	render.JSON(w, res, http.StatusOK)
}

func MakeReportLine(w http.ResponseWriter, r *http.Request) {
	result, err := model.MakeReportProgram()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportLine{}
	m := [][2]string{}
	label := result[0].Creator
	color := [5]string{"#00BCD4", "#CDDC39", "#FF5722", "#42f44b", "#ff0000"}
	indexColor := 0
	for _, v := range result {
		if v.Creator != label {
			fmt.Println("go " + v.Creator)
			temp := ReportLine{
				Label: label,
				Color: color[indexColor%5],
				Data:  m,
			}

			resultVal = append(resultVal, temp)

			m = [][2]string{}
			label = v.Creator
			indexColor++
		}

		mm := [2]string{}
		mm[0] = v.Month
		mm[1] = v.Total
		m = append(m, mm)
	}

	temp := ReportLine{
		Label: label,
		Color: color[indexColor%5],
		Data:  m,
	}

	resultVal = append(resultVal, temp)
	res := NewResponse(resultVal)
	render.JSON(w, res, http.StatusOK)
}
