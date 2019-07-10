package models

import (
	"bethapi/util"
)

type BlockByNumberResponse struct {
	Version string `json:"jsonrpc"`
	Id int	`json:"id"`
	Result BlockHead `json:"result"`
}

type BlockHead struct {
	Difficulty string `json:"difficulty"`
	ExtraData string `json:"extraData"`
	GasLimit string `json:"gasLimit"`
	GasUsed string `json:"gasUsed"`
	Hash string                `json:"hash"`
	LogsBloom string           `json:"logsBloom"`
	Miner string               `json:"miner"`
	MixHash string             `json:"mixHash"`
	Nonce string               `json:"nonce"`
	Number string              `json:"number"`
	ParentHash string          `json:"parentHash"`
	ReceiptsRoot string        `json:"receiptsRoot"`
	Sha3Uncles string          `json:"sha3Uncles"`
	Size string                `json:"size"`
	StateRoot string           `json:"stateRoot"`
	Timestamp string           `json:"timestamp"`
	TotalDifficulty string     `json:"totalDifficulty"`
	TransactionsRoot string    `json:"transactionsRoot"`
	Transactions []Transaction `json:"transactions"`
	Uncles []string            `json:"uncles"`
}

type Transaction struct {
	BlockHash string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
	From string `json:"from"`
	Gas string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Hash string `json:"hash"`
	Input string `json:"input"`
	Nonce string `json:"nonce"`
	R string `json:"r"`
	S string `json:"s"`
	V string `json:"v"`
	To string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value string `json:"value"`
}

type BlockRewardResponse struct {
	Message string `json:"message"`
	Result BlockReward `json:"result"`
	Status string  `json:"status"`
}

type BlockReward struct {
	BlockMiner string `json:"blockMiner"`
	BlockNumber string `json:"blockNumber"`
	BlockReward string `json:"blockReward"`
	TimeStamp string `json:"timeStamp"`
	UncleInclusionReward string `json:"uncleInclusionReward"`
	Uncles []UncleReward `json:"uncles"`
}

type UncleReward struct {
	Blockreward string `json:"blockreward"`
	Miner string `json:"miner"`
	UnclePosition string `json:"unclePosition"`
}

type Block struct {
	Head BlockHead
	Reward BlockReward
}

type BlockSummary struct {
	Number int64		`json:"number"`
	Hash string 	`json:"hash"`
	Time int64		`json:"time"`
	Txs int			`json:"txs"`
	Miner string	`json:"miner"`
	Reward string	`json:"reward"`
	UncleReward string	`json:"unclereward"`
	Difficult string	`json:"difficult"`
	TotalDifficulty string	`json:"totaldifficulty"`
	Size int64	`json:"size"`
	GasUsed int64 `json:"gasused"`
	GasLimit int64	`json:"gaslimit"`
	Extra string `json:"extra"`
}

type BlockResponse struct {
	Status int `json:"status"`
	Result BlockSummary `json:"result"`
	Error string `json:"error"`
}

type BlockBundle struct {
	Blocks []BlockSummary `json:"blocks"`
}

type BlockBundleResponse struct {
	Status int `json:"status"`
	Result BlockBundle `json:"result"`
	Error string `json:"error"`
}

type Transactions struct {
	Txs []TransactionSummary `json:"txs"`
}

type TransactionSummary struct {
	Hash string `json:"hash"`
	From string `json:"from"`
	To string `json:"to"`
	Fee string `json:"fee"`
	Value string `json:"value"`
}

type TransactionsResponse struct {
	Status int `json:"status"`
	Result Transactions `json:"result"`
	Error string `json:"error"`
}

type TransactionResponse struct {
	Status int `json:"status"`
	Result Transaction `json:"result"`
	Error string `json:"error"`
}

type ResponseTransaction struct {
	Version string `json:"jsonrpc"`
	Id int	`json:"id"`
	Result Transaction `json:"result"`
}

func GenerateBlockSummaryFromBlock(block Block)(summary BlockSummary){
	summary.Number = util.HexToDec(block.Head.Number)
	summary.Hash = block.Head.Hash
	summary.Time = util.HexToDec(block.Head.Timestamp)
	summary.Txs = len(block.Head.Transactions)
	summary.Miner = block.Head.Miner
	summary.Reward = util.IntWeiToEth(block.Reward.BlockReward)
	summary.UncleReward = util.IntWeiToEth(block.Reward.UncleInclusionReward)
	summary.Difficult = util.HexToIntString(block.Head.Difficulty)
	summary.TotalDifficulty = util.HexToIntString(block.Head.TotalDifficulty)
	summary.Size = util.HexToDec(block.Head.Size)
	summary.GasUsed = util.HexToDec(block.Head.GasUsed)
	summary.GasLimit = util.HexToDec(block.Head.GasLimit)
	summary.Extra = util.HexToString(block.Head.ExtraData)
	return
}
