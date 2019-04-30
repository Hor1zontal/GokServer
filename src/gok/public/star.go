package main


import (
	"gok/module/star"
	"gok/module/cluster"
	"gok/module/log"
	"gok/module/database"
	"gok/module/statistics"
	"gok/app"
)


func main() {
	app.Run(
		cluster.Module,
		log.Module,
		database.Module,
		statistics.Module,

		star.Module,
	)

}
