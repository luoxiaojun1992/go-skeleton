package bootstrap

import (
	"flag"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/ini/v2/dotnev"
	"github.com/luoxiaojun1992/go-skeleton/services/db/sql/mysql"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"github.com/luoxiaojun1992/go-skeleton/services/phabricator"
	"github.com/luoxiaojun1992/go-skeleton/services/utils/time/timezone"
	"os"
)

type App struct {
	ConfigPath  string
	EnvDir      string
	CommandName string

	Flag *flag.FlagSet
}

func Create() *App {
	return &App{}
}

func (app *App) Run(router func(app *App), configure func(app *App)) {
	defer app.Shutdown()
	app.InitFramework(configure)
	router(app)
}

func (app *App) InitArgs(configure func(app *App)) {
	if len(os.Args) > 1 {
		app.CommandName = os.Args[1]
		app.Flag = app.NewFlag()
		configure(app)
		app.ParseFlags(app.Flag)
	} else {
		app.CommandName = "help"
	}
}

func (app *App) NewFlag() *flag.FlagSet {
	newFlag := flag.NewFlagSet(app.CommandName, flag.ExitOnError)
	newFlag.StringVar(&app.ConfigPath, "config", "../../config/common.yaml", "配置文件路径")
	newFlag.StringVar(&app.EnvDir, "env", "../../", "env文件路径")
	return newFlag
}

func (app *App) ParseFlags(flag *flag.FlagSet) {
	flag.Parse(os.Args[2:])
}

func (app *App) InitEnv() {
	errLoadEnv := dotnev.Load(app.EnvDir, ".env")
	helper.CheckErrThenAbort("failed to load env", errLoadEnv)
}

func (app *App) InitConfig() {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	errParseConfig := config.LoadFiles(app.ConfigPath)
	helper.CheckErrThenAbort("failed to parse config", errParseConfig)
}

func (app *App) InitFramework(configure func(app *App)) {
	app.InitArgs(configure)
	app.InitEnv()
	app.InitConfig()
	timezone.Setup()
	mysql.Setup()
	phabricator.Setup()
}

func (app *App) Shutdown() {
	mysql.CloseClients()
}
