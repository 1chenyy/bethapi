package sync

import (
	"bethapi/sync/db"
	"bethapi/sync/net"
	"github.com/astaxie/beego/logs"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var asyncNum = 5

func StartSync(blockDB *db.BlockDB,transactionDB *db.TransactionDB,priceDB *db.PriceDBUtil){
	queryCh := make(chan int64,1)
	var wg sync.WaitGroup
	stopSign := make(chan os.Signal, 1)
	signal.Notify(stopSign, syscall.SIGINT)
	stop := make(chan struct{},asyncNum)

	wg.Add(asyncNum)

	go net.StartQuery(queryCh,stop,&wg,priceDB)
	go net.StartSyncFromNet(queryCh,stop, blockDB.DBCh,&wg)
	go blockDB.StartWriteToDB(stop,&wg,transactionDB.DBCh,transactionDB.StartCh)
	go ListenerStop(stopSign,stop)
	go transactionDB.StartWriteToDB(stop,&wg)
	go transactionDB.StartTrace(stop,&wg)

	wg.Wait()

	logs.Warn("sync exit!")
}

func StartTransactionsStatistics(){

}

func ListenerStop(stopSign <-chan os.Signal,stop chan<- struct{}){
	select {
	case <-stopSign:
		for i := 0;i<asyncNum;i++{
			stop<- struct{}{}
		}
		return
	}
}
