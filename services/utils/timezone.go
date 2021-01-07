package utils

import (
	"github.com/luoxiaojun1992/go-skeleton/consts"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"time"
)

func DefaultTimezone() *time.Location {
	timeZoneName := consts.DEFAULT_TIMEZONE
	location, errLoadLocation := time.LoadLocation(timeZoneName)
	helper.CheckErrThenPanic("Failed to load time location ("+timeZoneName+")", errLoadLocation)

	return location
}
