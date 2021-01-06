package commands

import (
	"flag"
	"github.com/luoxiaojun1992/go-skeleton/bootstrap/command"
	"os"
)

type CommandInterface interface {
	ParseOptions(flag *flag.FlagSet)
	Handle()
}

type BaseCommand struct {
	App  *command.App
	Flag *flag.FlagSet
}

func (bc *BaseCommand) Run(command CommandInterface, app *command.App) {
	bc.Init(command, app).Handle()
}

func (bc *BaseCommand) Init(command CommandInterface, app *command.App) CommandInterface {
	bc.App = app
	if len(os.Args) > 1 {
		bc.Flag = app.NewFlag()
		command.ParseOptions(bc.Flag)
		app.ParseFlags(bc.Flag)
	}
	return command
}
