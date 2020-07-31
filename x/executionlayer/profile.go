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
	beforeTime int64
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

func (po *ProfileObject) AddLogs(prefix string, postfix string, simulate bool) {
	testCount := 100
	if po.count == testCount {
		po.WriteFile()
		po.count++
		return
	} else if po.count > testCount {
		return
	}

	command := ""
	if simulate == true {
		command += "CheckTx"
	} else {
		command += "DeliveTx"
	}

	now := time.Now()
	timeStamp := now.UnixNano() / 1000000
	log := strings.Join([]string{prefix, command, "TimeStamp", strconv.FormatInt(timeStamp, 10), postfix}, " ") + "\n"
	po.logs = append(po.logs, log)
	if strings.Contains(prefix, "Before") {
		po.beforeTime = timeStamp
		return
	}
	gap := timeStamp - po.beforeTime
	log = strings.Join([]string{prefix, command, "Gap", strconv.FormatInt(gap, 10), postfix}, " ") + "\n"
	po.logs = append(po.logs, log)
	if strings.Contains(prefix, "After Commit") {
		po.count++
	}
}
