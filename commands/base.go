package commands

import (
	"flag"
	"github.com/luoxiaojun1992/go-skeleton/bootstrap"
	"os"
)

type CommandInterface interface {
	ParseOptions(flag *flag.FlagSet)
	Handle()
}

type BaseCommand struct {
	App  *bootstrap.App
	Flag *flag.FlagSet
}

func (bc *BaseCommand) Run(command CommandInterface, app *bootstrap.App) {
	bc.Init(command, app).Handle()
}

func (bc *BaseCommand) Init(command CommandInterface, app *bootstrap.App) CommandInterface {
	bc.App = app
	if len(os.Args) > 1 {
		bc.Flag = app.NewFlag()
		command.ParseOptions(bc.Flag)
		app.ParseFlags(bc.Flag)
	}
	return command
}
