package utils

import (
	"testing"
	"time"
)

func getTimeNowVietNam() time.Time {
	var LocNowVietNam, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(LocNowVietNam)
}
func TestBeginningOfMinute(t *testing.T) {

	var now = Now{
		Time: getTimeNowVietNam(),
	}
	if m1, m2 := now.GetCurrentMonth(); m1 == 1 && m2 == 2 {
		t.Log("ok")
	} else {
		t.Log(now.GetCurrentDay())
		t.Fail()
	}
}
