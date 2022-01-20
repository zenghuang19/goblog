package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/requests"
	"goblog/pkg/flash"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

type CategoriesController struct {
	BaseController
}

// Create 文章分类创建页面
func (*CategoriesController) Create(w http.ResponseWriter,r *http.Request)  {
	view.Render(w,view.D{}, "categories.create")
}

func (*CategoriesController) Store(w http.ResponseWriter,r *http.Request)  {
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}

	//表单验证
	errors := requests.ValidateCategoryForm(_category)

	//检查错误
	if len(errors) == 0 {
		//创建分类
		_category.Create()
		if _category.ID > 0 {
			flash.Success("分类创建成功")
			indexURL := route.Name2URL("home")
			http.Redirect(w, r, indexURL, http.StatusFound)
		}else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建分类失败，请联系管理员")
		}
	}else {
		view.Render(w, view.D{
			"Category": _category,
			"Errors" : errors,
		}, "categories.create")
	}
}

func (cc *CategoriesController) Show(w http.ResponseWriter, r *http.Request)  {
	// 获取URL参数
	id := route.GetRouteVariable("id", r)

	//读取对应参数
	_category,err := category.Get(id)

	//获取结果集
	articles, pagerData, err := article.GetByCategoryID(_category.GetStringID(), r, 2)

	if err != nil{
		cc.ResponseForSQLError(w,err)
	}else {
		// 加载模板
		view.Render(w,view.D{
			"Articles": articles,
			"pagerData":pagerData,
		}, "articles.index", "articles._article_meta")
	}
}