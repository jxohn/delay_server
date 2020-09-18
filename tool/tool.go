package tool

import "time"

func BuildFileNameByTime(time time.Time) string {
	return time.Format("2006-01-02-15")
}
