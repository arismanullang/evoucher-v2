package controller

import (
	"github.com/gilkor/evoucher/internal/model"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
	"time"
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

			if !strings.Contains(r.URL.Path, "assets") {
				logger.Info("API-Log")
			}
		})
}
