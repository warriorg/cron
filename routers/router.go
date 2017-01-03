package routers

import (
	"cron/controllers"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/hello", controllers.Hello).Methods("GET")
	router.HandleFunc("/tasks", controllers.Tasks).Methods("GET")
	router.HandleFunc("/task/add", controllers.TaskAdd).Methods("POST")
	router.HandleFunc("/task/del/{sid}", controllers.TaskDel).Methods("GET")

	return router
}
