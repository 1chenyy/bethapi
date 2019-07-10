package models

type TransactionsStatisticResponse struct {
	Status int `json:"status"`
	Result TransactionsStatistic `json:"result"`
	Error string `json:"error"`
}

type TransactionsStatistic struct {
	Number []int64 `json:"number"`
}

type ResponsePriceAndMakcap struct {
	Raw ResponsePriceAndMakcapRaw `json:"RAW"`
}

type ResponsePriceAndMakcapRaw struct {
	Eth ResponsePriceAndMakcapRawEth `json:"ETH"`
}

type ResponsePriceAndMakcapRawEth struct {
	USD ResponsePriceAndMakcapRawEthUSD `json:"USD"`
}

type ResponsePriceAndMakcapRawEthUSD struct {
	Price float64 `json:PRICE`
	Mktcap float64 `json:"MKTCAP"`
}

type PriceAndMakcapResponse struct {
	Status int `json:"status"`
	Result PriceAndMakcap `json:"result"`
	Error string `json:"error"`
}

type PriceAndMakcap struct {
	Price string `json:"price"`
	Mktcap string `json:"mktcap"`
}

type ResponseHistoryPrice struct {
	Data []ResponseHistoryPriceDaily `json:"Data"`
}

type ResponseHistoryPriceDaily struct {
	Close float64 `json:"close"`
}

type HistoryPrice struct {
	Price []string `json:"price"`
}

type HistoryPriceResponse struct {
	Status int `json:"status"`
	Result HistoryPrice `json:"result"`
	Error string `json:"error"`
}