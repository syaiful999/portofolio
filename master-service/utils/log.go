package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
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
	fileName := "service.json"

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		trackError(logMessage, err)
		loc = time.UTC
	}
	today := time.Now().In(loc)
	data, err := os.ReadFile(fileName)
	loging := Logs{}
	if err == nil {
		err = json.Unmarshal(data, &loging)
		if err != nil {
			trackError(logMessage, err)
		}
	} else if !os.IsNotExist(err) {
		trackError(logMessage, err)
	}

	// List of sensitive words to mask in logs
	sensitiveWords := []string{"password", "token", "nik", "nik_number", "access_key", "secret_key"}

	// Combine messages for masking
	fullMessage := fmt.Sprintf("%s, %s", logMessage, errMessage)

	// Perform masking of sensitive data (case-insensitive)
	maskedMessage := fullMessage
	for _, word := range sensitiveWords {
		re := regexp.MustCompile(`(?i)(` + word + `[:= ]+)[^, ]+`)
		maskedMessage = re.ReplaceAllString(maskedMessage, "$1***MASKED***")
	}

	str := strings.ReplaceAll(maskedMessage, "\"", " ")

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
