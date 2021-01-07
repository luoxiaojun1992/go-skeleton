package phabricator

import (
	"flag"
	"github.com/luoxiaojun1992/go-skeleton/commands"
	phabricatorLogic "github.com/luoxiaojun1992/go-skeleton/logics/phabricator"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"github.com/luoxiaojun1992/go-skeleton/services/utils"
	"log"
	"time"
)

type BugReporter struct {
	commands.BaseCommand
	OptionStartTime string
	OptionEndTime   string
}

func (br *BugReporter) validOptions() {
	location := utils.DefaultTimezone()

	modifiedStart, errModifiedStart := time.ParseInLocation("2006-01-02 15:04:05", br.OptionStartTime, location)
	helper.CheckErrThenPanic("Failed to parse task modified start", errModifiedStart)
	if modifiedStart.Format("2006-01-02 15:04:05") != br.OptionStartTime {
		log.Panicln("Invalid format of task modified start")
	}

	modifiedEnd, errModifiedEnd := time.ParseInLocation("2006-01-02 15:04:05", br.OptionEndTime, location)
	helper.CheckErrThenPanic("failed to parse task modified end", errModifiedEnd)
	if modifiedEnd.Format("2006-01-02 15:04:05") != br.OptionEndTime {
		log.Panicln("Invalid format of task modified end")
	}
}

func (br *BugReporter) Handle() {
	log.Println("Start...")

	br.validOptions()
	(&phabricatorLogic.BugExporter{}).Export(br.OptionStartTime, br.OptionEndTime)

	log.Println("Finished.")
}

func (br *BugReporter) ParseOptions(flag *flag.FlagSet) {
	flag.StringVar(&br.OptionStartTime, "start", time.Now().Add(-15*time.Minute).Format("2006-01-02 15:04")+":00", "Start Time")
	flag.StringVar(&br.OptionEndTime, "end", time.Now().Add(-1*time.Minute).Format("2006-01-02 15:04")+":59", "End Time")
}
