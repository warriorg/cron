package gocron

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Job struct {
	// pause interval * unit bettween runs
	interval uint64

	jobId string

	lastRun time.Time
	// datetime of next run
	nextRun time.Time
	// cache the period between last an next run
	period time.Duration

	runing bool

	unit string
	// Map for the function task store
	funcs interface{}
	// Map for function and  params of function
	fparams []interface{}
	delete  bool
}

func NewJob(id string, interval uint64, unit string, time time.Time) *Job {
	return &Job{
		interval: interval,
		jobId:    id,
		unit:     unit,
		nextRun:  time,
		runing:   false,
		delete:   false,
	}
}

func (j *Job) String() string {
	return fmt.Sprintf("jobId: %s interval: %d unit:%s lastRun:%s nextRun:%s period:%s",
		j.jobId, j.interval, j.unit, j.lastRun, j.nextRun, j.period)
}

func (j *Job) NextRun() time.Time {
	return j.nextRun
}

func (j *Job) LastRun() time.Time {
	return j.lastRun
}

func (j *Job) Delete() {
	j.delete = true
}

func (j *Job) shouldRun() bool {
	// log.Println("当前时间", time.Now())
	// log.Println("下次执行时间", j.nextRun)
	return time.Now().After(j.nextRun)
}

func (j *Job) run() {
	j.runing = true
	log.Println("run task")

	f := reflect.ValueOf(j.funcs)
	params := j.fparams
	if len(params)+1 != f.Type().NumIn() {
		log.Println("The number of param is not adapted.")
	}

	in := make([]reflect.Value, len(params)+1)
	in[0] = reflect.ValueOf(j)
	for k, param := range params {
		in[k+1] = reflect.ValueOf(param)
	}

	j.lastRun = time.Now()
	j.scheduleNextRun()
	f.Call(in)
	j.runing = false
}

func (j *Job) Do(jobFun interface{}, params ...interface{}) {
	typ := reflect.TypeOf(jobFun)
	if typ.Kind() != reflect.Func {
		panic("only function can be schedule into the job queue.")
	}

	j.funcs = jobFun
	j.fparams = params

	// j.scheduleNextRun()
}

//计算下次运行时间
//second, minute, hour, day, week, month, year
func (j *Job) scheduleNextRun() {
	if j.lastRun == time.Unix(0, 0) {
		j.lastRun = time.Now()
	}

	if j.interval == 0 {
		return
	}

	if j.period != 0 {
		// translate all the units to the Seconds
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	} else {
		switch strings.ToLower(j.unit) {
		case "second":
			j.period = time.Duration(j.interval)
			j.nextRun = j.lastRun.Add(j.period * time.Second)
			if j.nextRun.Before(time.Now()) {
				j.nextRun = time.Now()
			}
		case "minute":
			j.period = time.Duration(j.interval * 60)
			j.nextRun = j.lastRun.Add(j.period * time.Second)
			if j.nextRun.Before(time.Now()) {
				j.nextRun = time.Now()
			}
		case "hour":
			j.period = time.Duration(j.interval * 60 * 60)
			j.nextRun = j.lastRun.Add(j.period * time.Second)
			if j.nextRun.Before(time.Now()) {
				duration := time.Since(j.nextRun)
				j.nextRun = j.nextRun.Add(time.Duration(duration.Hours()) * time.Second)
			}
		case "day":
			j.nextRun = j.nextRun.AddDate(0, 0, 1)
			if j.nextRun.Before(time.Now()) {
				day := time.Now().Day() - j.nextRun.Day()
				j.nextRun = j.nextRun.AddDate(0, 0, day)
			}
		case "week":
			j.period = time.Duration(j.interval * 60 * 60 * 24 * 7)
			j.nextRun = j.lastRun.Add(j.period * time.Second)
			if j.nextRun.Before(time.Now()) {
				day := time.Now().Day() - j.nextRun.Day()
				j.nextRun = j.nextRun.AddDate(0, 0, int(math.Ceil(float64(day/7.0)))*7)
			}
		case "month":
			j.nextRun = j.nextRun.AddDate(0, 1, 0)
			if j.nextRun.Before(time.Now()) {
				month := int(time.Now().Month()) - int(j.nextRun.Month())
				j.nextRun = j.nextRun.AddDate(0, month, 0)
			}
		case "year":
			j.nextRun = j.nextRun.AddDate(1, 0, 0)
			if j.nextRun.Before(time.Now()) {
				year := time.Now().Year() - j.nextRun.Year()
				j.nextRun = j.nextRun.AddDate(year, 0, 0)
			}
		}

	}

}

type Scheduler struct {
	jobs []*Job
}

var instance *Scheduler
var once sync.Once

func GetScheduler() *Scheduler {
	once.Do(func() {
		instance = &Scheduler{[]*Job{}}
	})
	return instance
}

func (s *Scheduler) Add(j *Job) {
	if s.GetJob(j.jobId) != nil {
		log.Println("任务已存在<-->" + j.jobId)
		return
	}

	s.jobs = append(s.jobs, j)
}

func (s *Scheduler) Remove(jobId string) {
	index := s.Index(jobId)
	if index < 0 {
		return
	}
	if len(s.jobs) < index+1 {
		s.jobs = append(s.jobs[:index], s.jobs[index+1:]...)
	} else {
		s.jobs = s.jobs[:index]
	}
}

func (s *Scheduler) CleanDelJob() {
	jobs := make([]*Job, 0)
	for _, job := range s.jobs {
		if !job.delete {
			jobs = append(jobs, job)
		}
	}
	s.jobs = jobs
}

func (s *Scheduler) Index(jobId string) int {
	for index, job := range s.jobs {
		if job.jobId == jobId {
			return index
		}
	}
	return -1
}

func (s *Scheduler) GetJob(jobId string) (job *Job) {
	for _, job = range s.jobs {
		if job.jobId == jobId {
			return
		}
	}
	return nil
}

func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			runtime.Gosched()
			select {
			case <-ticker.C:
				s.RunPending()
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}

// Run all the jobs that are scheduled to run.
func (s *Scheduler) RunPending() {
	s.CleanDelJob()
	for _, j := range s.jobs {
		if !j.runing && j.shouldRun() {
			go j.run()
		}
	}
}

// for given function fn , get the name of funciton.
func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf((fn)).Pointer()).Name()
}
