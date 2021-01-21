package phabricator

import (
	"flag"
	"github.com/luoxiaojun1992/go-skeleton/commands"
	phabricatorLogic "github.com/luoxiaojun1992/go-skeleton/logics/phabricator"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"github.com/luoxiaojun1992/go-skeleton/services/utils/time/timezone"
	"log"
	"time"
)

type BugReporter struct {
	commands.BaseCommand
	OptionStartTime string
	OptionEndTime   string
	Debug           bool
}

func (br *BugReporter) validOptions() {
	location := timezone.Timezone()

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

func (br *BugReporter) handleWithDebug(handler func()) {
	log.Println("Start...")
	log.Println("Task Start Time: " + br.OptionStartTime)
	log.Println("Task End Time: " + br.OptionEndTime)

	handler()

	log.Println("Finished.")
}

func (br *BugReporter) Handle() {
	handler := func() {
		br.validOptions()
		(&phabricatorLogic.BugExporter{}).Export(br.OptionStartTime, br.OptionEndTime)
	}

	if br.Debug {
		br.handleWithDebug(handler)
	} else {
		handler()
	}
}

func (br *BugReporter) ParseOptions(flag *flag.FlagSet) {
	location := timezone.Timezone()
	flag.StringVar(&br.OptionStartTime, "start", time.Now().Add(-15*time.Minute).In(location).Format("2006-01-02 15:04")+":00", "Start Time")
	flag.StringVar(&br.OptionEndTime, "end", time.Now().Add(-1*time.Minute).In(location).Format("2006-01-02 15:04")+":59", "End Time")
	flag.BoolVar(&br.Debug, "debug", false, "Debug")
}
