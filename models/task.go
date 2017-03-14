package models

import (
	"cron/lib"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CTasks         = "tasks"
	CTaskHistories = "task_histories"
)

// Task 精度支持到分钟
type Task struct {
	Id string `bson:"_id" json:"id"`
	//2012-06-12 12:22
	Time lib.Timestamp `json:"time"`
	//2012-06-12 12:22
	EndTime lib.Timestamp `json:"endTime"`
	// 计数器
	MaxCount uint64 `json:"maxCount"`

	Every uint64 `json:"every"`
	//second, Minute, Hour, Day, Week, Month, Year
	// "" 代表按照Time 执行的一次性任务
	Unit string `json:"unit"`
	// 回调的地址
	URL string `json:"url"`
	// 回调的方式
	Method string `json:"method"`
	// 回传的数据
	Body string `json:"body"`
	// Header 中回传的数据
	Header map[string]string `json:"header"`

	Count int64 `json:"count"`

	LastRun time.Time `json:"lastRun"`

	RunResult string `runResult`
}

func NewTaskById(id string) *Task {
	return &Task{Id: id}
}

func FindById(id string) (*Task, error) {
	var task *Task
	err := db.C(CTasks).FindId(id).One(&task)
	return task, err
}

// TaskList 列表
func (task *Task) All() ([]*Task, error) {
	list := make([]*Task, 5)

	err := db.C(CTasks).Find(nil).All(&list)
	return list, err
}

func (task *Task) FindById() (*Task, error) {
	err := db.C(CTasks).FindId(task.Id).One(&task)
	return task, err
}

func (task *Task) Save() error {
	if task == nil {
		return errors.New("task 不能为空")
	}

	err := db.C(CTasks).Insert(task)
	return err
}

func (task *Task) Update() error {
	if task == nil {
		return errors.New("task 不能为空")
	}
	err := db.C(CTasks).UpdateId(task.Id, task)
	return err
}

func (task *Task) Delete() error {
	if task == nil {
		return errors.New("task 不能为空")
	}
	err := db.C(CTasks).RemoveId(task.Id)
	return err
}

func (task *Task) Callback() (err error) {
	if task.URL == "" {
		log.Println("URL nil，stop callback")
		return
	}

	if task.Method == "" {
		task.Method = "GET"
	}

	client := &http.Client{}
	req, _ := http.NewRequest(task.Method, task.URL, nil)
	for key, value := range task.Header {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("调用 ", task.URL, " 返回：", string(body))
		return err
	}

	if result["success"].(float64) != 0 {
		log.Println("---------" + result["message"].(string))
		return errors.New(result["message"].(string))
	}
	return nil
}

func SaveHistory(task *Task) {
	hist := &TaskHistory{}
	hist.Task = task
	hist.RunTime = time.Now()
	err := db.C(CTaskHistories).Insert(hist)
	CheckErr(err)
}

func (task *Task) String() string {
	value, _ := json.Marshal(task)
	return string(value)
}

type TaskHistory struct {
	*Task
	Id      bson.ObjectId `bson:"_id,omitempty"`
	RunTime time.Time
}
