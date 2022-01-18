package controllers

import (
	"fmt"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"gorm.io/gorm"
	"net/http"
)

type BaseController struct {
}

// ResponseForSQLError 处理 SQL 错误并返回
func (bc BaseController) ResponseForSQLError(w http.ResponseWriter,err error)  {
	if err == gorm.ErrRecordNotFound {
		// 未找到数据
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 文章未找到")
	}else {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	}
}

func (bc BaseController) ResponseForUnauthorized(w http.ResponseWriter, r *http.Request)  {
	flash.Warning("未授权操作！")
	http.Redirect(w, r, "/", http.StatusFound)
}