package data_types

import (
	"time"
)

func parseData(data time.Time) string {
	auxData := data.Format("2006-01-02")
	if auxData == "0001-01-01" {
		return ""
	}
	return auxData
}

func checkDateLayout(value string) string {
	var layout string
	if value[len(value)-1] == 'Z' {
		layout = "2006-01-02 15:04:05Z"
	} else {
		layout = "2006-01-02 15:04:05"
	}
	return layout
}
