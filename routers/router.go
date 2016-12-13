package routers

import (
	"cron/controllers"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/hello", controllers.Hello).Methods("GET")
	router.HandleFunc("/tasks", controllers.Tasks).Methods("GET")

	router.HandleFunc("/scheduler/add", controllers.SchedulerAdd).Methods("POST")
	router.HandleFunc("/scheduler/del/{sid}", controllers.SchedulerDel).Methods("GET")

	return router
}
