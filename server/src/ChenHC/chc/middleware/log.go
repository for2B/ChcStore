package middleware

import (
	"net/http"
	"ChenHC/internal/httpapi"
	"ChenHC/internal/httpparse"
	"time"

	"go.uber.org/zap"
)

func Log(logger *zap.Logger) httpapi.MiddlewareFunc {
	return func(next httpapi.APIHandler) httpapi.APIHandler {
		return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
			start := time.Now()
			response, err := next(w, r)
			errMsg := "success"
			if err != nil {
				errMsg = err.Error()
			}
			elapsed := time.Since(start)
			body, _ := httpparse.GetBody(r)
			form, _ := httpparse.GetForm(r)
			multipartFormFile, _ := httpparse.GetMultipartFormFile(r)
			multipartFormValue, _ := httpparse.GetMultipartFormValue(r)

			logger.Info("api log",
				zap.String("url", r.URL.RequestURI()),
				zap.String("method", r.Method),
				zap.String("remoteAddr", r.RemoteAddr),
				zap.String("errMsg", errMsg),
				zap.String("elapsed", elapsed.String()),
				zap.String("body", body),
				zap.Reflect("form", form),
				zap.Reflect("multipartFormFile", multipartFormFile),
				zap.Reflect("multipartFormValue", multipartFormValue),
			)

			return response, err
		}
	}
}
