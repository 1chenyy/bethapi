package net

import (
	"bethapi/models"
	"bethapi/util"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var RPCClient = http.Client{Timeout:30*time.Second}

func GetLocalLatestBlockNumber()int64{
	req:=`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`
	resp,err:=RPCClient.Post(RPC_URL,JSON,strings.NewReader(req))
	if err!=nil{
		logs.Warn("GetLocalLatestBlockNumbe PostError:",err)
		return -1
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetLocalLatestBlockNumbe ReadError:",err)
		return -1
	}
	result := &models.BlockNumber{}
	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetLocalLatestBlockNumbe UnmarshalError:",err)
		return -1
	}
	return util.HexToDec(result.Result)
}

func GetLocalBlockByNumber(i int64)models.BlockHead{
	result := &models.BlockByNumberResponse{}
	numHex:="0x"+util.DecToHex(i)
	req:=`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["`+numHex+`", true],"id":1}`
	resp,err:=RPCClient.Post(RPC_URL,JSON,strings.NewReader(req))
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

func GetLocalTransactionByHash(hash string)models.Transaction{
	result:=&models.ResponseTransaction{}
	req:=`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["`+hash+`"],"id":1}`
	resp,err:=RPCClient.Post(RPC_URL,JSON,strings.NewReader(req))
	if err != nil{
		logs.Warn("GetLocalTransactionByHash GetError:",err)
		return result.Result
	}
	defer resp.Body.Close()
	response,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		logs.Warn("GetLocalTransactionByHash ReadError:",err)
		return result.Result
	}
	err=json.Unmarshal(response,result)
	if err!=nil {
		logs.Warn("GetLocalTransactionByHash UnmarshalError:",err)
		result.Result = models.Transaction{}
		return result.Result
	}
	return result.Result
}
