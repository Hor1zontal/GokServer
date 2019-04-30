package main

import (
	"gok/module/cluster"
	"gok/module/log"
	"gok/module/database"
	"gok/module/statistics"
	"gok/module/mail"
	"gok/app"
)

func main() {
	app.Run(
		cluster.Module,
		log.Module,
		database.Module,
		statistics.Module,

		mail.Module,
	)

}
