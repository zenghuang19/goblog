package middlewares

import (
	"goblog/pkg/session"
	"net/http"
)

// StartSession 开启 session 会话控制
func StartSession(next http.Handler)http.Handler  {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//1 启动会话
		session.StartSession(writer,request)

		// 2继续出来
		next.ServeHTTP(writer,request)
	})
}