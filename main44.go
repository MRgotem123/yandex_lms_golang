package main

import (
	"time"
)

func FormatTimeToString(timestamp time.Time, format string) string {
	formattedTime := timestamp.Format(format)
	return formattedTime
}
