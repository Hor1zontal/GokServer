package main


import (
	"gok/module/cluster"
	"gok/module/log"
	"gok/module/database"
	"gok/module/game"
	"gok/module/gate"
	"gok/module/statistics"
	"gok/app"
)


func main() {
	app.Run(
		cluster.Module,
		log.Module,
		database.Module,
		statistics.Module,

		game.Module,
		gate.Module,
	)

}
