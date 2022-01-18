package middlewares

import (
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"net/http"
)

// Auth 登录用户才可访问
func Auth(next http.HandlerFunc)http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !auth.Check() {
			flash.Warning("登录用户才能访问此页面")
			http.Redirect(writer,request,"/", http.StatusFound)
		}

		next(writer,request)
	}
}
