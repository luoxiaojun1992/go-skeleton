package phabricator

import (
	"flag"
	"github.com/luoxiaojun1992/go-skeleton/commands"
	phabricatorLogic "github.com/luoxiaojun1992/go-skeleton/logics/phabricator"
	"log"
	"time"
)

type BugReporter struct {
	commands.BaseCommand
	OptionStartTime string
	OptionEndTime   string
}

func (br *BugReporter) Handle() {
	log.Println("Start...")

	(&phabricatorLogic.BugExporter{}).Export(br.OptionStartTime, br.OptionEndTime)

	log.Println("Finished.")
}

func (br *BugReporter) ParseOptions(flag *flag.FlagSet) {
	flag.StringVar(&br.OptionStartTime, "start", time.Now().Add(-15*time.Minute).Format("2006-01-02 15:04")+":00", "Start Time")
	flag.StringVar(&br.OptionEndTime, "end", time.Now().Add(-1*time.Minute).Format("2006-01-02 15:04")+":59", "End Time")
}
