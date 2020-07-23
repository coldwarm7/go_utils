package convert

import "time"

func TimestampToDate(timestamp int64) string {
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
