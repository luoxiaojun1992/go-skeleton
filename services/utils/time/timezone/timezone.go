package timezone

import (
	timezoneConsts "github.com/luoxiaojun1992/go-skeleton/consts/time/timezone"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"time"
)

func DefaultTimezone() *time.Location {
	timeZoneName := timezoneConsts.DEFAULT_TIMEZONE
	location, errLoadLocation := time.LoadLocation(timeZoneName)
	helper.CheckErrThenPanic("Failed to load time location ("+timeZoneName+")", errLoadLocation)

	return location
}
