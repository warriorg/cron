package routers

import (
	"cron/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes 配置路由
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/tasks", controllers.Tasks).Methods("GET")
	router.HandleFunc("/task/add", controllers.TaskAdd).Methods("POST")
	router.HandleFunc("/task/del/{sid}", controllers.TaskDel).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", router)

	return router
}
