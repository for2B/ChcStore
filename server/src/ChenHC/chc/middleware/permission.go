package middleware

import (
	"net/http"
	"ChenHC/internal/httpapi"
	"ChenHC/chc/constant"
)

func Permission(LevelNeeded string) httpapi.MiddlewareFunc {
	return func(next httpapi.APIHandler) httpapi.APIHandler {
		return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {

			// return next(w, r)
			session, ok := r.Context().Value("session").(*constant.Session)
			if !ok {
				return nil, httpapi.NewErr(constant.GLOBAL_SYS_ERR, "Permission->session assert failed", nil)
			}

			userLevel := session.Level

			//judge the level
			switch userLevel {
			case constant.USER_SUPER_ADMIN:
				return next(w, r)
			case constant.USER_ADMIN:
				if LevelNeeded != constant.USER_SUPER_ADMIN {
					return next(w, r)
				}
			case constant.USER_NORMAL:
				if LevelNeeded == constant.USER_NORMAL {
					return next(w, r)
				}
			}
			return nil, httpapi.Err{constant.GLOBAL_NO_AUTH, "No permission!"}

		}
	}
}
