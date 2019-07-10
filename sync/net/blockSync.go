package net

import (
	"bethapi/sync/db"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

func StartSyncFromNet(ch <-chan int64,stop <-chan struct{},dbCh chan<- db.BLockValue,wg *sync.WaitGroup){
	defer wg.Done()
	list:=make(chan int64,500)
	syncCh := make(chan struct{},3)
	first:=true
	var latest int64 = 0
	for{
		select {
		case i:=<-ch:
			if first {
				latest = i
				first = false
				logs.Debug("Sync：区块",latest,"进入队列，等待同步")
				list<-i
			}else{
				if latest < i{
					for latest+1<i  {
						latest++
						logs.Debug("Sync：区块",latest,"进入队列，等待同步")
						list<-latest
					}
					latest = i
					logs.Debug("Sync：区块",latest,"进入队列，等待同步")
					list<-i
				}
			}
		case n:=<-list:
			go SyncBlockAndSend(n,list,syncCh,dbCh)
		case <-stop:
			return
		}
	}
}


func SyncBlockAndSend(n int64,list chan<- int64,syncLock chan struct{},dbCh chan<- db.BLockValue){
	syncLock<-struct{}{}
	logs.Debug("Sync：开始同步区块",n)
	result:=GetFullBlock(n,dbCh)
	if result.Head.Hash==""{
		logs.Warn("Sync：区块",n,"同步出错，准备重试")
		time.Sleep(2*time.Second)
		list<-n
	}
	<-syncLock
}