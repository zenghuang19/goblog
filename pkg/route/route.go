package route

import "github.com/gorilla/mux"

//router 路由对象
var Router *mux.Router

//初始化
func Initialize()  {
	Router = mux.NewRouter()
}

func Name2URL(routName string,pairs ...string)string  {
	url,err := Router.Get(routName).URL(pairs...)

	if err != nil{
		return ""
	}

	return url.String()
}
