package controller

import (
	"fmt"
	"net/http"

	"github.com/ruizu/render"

	"github.com/gilkor/evoucher/internal/model"
)

type (
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
	color := [3]string{"#00BCD4", "#CDDC39", "#FF5722"}
	indexColor := 0
	for _, v := range result {
		if v.Creator != label {
			fmt.Println("go " + v.Creator)
			temp := ReportLine{
				Label: label,
				Color: color[indexColor],
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
		Color: color[indexColor],
		Data:  m,
	}

	resultVal = append(resultVal, temp)
	res := NewResponse(resultVal)
	render.JSON(w, res, http.StatusOK)
}

func MakeReportVariantFlotBar(w http.ResponseWriter, r *http.Request) {
	result, err := model.MakeReportVariant()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportFlotBar{}
	m := [][2]string{}
	label := result[0].Creator
	color := [3]string{"#00BCD4", "#CDDC39", "#FF5722"}
	indexColor := 0
	for _, v := range result {
		if v.Creator != label {
			fmt.Println("go " + v.Creator)
			bars := FlotBar{
				Order:     indexColor,
				FillColor: color[indexColor],
			}
			temp := ReportFlotBar{
				Label: label,
				Bars:  bars,
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

	bars := FlotBar{
		Order:     indexColor,
		FillColor: color[indexColor],
	}
	temp := ReportFlotBar{
		Label: label,
		Bars:  bars,
		Data:  m,
	}

	resultVal = append(resultVal, temp)
	res := NewResponse(resultVal)
	render.JSON(w, res, http.StatusOK)
}

func MakeReportLine(w http.ResponseWriter, r *http.Request) {
	result, err := model.MakeReportVariant()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	resultVal := []ReportLine{}
	m := [][2]string{}
	label := result[0].Creator
	color := [3]string{"#00BCD4", "#CDDC39", "#FF5722"}
	indexColor := 0
	for _, v := range result {
		if v.Creator != label {
			fmt.Println("go " + v.Creator)
			temp := ReportLine{
				Label: label,
				Color: color[indexColor],
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
		Color: color[indexColor],
		Data:  m,
	}

	resultVal = append(resultVal, temp)
	res := NewResponse(resultVal)
	render.JSON(w, res, http.StatusOK)
}
