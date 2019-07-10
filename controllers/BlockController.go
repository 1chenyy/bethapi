package controllers

import (
	"bethapi/models"
	"bethapi/sync/db"
	"bethapi/sync/net"
	"bethapi/util"
	"encoding/json"
	"github.com/astaxie/beego"
	"strings"
)

type BlockController struct {
	beego.Controller
	DB *db.BlockDB
}

func NewBlockController(db *db.BlockDB)*BlockController{
	return &BlockController{DB:db}
}

func(c *BlockController)GetBlockByNumber(){
	num:=util.StringToInt64(c.Ctx.Input.Param(":num"))
	block:=models.BlockSummary{}
	if num == -1{
		c.Ctx.Output.Body([]byte(net.INPUT_ERROR_RESPONSE))
		return
	}
	if num>c.DB.GetLatestBlockNumber() {
		c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
		return
	}
	block = c.DB.GetBlockSummaryByNumberFromDB(num)
	if block.Miner == "" {
		fullblock:=net.GetFullBlock(num,c.DB.DBCh)
		if fullblock.Head.Miner == "" {
			c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
			return
		}
		block = models.GenerateBlockSummaryFromBlock(fullblock)
	}
	response:=models.BlockResponse{}
	response.Status = 1
	response.Result = block
	b,err:=json.Marshal(response)
	if err !=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
		return
	}
	c.Ctx.Output.Body(b)
}

func (c *BlockController)GetLatest15Blocks(){
	bundle:=c.DB.GetLatest15Blocks()
	response:=models.BlockBundleResponse{}
	response.Status = 1
	response.Result = bundle
	b,err:=json.Marshal(response)
	if err !=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
		return
	}
	c.Ctx.Output.Body(b)
}

func (c *BlockController)GetLatestBlocks(){
	num:=util.StringToInt64(c.Ctx.Input.Param(":num"))
	bundle:=models.BlockBundle{}
	if num == -1{
		c.Ctx.Output.Body([]byte(net.INPUT_ERROR_RESPONSE))
		return
	}
	if num>c.DB.GetLatestBlockNumber() {
		c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
		return
	}
	bundle=c.DB.GetLatestBlocks(num)
	response:=models.BlockBundleResponse{}
	response.Status = 1
	response.Result = bundle
	b,err:=json.Marshal(response)
	if err !=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
		return
	}
	c.Ctx.Output.Body(b)
}

func (c *BlockController)GetTransactionsByNumber(){
	num:=util.StringToInt64(c.Ctx.Input.Param(":num"))
	txs:=models.Transactions{}
	if num == -1{
		c.Ctx.Output.Body([]byte(net.INPUT_ERROR_RESPONSE))
		return
	}
	if num>c.DB.GetLatestBlockNumber() {
		c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
		return
	}
	response:=models.TransactionsResponse{}

	txs,len:=c.DB.GetTransactionsByNumberFromDB(num)
	if len == -1 {
		localBlock:=net.GetLocalBlockByNumber(num)
		if localBlock.Miner == "" {
			fullblock:=net.GetFullBlock(num,c.DB.DBCh)
			if fullblock.Head.Miner == "" {
				c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
				return
			}else{
				txs.Txs = db.GenerateTransactionsSummary(fullblock.Head.Transactions)
			}
		}else{
			txs.Txs = db.GenerateTransactionsSummary(localBlock.Transactions)
		}
	}
	response.Status = 1
	response.Result = txs
	b,err:=json.Marshal(response)
	if err !=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
		return
	}
	c.Ctx.Output.Body(b)
}

func (c *BlockController)GetTransactionDetails(){
	tx:=models.Transaction{}
	hashnums:=strings.Split(c.Ctx.Input.Param(":hashnum"),"-")
	if len(hashnums)!=2 {
		c.Ctx.Output.Body([]byte(net.INPUT_ERROR_RESPONSE))
		return
	}
	hash:=hashnums[0]
	num:=util.StringToInt64(hashnums[1])
	if num == -1 {
		c.Ctx.Output.Body([]byte(net.INPUT_ERROR_RESPONSE))
		return
	}
	dbblock := c.DB.GetBlockByNumberFromDB(num)
	if dbblock.Head.Miner == "" {
		localBlock:=net.GetLocalBlockByNumber(num)
		if localBlock.Miner == "" {
			fullblock:=net.GetFullBlock(num,c.DB.DBCh)
			if fullblock.Head.Miner == "" {
				c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
				return
			}else{
				tx = db.FindTransaction(hash,fullblock.Head.Transactions)
			}
		}else{
			tx = db.FindTransaction(hash,localBlock.Transactions)
		}
	}else{
		tx=db.FindTransaction(hash, dbblock.Head.Transactions)
	}
	response:=models.TransactionResponse{}
	response.Status = 1
	response.Result = tx
	b,err:=json.Marshal(response)
	if err !=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
		return
	}
	c.Ctx.Output.Body(b)
}

func(c *BlockController)GetOnePageBlocks(){
	num:=util.StringToInt64(c.Ctx.Input.Param(":num"))
	bundle:=models.BlockBundle{}
	if num == -1{
		c.Ctx.Output.Body([]byte(net.INPUT_ERROR_RESPONSE))
		return
	}
	if num>c.DB.GetLatestBlockNumber() {
		c.Ctx.Output.Body([]byte(net.NO_DATA_ERROR_RESPONSE))
		return
	}
	bundle = c.DB.GetOnePageBlocks(num)
	response:=models.BlockBundleResponse{}
	response.Status = 1
	response.Result = bundle
	b,err:=json.Marshal(response)
	if err !=nil {
		c.Ctx.Output.Body([]byte(net.MARSHAL_ERROR_RESPONSE))
		return
	}
	c.Ctx.Output.Body(b)
}
