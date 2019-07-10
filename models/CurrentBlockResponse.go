package models

type CurrentBlockResponse struct {
	Status int `json:"status"`
	Result CurrentBlock `json:"result"`
	Error string `json:"error"`
}

type CurrentBlock struct {
	Number int64 `json:"number"`
	TotalDifficult string `json:"total_difficult"`
}

type BlockNumber struct {
	Version string `json:"jsonrpc"`
	Id int	`json:"id"`
	Result string `json:"result"`
}