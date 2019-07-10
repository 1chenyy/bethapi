package net

import (
	"bethapi/models"
	"bethapi/sync/db"
	"bethapi/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var GetLatestClient = &http.Client{Timeout:time.Second*60}

func GetLatestBlock()int64{
	resp,err:=GetLatestClient.Get(GetBlockNumberUrl)
	if err != nil{
		logs.Warn("GetBlockNumber GetError:",err)
		return -1
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetBlockNumber ReadError:",err)
		return -1
	}
	result := &models.BlockNumber{}
	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetBlockNumber UnmarshalError:",err)
		return -1
	}
	return util.HexToDec(result.Result)
}

var GetBlockClient = &http.Client{Timeout:time.Second*90}

func GetBlockByNumber(i int64)models.BlockHead{
	head := GetLocalBlockByNumber(i)
	if head.Miner!=""{
		return head
	}
	result := &models.BlockByNumberResponse{}
	numHex:="0x"+util.DecToHex(i)
	resp,err:=GetBlockClient.Get(GetBlockByNumber1+numHex+GetBlockByNumber2)
	if err != nil{
		logs.Warn("GetBlockByNumber GetError:",err)
		return result.Result
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetBlockByNumber ReadError:",err)
		return result.Result
	}

	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetBlockByNumber UnmarshalError:",err)
		result.Result = models.BlockHead{}
		return result.Result
	}
	return result.Result
}

func GetBlockRewardByNumber(i int64)models.BlockReward{
	result := &models.BlockRewardResponse{}
	resp,err:=GetBlockClient.Get(GetBlockRewardByNumber1+util.DecToString(i)+GetBlockRewardByNumber2)
	if err != nil{
		logs.Warn("GetBlockRewardByNumber GetError:",err)
		return result.Result
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetBlockRewardByNumber ReadError:",err)
		return result.Result
	}

	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetBlockRewardByNumber UnmarshalError:",err)
		result.Result = models.BlockReward{}
		return result.Result
	}
	return result.Result
}

func GetFullBlock(num int64,dbCh chan<- db.BLockValue)models.Block{
	var wg sync.WaitGroup
	var head models.BlockHead
	var reward models.BlockReward
	wg.Add(2)
	go func(n int64) {
		head = GetBlockByNumber(num)
		wg.Done()
	}(num)
	go func(n int64) {
		reward = GetBlockRewardByNumber(num)
		wg.Done()
	}(num)
	wg.Wait()
	block:=models.Block{}
	if head.Hash != "" && reward.BlockReward != "" {
		block.Head = head
		block.Reward = reward
		b,err:=json.Marshal(block)
		if err!=nil {
			logs.Warn("MarshalErr",err)
			block.Head.Hash = ""
			return block
		}else{
			logs.Debug("Sync：区块",num,"同步成功，准备写入数据库")
			dbCh<-db.BLockValue{
				Num: num,
				Data:b,
				Hash:head.Hash,
				TxNum:int64(len(block.Head.Transactions)),
				Time:time.Unix(util.HexToDec(block.Head.Timestamp),0).String()[:10],
			}
			return block
		}
	}else {
		return block
	}
}

func GetBlockSummaryFromNet(num int64,dbCh chan<- db.BLockValue)(summary models.BlockSummary){
	block:=GetFullBlock(num,dbCh)
	if block.Head.Hash!="" {
		summary = models.GenerateBlockSummaryFromBlock(block)
	}
	return
}

var PriceClient = &http.Client{Timeout:time.Second*90}

func GetPriceAndMktcap()models.PriceAndMakcap{
	logs.Debug("开始查询价格和市值")
	pam:=models.PriceAndMakcap{}

	resp,err:=PriceClient.Get(GetETHPriceAndMarket)
	if err != nil{
		logs.Warn("GetPriceAndMktcap GetError:",err)
		return pam
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetPriceAndMktcap ReadError:",err)
		return pam
	}
	result:=&models.ResponsePriceAndMakcap{}
	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetPriceAndMktcap UnmarshalError:",err)
		return pam
	}
	if result.Raw.Eth.USD.Mktcap == 0 || result.Raw.Eth.USD.Price == 0 {
		return pam
	}
	pam.Mktcap = fmt.Sprintf("%.2f", result.Raw.Eth.USD.Mktcap)
	pam.Price = fmt.Sprintf("%.2f", result.Raw.Eth.USD.Price)
	logs.Debug("价格与市值查询成功：",pam.Price,"---",pam.Mktcap)
	return pam
}

func GetHistoryPrice()models.HistoryPrice{
	logs.Debug("查询历史价格")
	history:=models.HistoryPrice{}
	resp,err:=PriceClient.Get(GetETHHistoryPrice)
	if err != nil {
		logs.Warn("GetHistoryPrice GetError:",err)
		return history
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetPriceAndMktcap ReadError:",err)
		return history
	}
	result:=&models.ResponseHistoryPrice{}
	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetHistoryPrice UnmarshalError:",err)
		return history
	}
	for _,v:=range result.Data{
		history.Price = append(history.Price, fmt.Sprintf("%.2f", v.Close))
	}
	logs.Debug("历史价格查询成功",history.Price)
	return history
}
