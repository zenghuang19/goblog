package main

import (
	"embed"
	"github.com/gorilla/mux"
	middlewares2 "goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"
)

var router *mux.Router

//go:embed resources/views/articles/*
//go:embed resources/views/auth/*
//go:embed resources/views/categories/*
//go:embed resources/views/layouts/*
var tplFS embed.FS

//go:embed public/*
var staticFS embed.FS

func init()  {
	config.Initialize()
}

func main() {
	//初始化SQL
	bootstrap.SetupDB()

	//初始化模板
	bootstrap.SetupTemplate(tplFS)

	// 初始化路由绑定
	router = bootstrap.SetupRoute(staticFS)

	//err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	err := http.ListenAndServe("127.0.0.1:" + c.GetString("app.port"), middlewares2.RemoveTrailingSlash(router))
	logger.LogError(err)
}
