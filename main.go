package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/larsfox/newochem-bot/core"
	"github.com/larsfox/newochem-bot/db"
	"github.com/larsfox/newochem-bot/tg"
)

var (
	app = cli.App("tgfinbot", "Launches a Test bot for Newochem")

	dbUser = app.StringOpt("db-user", "", "Database user")
	dbPass = app.StringOpt("db-pass", "", "Database pass")
	dbName = app.StringOpt("db-name", "", "Database name")
	dbHost = app.StringOpt("db-host", "", "Database host")

	tgToken = app.StringOpt("tg-token", "", "Telegram Bot Token")
	tgUsers = app.StringsOpt("tg-user", nil, "List of users who have access to bot")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
	app.Action = appAction
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("cli: ", err)
	}
}

func appAction() {
	tgClient := tg.NewClient(*tgToken)

	var host string
	if *dbHost != "" {
		host = fmt.Sprintf("tcp(%s)", *dbHost)
	}

	dbClient, err := db.NewClient(
		fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
			*dbUser, *dbPass, host, *dbName))

	if err != nil {
		log.Println(err)
		return
	}
	defer dbClient.CloseDB()

	appManager := core.NewManager(dbClient, tgClient, *tgUsers)
	appManager.Listen()
}
