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
		instance = NewProfileObject("result.txt")
	}
	return instance
}


func NewProfileObject(fileName string) *ProfileObject {
	return &ProfileObject{
		"/tmp/" + fileName,
		[]string{},
		0,
	}
}

func (po *ProfileObject) WriteFile() {
	content := ""
	for i:= 0; i < len(po.logs); i++ {
		content += po.logs[i]
	}
	err := ioutil.WriteFile(po.filePath, []byte(content), 0777)
	if err != nil {
		panic(err)
	}
}

func (po *ProfileObject) AddLogs(prefix string, postfix string) {
	if po.count > 100 {
		return
	} else if po.count == 100 {
		po.WriteFile()	
	}
	now := time.Now()
	timeStamp := now.UnixNano() / 1000000
	log := strings.Join([]string{prefix, strconv.FormatInt(timeStamp, 10), postfix}, " ") + "\n"
	po.logs = append(po.logs, log)
	po.count++
}