package controllers

import (
	"cron/services"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

	hash := md5.Sum(body)
	id := hex.EncodeToString(hash[:])
	log.Println("md5:" + id + " body:" + string(body))

	var task *services.Task
	err := json.Unmarshal(body, &task)
	if err != nil {
		log.Fatal("参数异常", err.Error())
		w.Write([]byte(`{"success":1}`))
		return
	}

	if task == nil {
		w.Write([]byte(`{"success":1}`))
		return
	}

	task.Save(id)
	task.Run(id)

	w.Write([]byte(`{"success":0}`))
}

func TaskDel(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["sid"]
	services.DeleteScheduler(sid)
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
