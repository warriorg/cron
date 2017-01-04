package controllers

import (
	"cron/services"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Hello 检查
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

// Tasks 任务列表
func Tasks(w http.ResponseWriter, r *http.Request) {
	list := services.Tasks()

	t, _ := template.ParseFiles("templates/tasks.html")
	tplValues := map[string]interface{}{}
	tplValues["tasks"] = list
	t.Execute(w, tplValues)
}

// TaskAdd 任务添加
func TaskAdd(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	err := services.TaskAdd(body)

	if err != nil {
		log.Fatal("参数异常", err.Error())
		w.Write([]byte(`{"success":1}`))
		return
	}

	w.Write([]byte(`{"success":0}`))
}

func TaskDel(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["sid"]
	services.TaskDelete(sid)
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
