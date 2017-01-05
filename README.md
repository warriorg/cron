# cron

run app

go run main.go

tasks list
>http://127.0.0.1:7000/tasks

put task
`{"time":"2017-01-04 12:42:35","maxCount":0,"every":0,"url":"http://127.0.0.1:9000/api/v1/home/hello","body":"","method":"GET"}`

time
>"2006-01-02 15:04"

unit type
>minute, hour, day, week, month, year


next:
  增加回调消息错误时自动重试功能
  增加最多调用册书限制
  集成日志
