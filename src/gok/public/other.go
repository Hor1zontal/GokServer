package main


import (
	"github.com/name5566/leaf"
	"gok/module/passport"
	"gok/module/community"
	"gok/module/cluster"
	"time"
	"math/rand"
	"gok/module/log"
	"gok/module/database"
	"gok/module/center"
	"gok/module/search"
	"gok/module/mail"
	"gok/module/statistics"
	"gok/module/trade"
)


func main() {
	rand.Seed(time.Now().UnixNano())
	leaf.Run(
		cluster.Module,
		log.Module,
		database.Module,
		statistics.Module,

		community.Module,
		trade.Module,
		//star.Module,
		passport.Module,
		center.Module,
		search.Module,
		mail.Module,
	)


}
