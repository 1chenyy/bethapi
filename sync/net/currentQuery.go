package net

import (
	"bethapi/sync/db"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

func StartQuery(ch chan<- int64,stop <-chan struct{},wg *sync.WaitGroup,price *db.PriceDBUtil){
	defer wg.Done()
	go GetPriceAndSend(price)
	latestTicker :=time.NewTicker(5*time.Second)
	priceTicker := time.NewTicker(10*time.Minute)
	for{
		select {
		case <-latestTicker.C:
			go GetAndSendLatestMsg(ch)
		case <-priceTicker.C:
			go GetPriceAndSend(price)
		case <-stop:
			return
		}
	}
}

func GetAndSendLatestMsg(ch chan<- int64){
	i:=GetLatestBlock()
	if i != -1 {
		logs.Debug("Query：查询到最新区块：",i)
		ch<-i
	}
}

func GetPriceAndSend(price *db.PriceDBUtil){
	pam:=GetPriceAndMktcap()
	price.SetPriceAndMakcap(pam)
	history:=GetHistoryPrice()
	price.SetHistoryPrice(history)
}