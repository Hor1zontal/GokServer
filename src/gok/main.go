package main


import (
	"gok/module/game"
	"gok/module/passport"
	"gok/module/community"
	"gok/module/star"
	"gok/module/cluster"
	"gok/module/log"
	"gok/module/database"
	"gok/module/center"
	"gok/module/mail"
	"gok/module/trade"
	"gok/module/search"
	"gok/app"
	"gok/module/gate"
	"gok/module/statistics"
)


func main() {
	app.Run(
		cluster.Module,
		log.Module,
		database.Module,
		statistics.Module,
		gate.Module,
		search.Module,
		community.Module,
		trade.Module,
		star.Module,
		passport.Module,
		center.Module,
		game.Module,
		mail.Module,
	)

}
