package controllers

import (
	"bethapi/models"
	"bethapi/sync/db"
	"bethapi/sync/net"
	"encoding/json"
	"github.com/astaxie/beego"
)

type CurrentBlockController struct {
	beego.Controller
	DB *db.BlockDB
}

func NewCurrentBlockController(db *db.BlockDB)*CurrentBlockController{
	return &CurrentBlockController{DB:db}
}

func(c *CurrentBlockController)Get(){
	resp:=models.CurrentBlockResponse{}
	num:=c.DB.GetLatestBlockNumber()
	resp.Result = models.CurrentBlock{
		Number:num,
		TotalDifficult:c.DB.GetBlockSummaryByNumberFromDB(num).TotalDifficulty}
	resp.Status = 1
	if resp.Result.Number == -1 {
		resp.Error = net.ERROR_DB
	}
	b,err:=json.Marshal(resp)
	if err != nil{
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
	}else{
		c.Ctx.Output.Body(b)
	}
}

