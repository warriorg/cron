package controllers

import (
	"cron/services"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Index 检查
func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t = t.Delims("<<:", ":>>")
	t.Execute(w, nil)
}

// Tasks 任务列表
func Tasks(w http.ResponseWriter, r *http.Request) {
	list := services.Tasks()

	switch r.Header.Get("Accept") {
	case "application/json":
		jb, err := json.Marshal(list)
		if err != nil {
			log.Println(err)
		}
		w.Write(jb)
	default:
		t, _ := template.ParseFiles("templates/tasks.html")
		tplValues := map[string]interface{}{}
		tplValues["tasks"] = list
		t.Execute(w, tplValues)
	}

}

// TaskAdd 任务添加
func TaskAdd(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	err := services.TaskAdd(body)

	if err != nil {
		log.Println("参数异常", err.Error())
		w.Write([]byte(`{"success":1}`))
		return
	}

	w.Write([]byte(`{"success":0}`))
}

// TaskDel 删除
func TaskDel(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["sid"]
	services.TaskDelete(sid)
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Write([]byte(`{"success":0}`))
	default:
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	}

}
