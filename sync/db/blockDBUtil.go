package db

import (
	"bethapi/models"
	"bethapi/util"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"sync"
)

type BLockValue struct {
	Num int64
	Data []byte
	Hash string
	TxNum int64
	Time string
}

var LATEST_BLOCK_KEY = []byte("latest_block")
var DIVIDE = []byte("divide")

type BlockDB struct {
	Lvdb *leveldb.DB
	DBCh chan BLockValue
	Count int64
}

func NewBlockDBUtil()*BlockDB {
	return &BlockDB{Lvdb: CreateBlockDB(),DBCh:make(chan BLockValue,500),Count:0}
}

func CreateBlockDB()*leveldb.DB{
	db, err :=leveldb.OpenFile("blockdb",nil)
	if err!=nil {
		log.Fatalln(err)
	}
	return db
}

func (u *BlockDB)StartWriteToDB(stop <-chan struct{},wg *sync.WaitGroup,txCh chan<- TransactionValue, startTrace chan<- int64)  {
	defer wg.Done()
	first:=false
	for {
		select {
		case v:=<-u.DBCh:
			logs.Debug("txDB：开始写入区块",v.Num,"长度",len(v.Data))
			err:=u.Lvdb.Put(BlockNumberKey(v.Num),v.Data,nil)
			if err != nil {
				logs.Warn("txDB：写入区块失败",err)
			}else {
				if u.GetLatestBlockNumber() < v.Num {
					err=u.Lvdb.Put(LATEST_BLOCK_KEY,util.Int64ToBytes(v.Num),nil)
					if err != nil {
						logs.Warn("txDB：写入高度失败：",v.Num)
					}else{
						logs.Debug("txDB：写入高度成功：",v.Num)
					}
				}
				if !first {
					u.Lvdb.Put(DIVIDE,util.Int64ToBytes(v.Num),nil)
					if u.GetDivideNumber() != -1 {
						startTrace<-v.Num
						first = true
					}
				}
				value:=TransactionValue{}
				value.Time = []byte(v.Time)
				value.Num = v.TxNum

				go func(value TransactionValue,txCh chan<- TransactionValue) {
					txCh<-value
				}(value,txCh)

				logs.Debug("txDB：区块",v.Num,"写入区块成功")
				u.Count++
			}
		case <-stop:
			return
		}
	}
}

func (u *BlockDB)GetLatestBlocks(num int64)models.BlockBundle{
	latest:=u.GetLatestBlockNumber()
	if num < latest {
		return u.GetLatestNBlocks(latest-num)
	}
	return models.BlockBundle{}
}

func (u *BlockDB)GetLatestNBlocks(n int64)models.BlockBundle{
	num:=u.GetLatestBlockNumber()
	return u.GetBlockBundle(num,n)
}

func (u *BlockDB)GetLatest15Blocks()models.BlockBundle{
	num:=u.GetLatestBlockNumber()
	return u.GetBlockBundle(num,15)
}

func (u *BlockDB)GetOnePageBlocks(num int64)models.BlockBundle{
	return u.GetBlockBundle(num,10)
}

func (u *BlockDB)GetBlockBundle(num int64, length int64)models.BlockBundle{
	bundle := models.BlockBundle{}
	if !u.IsValidNum(num) {
		return bundle
	}
	if length > num {
		return bundle
	}
	if length > u.Count {
		length = u.Count
	}
	for i := int64(0); i< length; i++ {
		summary:=u.GetBlockSummaryByNumberFromDB(num-i)
		if summary.Number != -1{
			bundle.Blocks = append(bundle.Blocks, summary)
		}else{
			bundle.Blocks = append(bundle.Blocks,models.BlockSummary{})
		}
	}
	return bundle
}

func (u *BlockDB)GetTransactionsByNumberFromDB(i int64)(models.Transactions,int){
	txs:=models.Transactions{}
	if !u.IsValidNum(i) {
		return txs,-1
	}
	v,err:=u.Lvdb.Get(BlockNumberKey(i),nil)
	if err != nil {
		return txs,-1
	}
	block := &models.Block{}
	err=json.Unmarshal(v,block)
	if err != nil {
		return txs,-1
	}
	txs.Txs = GenerateTransactionsSummary(block.Head.Transactions)
	return txs,len(txs.Txs)
}

func GenerateTransactionsSummary(txs []models.Transaction)[]models.TransactionSummary{
	var stxs []models.TransactionSummary
	for _,v:=range txs{
		stxs = append(stxs, models.TransactionSummary{
			Hash:v.Hash,
			From:v.From,
			To:v.To,
			Fee:util.CalculateFee(v.Gas,v.GasPrice),
			Value:util.HexWeiToEth(v.Value),
		})
	}
	return stxs
}

func FindTransaction(hash string,txs []models.Transaction)models.Transaction{
	tx:=models.Transaction{}
	for _,v:=range txs{
		if v.Hash == hash {
			tx = v
		}
	}
	return tx
}

func (u *BlockDB)GetBlockSummaryByNumberFromDB(i int64)models.BlockSummary{
	summary := models.BlockSummary{}
	if !u.IsValidNum(i) {
		return summary
	}
	summary.Number = -1
	v,err:=u.Lvdb.Get(BlockNumberKey(i),nil)
	if err != nil {
		return summary
	}
	block := &models.Block{}
	err=json.Unmarshal(v,block)
	if err != nil {
		return summary
	}
	return models.GenerateBlockSummaryFromBlock(*block)
}

func(u *BlockDB)GetBlockByNumberFromDB(i int64)models.Block{
	block:=&models.Block{}
	if !u.IsValidNum(i) {
		return models.Block{}
	}
	v,err:=u.Lvdb.Get(BlockNumberKey(i),nil)
	if err != nil {
		return models.Block{}
	}
	err=json.Unmarshal(v,block)
	if err != nil {
		return models.Block{}
	}
	return *block
}

func (u *BlockDB)GetLatestBlockNumber()int64{
	v,err:=u.Lvdb.Get(LATEST_BLOCK_KEY,nil)
	if err != nil {
		return -1
	}
	return util.BytesToInt64(v)
}

func (u *BlockDB)GetDivideNumber()int64{
	v,err:=u.Lvdb.Get(DIVIDE,nil)
	if err != nil {
		return -1
	}
	return util.BytesToInt64(v)
}

func (u *BlockDB)IsValidNum(num int64)bool{
	return u.GetLatestBlockNumber() >= num
}

func BlockNumberKey(num int64)[]byte{
	return []byte("n"+util.DecToString(num))
}

func BlockHashKey(hash string)[]byte{
	return []byte(hash)
}