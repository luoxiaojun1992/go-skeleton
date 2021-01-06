package main

import (
	"flag"
	"github.com/luoxiaojun1992/go-skeleton/bootstrap"
	"github.com/luoxiaojun1992/go-skeleton/commands/phabricator"
	"log"
)

func main() {
	bootstrap.Create().Run(func(app *bootstrap.App) {
		var handlers map[string]func()
		handlers = make(map[string]func())

		handlers["phabricator_bug_exporter"] = func() {
			bugReporter := &phabricator.BugReporter{}
			bugReporter.Run(bugReporter, app)
		}
		handlers["help"] = func() {
			flag.PrintDefaults()
		}

		if handler, hasHandler := handlers[app.CommandName]; hasHandler {
			handler()
		} else {
			log.Panic("Handler not found")
		}
	}, func(app *bootstrap.App) {
		var configures map[string]func()
		configures = make(map[string]func())

		configures["phabricator_bug_exporter"] = func() {
			(&phabricator.BugReporter{}).ParseOptions(app.Flag)
		}

		if configure, hasConfigure := configures[app.CommandName]; hasConfigure {
			configure()
		} else {
			log.Panic("Configure not found")
		}
	})
}
