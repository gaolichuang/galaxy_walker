package time_util

import (
	"time"
        "fmt"
)

func GetCurrentTimeStamp() int64 {
	return time.Now().Unix()
}
func GetTimeInMs() int64 {
	return time.Now().UnixNano() / int64(1000000)
}
func GetReadableTimeNow() string {
	return GetReadableTime(GetCurrentTimeStamp())
}
func GetReadableTimeNowS() string {
        y,m,d := time.Now().Date()
        time.Now().Clock()
        return fmt.Sprintf("%.4d%.2d%.2d",y,int(m),d)
}
func GetCurrentTime() string {
        y,m,d := time.Now().Date()
        h,mi,s := time.Now().Clock()
        return fmt.Sprintf("%.4d%02d%02d%02d%02d%02d",y,int(m),d,h,mi,s)
}
func GetCurrentDate() string {
        y,m,d := time.Now().Date()
        return fmt.Sprintf("%.4d%.2d%.2d",y,int(m),d)
}
// input is timestamp
func GetReadableTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04 MST")
}
func Sleep(s int) {
	time.Sleep(time.Second * time.Duration(s))
}
func Usleep(n int) {
	time.Sleep(time.Microsecond * time.Duration(n))
}
func GetCurrentTimeInNumSeries() string {
	return time.Unix(time.Now().Unix(), 0).Format("200601021504")
}
