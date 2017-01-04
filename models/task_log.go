package models

import (
	"cron/lib"
	"time"
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
	Header map[string]string `json:"header"`
}

// TaskLog 工作日志
type TaskWork struct {
	Count   int64     `json:"count"`
	LastRun time.Time `json:"lastRun"`
}
