package services

import (
	"cron/gocron"
	"cron/lib"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// TODO id md5码 key 前置 + id

// 常量
const (
	TaskTable        = "task-"
	TaskLogTable     = "tasklog-"
	TaskHistoryTable = "taskhistory-"
)

// Task 精度支持到分钟
type Task struct {
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
	Header map[string]string
}

// TaskLog 工作日志
type TaskLog struct {
	Count   int64     `json:"count"`
	LastRun time.Time `json:"lastRun"`
}

// DbTask 数据库存储的
type DbTask struct {
	Key     string
	Task    *Task
	TaskLog *TaskLog
}

var db *leveldb.DB

func init() {
	_db, err := leveldb.OpenFile("data/db", nil)

	if err != nil {
		panic(err)
	}
	db = _db

	s := gocron.GetScheduler()
	s.Start()
	initJob()
}

func initJob() {
	iter := db.NewIterator(util.BytesPrefix([]byte(TaskTable)), nil)

	for iter.Next() {
		log.Println("load task key: %s, %s", string(iter.Key()[:]), string(iter.Value()[:]))
		key := iter.Key()
		task := TaskFromJSON(iter.Value())

		task.Run(strings.TrimLeft(string(key[:]), TaskTable))
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal(err)
	}

}

// Tasks 列表
func Tasks() []*DbTask {
	iter := db.NewIterator(util.BytesPrefix([]byte(TaskTable)), nil)
	list := []*DbTask{}
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		dbTask := new(DbTask)
		dbTask.Key = string(key[:])
		dbTask.Task = TaskFromJSON(value)
		list = append(list, dbTask)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal(err)
	}
	return list
}

// DeleteScheduler 删除job
func DeleteScheduler(key string) {
	db.Delete([]byte(key), nil)
	s := gocron.GetScheduler()
	s.Remove(key)
}

func getTaskByID(id string) (task *Task) {
	value, _ := db.Get([]byte(TaskTable+id), nil)
	return TaskFromJSON(value)
}

// Delete 删除
func (task *Task) Delete(id string) {
	db.Delete([]byte(TaskTable+id), nil)
}

// Save 保存任务
func (task *Task) Save(id string) {
	log.Println("保存", task.json())
	db.Put([]byte(TaskTable+id), []byte(task.json()), nil)
}

// SaveHistory 保存历史
func (task *Task) SaveHistory(id string) {
	db.Put([]byte(TaskHistoryTable+id), []byte(task.json()), nil)
}

// Run 执行任务
func (task *Task) Run(id string) error {
	log.Println("加入任务-->" + task.json())
	j := gocron.NewJob(id, task.Every, task.Unit, task.Time.Time)
	j.Do(taskRun, id)

	s := gocron.GetScheduler()
	s.Add(j)

	return nil
}

func (task *Task) callback() (err error) {
	if task.URL == "" {
		return errors.New("URL 不能为空")
	}

	if task.Method == "" {
		task.Method = "GET"
	}

	client := &http.Client{}
	req, _ := http.NewRequest(task.Method, task.URL, nil)
	for key, value := range task.Header {
		req.Header.Add(key, value)
	}

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result["success"].(float64) != 0 || result["data"] == nil {
		log.Println("---------" + result["message"].(string))
		return errors.New(result["message"].(string))
	}
	return nil
}

// json 转化为Json
func (task *Task) json() string {
	js, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	return string(js)
}

func taskRun(j *gocron.Job, id string) {
	fmt.Println(" Run task : ", time.Now(), id)

	task := getTaskByID(id)
	if task == nil {
		fmt.Println("Task nil, Remove task : " + id)
		s := gocron.GetScheduler()
		s.Remove(id)
		return
	}

	// task.Time = j.NextRun().Format(DATE_FORMAT)
	nextRun := j.NextRun()
	task.Time = lib.Timestamp{nextRun}
	task.Save(id)
	task.callback()

	logJSON, _ := db.Get([]byte(TaskLogTable+id), nil)
	var taskLog *TaskLog
	if logJSON == nil {
		taskLog = &TaskLog{}
	}
	_ = json.Unmarshal(logJSON, &taskLog)
	taskLog.Count++
	taskLog.LastRun = j.LastRun()
	taskLog.save(TaskLogTable + id)

	if (task.EndTime.After(time.Unix(0, 0)) && time.Now().After(task.EndTime.Time)) || task.Unit == "" {
		fmt.Println("Remove task : ", id, task.Time, task.EndTime)
		s := gocron.GetScheduler()
		s.Remove(id)
		task.SaveHistory(id)
		task.Delete(id)
	}
}

// TaskFromJSON json 反序列化
func TaskFromJSON(value []byte) (task *Task) {
	json.Unmarshal(value, &task)
	return
}

func (taskLog *TaskLog) save(id string) {
	db.Put([]byte(TaskLogTable+id), []byte(taskLog.json()), nil)
}

// json 转化为Json
func (taskLog *TaskLog) json() string {
	js, err := json.Marshal(taskLog)
	if err != nil {
		log.Fatal(err)
	}

	return string(js)
}
