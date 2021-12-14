package controllers

import (
	"fmt"
	"net/http"
)

// PagesController 处理静态页面
type PagesController struct {
}

// Home 首页
func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>欢迎来到我的世界！</h1>")
}
// About 关于我们页面
func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "唯有生活不可辜负！")
}

// NotFound 404 页面
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们1,。</p>")
}
