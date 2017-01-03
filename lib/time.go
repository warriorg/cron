package lib

import (
	"log"
	"time"
)

const (
	DATE_FORMAT = `"2006-01-02 15:04:05"`
)

type Timestamp struct {
	time.Time
}

// UnmarshalJSON 反序列化
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ct, err := time.ParseInLocation(DATE_FORMAT, string(b[:]), time.Local)
	if err != nil {
		log.Println("解析时间错误", err)
		return err
	}
	*t = Timestamp{ct}

	return err
}

// MarshalJSON to json
func (t Timestamp) MarshalJSON() ([]byte, error) {
	// log.Println("序列化 ：", t.Format(DATE_FORMAT))
	// if t.Time == time.Time{} {
	// 	return nil, nil
	// }
	return []byte(t.Format(DATE_FORMAT)), nil
}
