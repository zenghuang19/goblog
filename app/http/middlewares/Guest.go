package middlewares

import (
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"net/http"
)

// Guest 只允许未登录用户访问
func Guest(next http.HandlerFunc)http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		if auth.Check() {
			flash.Warning("登录用户无法访问次页面")
			http.Redirect(writer,request,"/",http.StatusFound)
		}

		next(writer,request)
	}
}
