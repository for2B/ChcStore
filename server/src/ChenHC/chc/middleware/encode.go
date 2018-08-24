package middleware

import (
	"net/http"
	"ChenHC/internal/httpapi"
	"ChenHC/chc/constant"
	"ChenHC/utils/feedback"
)

func Encode(next httpapi.APIHandler) httpapi.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		fb := feedback.NewFeedBack(w)
		fb.Code(constant.GLOBAL_SUCCESS)
		response, err := next(w, r)
		if err != nil {
			apiErr, ok := err.(httpapi.Err)
			if !ok {
				fb.Code(constant.GLOBAL_SYS_ERR).Msg("response err invalid").Response()
				return nil, nil
			}
			fb.Code(apiErr.Code)
			fb.Msg(apiErr.Msg)
		}

		fb.Data(response)

		fb.Response()
		return nil, nil
	}
}
