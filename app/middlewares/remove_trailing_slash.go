package middlewares

import (
	"net/http"
	"strings"
)

// RemoveTrailingSlash 除首页以外，移除所有请求路径后面的斜杆
func RemoveTrailingSlash(next http.Handler)http.Handler  {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 1.除首页以外，移除所有请求路径后面的斜杠
		if request.URL.Path != "/" {
			request.URL.Path = strings.TrimSuffix(request.URL.Path, "/")
		}

		//2.将请求传递下去
		next.ServeHTTP(writer,request)
	})
}