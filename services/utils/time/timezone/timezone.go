package timezone

import (
	"github.com/gookit/config/v2"
	timezoneConsts "github.com/luoxiaojun1992/go-skeleton/consts/time/timezone"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"time"
)

var LocationName string

var Loc *time.Location

func Setup() {
	LocationName := config.String("app.location", timezoneConsts.DEFAULT_TIMEZONE)
	Loc = GetTimezone(LocationName)
}

func Timezone() *time.Location {
	if Loc == nil {
		return time.Local
	}

	return Loc
}

func TimezoneName() string {
	if len(LocationName) <= 0 {
		return timezoneConsts.DEFAULT_TIMEZONE
	}

	return LocationName
}

func GetTimezone(timezoneName string) *time.Location {
	loc, errLoadLoc := time.LoadLocation(timezoneName)
	helper.CheckErrThenPanic("Failed to load time location ("+timezoneName+")", errLoadLoc)
	return loc
}
