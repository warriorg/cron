package services

import (
	"cron/gocron"
	"cron/lib"
	"cron/models"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

// Task 精度支持到分钟
type TaskTo struct {
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
}

func (to *TaskTo) toTask(md5 string) *models.Task {
	task := &models.Task{}
	task.Id = md5
	task.Time = to.Time
	task.EndTime = to.EndTime
	task.MaxCount = to.MaxCount
	task.Every = to.Every
	task.Unit = to.Unit
	task.URL = to.URL
	task.Method = to.Method
	task.Body = to.Body
	task.Header = to.Header

	return task
}

func init() {
	s := gocron.GetScheduler()
	s.Start()
	initJob()
}

func initJob() {
	list := Tasks()
	for _, task := range list {
		joinTask(task)
	}
}

func TaskAdd(body []byte) error {
	hash := md5.Sum(body)
	md5 := hex.EncodeToString(hash[:])
	log.Println("md5:" + md5 + " body:" + string(body))
	var to *TaskTo
	err := json.Unmarshal(body, &to)
	if err != nil {
		return err
	}
	task := to.toTask(md5)
	task.Save()

	joinTask(task)
	return nil
}

// Tasks 列表
func Tasks() []*models.Task {
	var task *models.Task
	list, err := task.All()
	if err != nil {
		log.Println(err)
	}
	return list
}

// DeleteScheduler 删除job
func TaskDelete(id string) {
	task := models.NewTaskById(id)
	task.Delete()
	s := gocron.GetScheduler()
	s.DeleteJobByID(id)
}

func joinTask(task *models.Task) error {
	log.Println("join task --> " + task.Id)
	j := gocron.NewJob(task.Id, task.Every, task.Unit, task.Time.Time)
	j.Do(taskRun, task.Id)

	s := gocron.GetScheduler()
	s.Add(j)

	return nil
}

func taskRun(j *gocron.Job, id string) {
	log.Println("run task: ", time.Now(), id)

	task, err := models.FindById(id)
	if task == nil {
		log.Println("task nil, remove job : ", j)
		j.Delete()
		return
	}

	time.Sleep(20 * time.Second)

	// task.Time = j.NextRun().Format(DATE_FORMAT)
	nextRun := j.NextRun()
	task.Time = lib.Timestamp{nextRun}
	task.LastRun = j.LastRun()
	task.Count++
	task.RunResult = "success"
	err = task.Callback()
	if err != nil {
		task.RunResult = err.Error()
		//回调失败
		log.Println("回调错误：", err.Error(), "任务：", task)
	}

	models.SaveHistory(task)
	if (task.EndTime.After(time.Unix(0, 0)) && time.Now().After(task.EndTime.Time)) || task.Unit == "" {
		log.Println("remove task : ", task)
		j.Delete()
		task.Delete()
	} else {
		task.Update()
	}

}
