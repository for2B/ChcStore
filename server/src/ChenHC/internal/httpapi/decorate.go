package httpapi

import (
	"net/http"
)
 //中间件相关定义和处理
type MiddlewareFunc func(APIHandler) APIHandler  //中间件函数定义

type APIHandler func(w http.ResponseWriter, r *http.Request) (response interface{}, err error)  //Handler函数

func Decorate(h APIHandler, ds ...MiddlewareFunc) http.HandlerFunc {
	decorated := h
	for i := len(ds) - 1; i >= 0; i-- {
		decorated = ds[i](decorated)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		decorated(w, r) //返回最顶层执行的函数体
	}
}
//获取MiddlewareFunc数组ds 返回将ds遍历执行一遍的函数MiddlewareFunc
func CombinationMiddleware(ds ...MiddlewareFunc) MiddlewareFunc {
	return func(w APIHandler) APIHandler {
		for i := len(ds) - 1; i >= 0; i-- {
			w = ds[i](w) //遍历构造一层接一层的中间件
		}
		return w
	}
}
