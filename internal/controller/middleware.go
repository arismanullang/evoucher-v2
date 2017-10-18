package controller

import (
	"github.com/gilkor/evoucher/internal/model"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
)

func LoggerMiddleware() negroni.Handler {
	return negroni.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
			logger := model.NewLog()
			logger.SetStart(time.Now()).
				SetMethod(r.Method).
				SetHost(r.Host).
				SetPath(r.URL.Path)

			n(w, r)

			res := w.(negroni.ResponseWriter)
			logger.SetEnd(time.Now()).
				SetDelta(time.Since(logger.StartTime).Seconds()).
				SetStatus(res.Status())

			var reqData string

			if r.Method == "GET" {
				reqData = r.URL.RawQuery
			}else {
				 body, _ := ioutil.ReadAll(r.Body)
				 reqData = string(body)
			}


			if !strings.Contains(r.URL.Path, "assets") {
				logger.Info("request: ",reqData )
			}
		})
}
