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
