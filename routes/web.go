package routes

import (
	"github.com/gorilla/mux"
	"goblog/app/http/controllers"
	"net/http"
)

// RegisterWebRoutes 注册网页相关路由
func RegisterWebRoutes(r *mux.Router) {
	// 静态页面
	pc := new(controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	//文章相关页面
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")

	//文章列表
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")

	//创建页面
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	//创建
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")

	//编辑回显
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")

	//更新内容
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")

	//删除
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	//静态资源
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 用户认证
	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register", auc.Register).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", auc.DoRegister).Methods("POST").Name("auth.doregister")

	//中间件：强制内容类型为 HTML
	//r.Use(middlewares.ForceHTML)
}
