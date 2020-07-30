package executionlayer

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type ProfileObject struct {
	filePath string
	logs []string
	count int
}

var instance *ProfileObject

func GetInstance() *ProfileObject {
	if instance == nil {
		instance = NewProfileObject("/tmp")
	}
	return instance
}


func NewProfileObject(filePath string) *ProfileObject {
	return &ProfileObject{
		"/tmp/" + filePath,
		[]string{},
		0,
	}
}

func (po *ProfileObject) WriteFile() {
	content := ""
	for i:= 0; i < len(po.logs); i++ {
		content += po.logs[i]
	}
	err := ioutil.WriteFile("/tmp/result.txt", []byte(content), 0)
	if err != nil {
		panic(err)
	}
}

func (po *ProfileObject) AddLogs(prefix string, postfix string) {
	if po.count > 1000 {
		return
	} else if po.count == 1000 {
		po.WriteFile()	
	}
	now := time.Now()
	timeStamp := now.UnixNano() / 1000000
	log := strings.Join([]string{prefix, strconv.FormatInt(timeStamp, 10), postfix}, " ") + "\n"
	po.logs = append(po.logs, log)
	po.count++
}
