package main

import (
	"time"
)

func TimeAgo(pastTime time.Time) string {
	diff := pastTime.Sub(time.now)
	formattedTime := diff.Format(string)
	return formattedTime
}
