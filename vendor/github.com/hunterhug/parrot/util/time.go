package util

import (
	"time"
)

func Sleep(wait_time int) {
	if wait_time <= 0 {
		return
	}
	time.Sleep(time.Duration(wait_time) * time.Second)
}

// Get Secord
func Second(times int) time.Duration {
	return time.Duration(times) * time.Second
}

// Get Timestamp
func GetTimestamp() int64 {
	return GetSecondTimes()
}

func GetSecondTimes() int64 {
	return time.Now().UTC().Unix()
}

//201611112113
func GetSecond2DateTimes(secord int64) string {
	tm := time.Unix(secord, 0)
	return tm.Format("20060102150405")

}

func GetDateTimes2Second(datestring string) int64 {
	tm2, _ := time.Parse("20060102150405", datestring)
	return tm2.Unix()

}
func TodayString(level int) string {
	formats := "20060102150405"
	switch level {
	case 1:
		formats = "2006"
	case 2:
		formats = "200601"
	case 3:
		formats = "20060102"
	case 4:
		formats = "2006010215"
	case 5:
		formats = "200601021504"
	default:

	}
	return time.Now().Format(formats)
}
