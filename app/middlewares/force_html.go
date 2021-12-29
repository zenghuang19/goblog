package middlewares

import "net/http"

func ForceHTML(next http.Handler)http.Handler  {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//1.设置标头
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		//2.继续出来请求
		next.ServeHTTP(writer,request)
	})
}