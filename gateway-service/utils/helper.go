package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
)

func FindInArray(needle interface{}, hystack interface{}) (int, bool) {
	switch key := needle.(type) {
	case string:
		for i, item := range hystack.([]string) {
			if key == item {
				return i, true
			}
		}
	case int:
		for i, item := range hystack.([]int) {
			if key == item {
				return i, true
			}
		}
	case int64:
		for i, item := range hystack.([]int64) {
			if key == item {
				return i, true
			}
		}
	default:
		return -1, false
	}
	return -1, false
}

func ReadJsonFile(file *multipart.FileHeader) (string, error) {
	var result string
	read, errOpen := file.Open()
	defer read.Close()

	buf := make([]byte, file.Size)

	len, _ := read.Read(buf)
	b := string(buf)
	if errOpen != nil {
		return result, errOpen
	}

	if len == 0 {
		return result, errors.New("invalid json file")
	}

	return b, nil
}

func ReadXMlFile(file *multipart.FileHeader) (string, error) {
	var result string
	read, errOpen := file.Open()
	defer read.Close()

	buf := make([]byte, file.Size)

	len, _ := read.Read(buf)
	b := string(buf)
	if errOpen != nil {
		return result, errOpen
	}

	if len == 0 {
		return result, errors.New("invalid xml file")
	}

	return b, nil
}

func TrimLastChar(data string) string {
	if len(data) == 0 {
		return ""
	}
	data = data[0 : len(data)-1]
	return data
}

func ArrToStrDelimiter(arrs []string, delimiter string) string {
	arrCheck := []string{}
	result := ""
	if delimiter == "" {
		delimiter = ","
	}
	for _, arr := range arrs {
		if _, err := FindInArray(arrCheck, arr); !err {
			arrCheck = append(arrCheck, arr)
			result += "'" + arr + "'" + delimiter
		}
	}
	return strings.TrimSuffix(result, delimiter)
}

func ConvStringToInt(data string) int {
	a, b := strconv.Atoi(data)
	if b != nil {
		return 0
	}
	return a
}

func ConvStringToFloat64(data string) float64 {
	a, b := strconv.ParseFloat(data, 64)
	if b != nil {
		return 0
	}
	return a
}

func CheckRubbishString(data string, rubbish string) string {
	pText := strings.Trim(data, rubbish)
	pTextCheck := strings.Trim(pText, " ")
	if len(pTextCheck) == 0 {
		return ""
	}
	return pText
}

func Percentage(a int, b int) float64 {
	if b == 0 {
		b = 1
	}
	c := fmt.Sprintf("%.2f", float64(a)/float64(b)*100)
	d, _ := strconv.ParseFloat(c, 64)
	return d
}
