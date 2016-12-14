package controllers

import (
	"cron/services"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func Tasks(w http.ResponseWriter, r *http.Request) {
	list := services.Tasks()

	t, _ := template.ParseFiles("templates/tasks.html")
	tplValues := map[string]interface{}{}
	tplValues["tasks"] = list
	t.Execute(w, tplValues)
}

func TaskAdd(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	hash := md5.Sum(body)
	id := services.TaskTable + hex.EncodeToString(hash[:])
	fmt.Println("md5:" + id + "\n body:" + string(body))

	var task services.Task
	err := json.Unmarshal(body, &task)
	if err != nil {
		fmt.Println("参数异常")
		log.Fatal(err)
		panic(err)
	}
	log.Println(task)

	task.Save(id)
	task.Run(id)

	w.Write([]byte(`{"success":200}`))
}

func TaskDel(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["sid"]
	services.DeleteScheduler(sid)
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
