package main

import (
	"github.com/gorilla/mux"
	"goblog/app/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/logger"
	"net/http"
)

var router *mux.Router

func main() {
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	//err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	err := http.ListenAndServe("127.0.0.1:3000", middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
