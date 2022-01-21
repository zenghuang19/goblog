package bootstrap

import (
	"embed"
	"github.com/gorilla/mux"
	"goblog/pkg/route"
	"goblog/routes"
	"io/fs"
	"net/http"
)

// SetupRoute 路由初始化
func SetupRoute(staticFS embed.FS) *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	route.SetRouter(router)

	//静态资源
	sub,_ := fs.Sub(staticFS, "public")
	router.PathPrefix("/").Handler(http.FileServer(http.FS(sub)))

	return router
}
