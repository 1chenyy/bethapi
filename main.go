package main

import (
	"bethapi/controllers"
	"bethapi/sync"
	"bethapi/sync/db"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	blockDB := db.NewBlockDBUtil()
	defer blockDB.Lvdb.Close()

	transactionDB := db.NewTransactionDBUtil()
	defer transactionDB.Lvdb.Close()

	priceDB := db.NewPriceDBUtil()

	go sync.StartSync(blockDB,transactionDB,priceDB)

	initRoute(blockDB,transactionDB,priceDB)

	logs.SetLogger("console")
	logs.Async(1000)
	beego.Run()
	logs.Warn("api exit")
}

func initRoute(blockDB *db.BlockDB,statisticDB *db.TransactionDB,priceDB *db.PriceDBUtil){
	currentBlockController:=controllers.NewCurrentBlockController(blockDB)
	blockController:=controllers.NewBlockController(blockDB)
	statisticController:=controllers.NewStatisticController(statisticDB,priceDB)

	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/current",
			currentBlockController),

		beego.NSRouter("/getblockbynum/:num",
			blockController,"get:GetBlockByNumber"),

		beego.NSRouter("/getlatest15blocks",
			blockController,"get:GetLatest15Blocks"),

		beego.NSRouter("/getlatestblocks/:num",
			blockController,"get:GetLatestBlocks"),

		beego.NSRouter("/gettransactionsbynum/:num",
			blockController,"get:GetTransactionsByNumber"),

		beego.NSRouter("/gettransactiondetails/:hashnum",
			blockController,"get:GetTransactionDetails"),

		beego.NSRouter("/getonepageblocks/:num",
			blockController,"get:GetOnePageBlocks"),

		beego.NSRouter("/gettransactionsstatistic",
			statisticController,"get:GetTransactionsStatistic"),

		beego.NSRouter("/getpriceandmakcap",
			statisticController,"get:GetPriceAndMakcap"),

		beego.NSRouter("/gethistoryprice",
			statisticController,"get:GetHistoryPrice"),
	)
	beego.AddNamespace(ns)
}


