package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Logs struct {
	Logs []LogsData `form:"logs" json:"logs" xml:"logs"`
}

type LogsData struct {
	Date string `form:"date" json:"date" xml:"date"`
	Log  string `form:"log" json:"log" xml:"log"`
}

func PushLogf(timezone, logMessage, errMessage string) {
	/**

	 */
	logDir := "./logs"
	currentDate := time.Now().Format("20060102")
	fileName := logDir + "/" + currentDate + ".json"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		trackError(logMessage, err)

	}
	today := time.Now().In(loc)
	data, err := os.ReadFile(fileName)
	if err != nil {
		trackError(logMessage, err)

	}

	loging := Logs{}
	err = json.Unmarshal(data, &loging)
	if err != nil {
		trackError(logMessage, err)

	}

	str := strings.ReplaceAll(fmt.Sprintf("[%s], %s", logMessage, errMessage), "\"", " ")

	loging.Logs = append(loging.Logs, LogsData{
		Date: today.Format("2006-01-02 15:04:05"),
		Log:  strings.ReplaceAll(str, "\\", ""),
	})

	bytes, _ := json.Marshal(loging)

	os.WriteFile(fileName, bytes, 0644)
}

func trackError(logMessage string, err error) {
	if err != nil {
		log.Printf("[log err] %s|%s", logMessage, err.Error())

	}
}
