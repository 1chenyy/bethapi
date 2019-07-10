package controllers

import (
	"bethapi/models"
	"bethapi/sync/db"
	"bethapi/sync/net"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
)

type StatisticController struct {
	beego.Controller
	TxDB    *db.TransactionDB
	PriceDB *db.PriceDBUtil
}

func NewStatisticController(txDB *db.TransactionDB,price *db.PriceDBUtil)*StatisticController{
	return &StatisticController{TxDB: txDB, PriceDB:price}
}

func (c *StatisticController)GetTransactionsStatistic(){
	logs.Error(c.PriceDB)
	i:=time.Duration(1)
	result := models.TransactionsStatistic{}
	for i=1;i<=15;i++{
		time:=time.Now().Add(-time.Hour*24*i).String()[:10]
		num:=c.TxDB.GetTransactionsStatisticByTime(time)
		if num == -1 {
			result.Number = append(result.Number, 0)
		}else{
			result.Number = append(result.Number, num)
		}
	}
	response:=models.TransactionsStatisticResponse{}
	response.Result = result
	response.Status = 1
	b,err:=json.Marshal(response)
	if err!=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
	}else{
		c.Ctx.Output.Body(b)
	}
}

func (c *StatisticController)GetPriceAndMakcap(){
	result:=c.PriceDB.PAM
	if result.Price=="" {
		c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
		return
	}
	response:=models.PriceAndMakcapResponse{}
	response.Result = result
	response.Status = 1
	b,err:=json.Marshal(response)
	if err!=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
	}else{
		c.Ctx.Output.Body(b)
	}
}

func (c *StatisticController)GetHistoryPrice(){
	result:=c.PriceDB.History
	if len(result.Price)==0 {
		c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
		return
	}
	response:=models.HistoryPriceResponse{}
	response.Result = result
	response.Status = 1
	b,err:=json.Marshal(response)
	if err!=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
	}else {
		c.Ctx.Output.Body(b)
	}
}