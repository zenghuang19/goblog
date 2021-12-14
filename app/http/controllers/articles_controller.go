package controllers

import (
	"fmt"
	article2 "goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

type ArticlesController struct {
}

func (* ArticlesController) Show(w http.ResponseWriter,r *http.Request)  {
	//1.获取参数
	id := route.GetRouteVariable("id", r)

	//2.读取文字数据
	article, err := article2.Get(id)

	//3.出现错误
	if err != nil {
		//3.1未找到数据
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}

	} else {
		//4.读取成功
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Uint64ToString": types.Uint64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
	fmt.Fprintf(w, "ID"+id)
}

func (* ArticlesController) Index(w http.ResponseWriter,r *http.Request)  {
	// 1. 获取结果集
	articles, err := article2.GetAll()

	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 2. 加载模板
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		// 3. 渲染模板，将所有文章的数据传输进去
		err = tmpl.Execute(w, articles)
		logger.LogError(err)
	}
}