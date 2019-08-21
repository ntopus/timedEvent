package data_types

import (
	"time"
)

func ParseData(data time.Time) string {
	auxData := data.Format("2006-01-02")
	if auxData == "0001-01-01" {
		return ""
	}
	return auxData
}

func GetTime(data string) (*time.Time, error) {
	const shortForm = "2006-01-02 15:04:05Z"
	t, err := time.Parse(shortForm, data)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CheckDateLayout(value string) string {
	var layout string
	if value[len(value)-1] == 'Z' {
		layout = "2006-01-02 15:04:05Z"
	} else {
		layout = "2006-01-02 15:04:05"
	}
	return layout
}
