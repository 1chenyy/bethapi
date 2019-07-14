package db

import (
	"bethapi/models"
	"bethapi/util"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type TransactionValue struct {
	Time []byte
	Num  int64
}

type TransactionDB struct {
	Lvdb *leveldb.DB
	DBCh chan TransactionValue
	StartCh chan int64
	Count int64
	StopTrace bool
}

func NewTransactionDBUtil()*TransactionDB {
	db:=&TransactionDB{Lvdb: CreateTransactionDB(),
		DBCh:make(chan TransactionValue,1000),
		Count:0,
		StartCh:make(chan int64),
		StopTrace:false,
	}
	isRestart,err:=beego.AppConfig.Bool("restart")
	if err == nil && isRestart {
		db.Lvdb.Put([]byte(time.Now().String()[:10]),util.Int64ToBytes(0),nil)
	}
	return db
}

func CreateTransactionDB()*leveldb.DB{
	db, err :=leveldb.OpenFile("transactiondb",nil)
	if err!=nil {
		log.Fatalln(err)
	}
	return db
}

func (u *TransactionDB)GetTransactionsStatisticByTime(time string)int64{
	result,err:=u.Lvdb.Get([]byte(time),nil)
	if err!=nil {
		return -1
	}
	return util.BytesToInt64(result)
}

func (u *TransactionDB)StartWriteToDB(stop <-chan struct{},wg *sync.WaitGroup)  {
	defer wg.Done()
	for {
		select {
		case v:=<-u.DBCh:
			i:=int64(0)
			result,err:=u.Lvdb.Get(v.Time,nil)
			if err == nil {
				i = util.BytesToInt64(result)
			}
			i+=v.Num
			logs.Debug("更新",string(v.Time),"交易信息，总数",i)
			u.Lvdb.Put(v.Time,util.Int64ToBytes(i),nil)
		case <-stop:
			return
		}
	}
}

func (u *TransactionDB)StartTrace(stop <-chan struct{},wg *sync.WaitGroup){
	defer wg.Done()
	for{
		select {
		case divide:=<-u.StartCh:
			go u.TraceBlocks(divide)
		case <-stop:
			u.StopTrace = true
			return
		}
	}
}

func (u TransactionDB)TraceBlocks(divide int64){
	isRestart,err:=beego.AppConfig.Bool("restart")
	if err!=nil {
		isRestart = false;
	}
	start:=int64(0)
	if isRestart {
		start = 7200
	}else{
		start = 100000
	}
	today:=time.Now().String()[:10]
	for i:=divide-start;i<divide && !u.StopTrace;i++ {
		localblock:=u.GetLocalBlockByNumber(i)
		blockTime := time.Unix(util.HexToDec(localblock.Timestamp),0).String()[:10]
		if isRestart && blockTime!=today {
			logs.Debug("跳过区块：",localblock.Number)
			continue
		}
		value:=TransactionValue{}
		value.Time = []byte(blockTime)
		value.Num = int64(len(localblock.Transactions))
		u.DBCh<-value
	}
}

func (u TransactionDB)GetLocalBlockByNumber(i int64)models.BlockHead{
	result := &models.BlockByNumberResponse{}
	numHex:="0x"+util.DecToHex(i)
	req:=`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["`+numHex+`", true],"id":1}`
	resp,err:=http.Post("http://10.166.33.85:8546","application/json",strings.NewReader(req))
	if err != nil{
		logs.Warn("GetLocalBlockByNumber GetError:",err)
		return result.Result
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetLocalBlockByNumber ReadError:",err)
		return result.Result
	}

	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetLocalBlockByNumber UnmarshalError:",err)
		result.Result = models.BlockHead{}
		return result.Result
	}
	return result.Result
}


