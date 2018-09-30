package common

import (
	"strconv"
)

func ConvertF32ToString(timePush float32) string {
	var timePushInt = int(timePush)
	if timePush > float32(timePushInt) {
		return strconv.Itoa(timePushInt) + "h 30p"
	}
	return strconv.Itoa(timePushInt) + "h"
}
