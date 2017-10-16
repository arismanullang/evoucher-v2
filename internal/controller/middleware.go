package controller

import (
	"time"
	"net/http"
	"github.com/urfave/negroni"
	"github.com/gilkor/evoucher/internal/model"
)

func LoggerMiddleware() negroni.Handler{
	return negroni.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request , n http.HandlerFunc) {
			logger  := model.NewLog()
			logger.SetStart(time.Now()).
				SetMethod(r.Method).
				SetHost(r.Host).
				SetPath(r.URL.Path)

			n(w,r)

			res := w.(negroni.ResponseWriter)
			logger.SetEnd(time.Now()).
				SetDelta(time.Since(logger.StartTime).Seconds()).
				SetStatus(res.Status())

			logger.Info("API-Log" )
		})
}
