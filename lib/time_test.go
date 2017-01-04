package lib

import (
	"encoding/json"
	"log"
	"testing"
)

type JSONTest struct {
	Current Timestamp `json:"current"`
	Bone    string    `json:"bone"`
}

func (j *JSONTest) json() string {
	js, err := json.Marshal(j)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(js)
	return string(js)
}

func (j *JSONTest) fromJSON(value string) {
	log.Println("fromJSON", value)
	json.Unmarshal([]byte(value), &j)
}

// TaskFromJSON json 反序列化
func FromJSON(value string) (task JSONTest) {
	// log.Println("FromJSON", value, task.Current.Time == &time.Time{})
	json.Unmarshal([]byte(value), &task)
	return task
}

func (j *JSONTest) String() {
	log.Println("String ", j.Bone)
}

func Test_Make(t *testing.T) {
	list := make([]int, 5, 8)
	println(list)
	list = append(list, 1)
	println(list)
	list = append(list, 1)
	println(list)
	list = append(list, 1)
	println(list)
	list = append(list, 1)
	println(list)
	list = append(list, 1)
	println(list)

}

func Test_JSON(t *testing.T) {

	// jt := JSONTest{Current: Timestamp{time.Now()}, Bone: "bbbb"}
	// log.Println(jt.json())

	str := `{"current":"2016-09-21 10:10:10","bone":"aaaa"}`
	j := FromJSON(str)
	log.Println(j)
	log.Println(j.json())

	jx := JSONTest{}
	jx.fromJSON(str)
	log.Println(jx)
	log.Println(jx.json())
}
