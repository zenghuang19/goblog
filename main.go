package main

import (
	"github.com/gorilla/mux"
	"goblog/app/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"
)

var router *mux.Router

func init()  {
	config.Initialize()
}

func main() {
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	//err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	err := http.ListenAndServe("127.0.0.1:" + c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
