package middleware

import (
	"context"
	"net/http"
	"ChenHC/internal/httpapi"
	"ChenHC/chc/constant"
	"ChenHC/chc/session"
)

func Session(sf *session.SessionManager) httpapi.MiddlewareFunc {
	return func(next httpapi.APIHandler) httpapi.APIHandler {
		return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {

			// return next(w, r)
			sess := &constant.Session{}
			//mock a session in local enviroment
			if sf.Build_type == constant.BUILD_TYPE_DEPLOY {
				var err error
				sess, err = sf.GetSession(w, r)
				if err != nil {
					return nil, err
				}

			} else {

				sess = &constant.Session{
					Userid:           "1606100101",
					Name:             "MockMan",
					Level:            "SUPER",
					Identity:         "中共党员",
					Department_array: []string{"%计算机%", "%生命%"},
				}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "session", sess)
			return next(w, r.WithContext(ctx))
		}
	}
}
