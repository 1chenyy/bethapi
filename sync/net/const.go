package net

const (
	key_1 = "YourApiKeyToken"
	key_2 = "YourApiKeyToken"
	key_3 = "YourApiKeyToken"
	key_4 = "YourApiKeyToken" //myregeth

	PriceKey = ""

	GetBlockNumberUrl = "https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey="+key_1

	GetBlockByNumber1 = "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag="
	GetBlockByNumber2 = "&boolean=true&apikey="+key_2

	GetBlockRewardByNumber1 = "https://api.etherscan.io/api?module=block&action=getblockreward&blockno="
	GetBlockRewardByNumber2 = "&apikey="+key_3

	GetETHPriceAndMarket = "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD&api_key="+PriceKey
	GetETHHistoryPrice = "https://min-api.cryptocompare.com/data/histoday?fsym=ETH&tsym=USD&limit=15&api_key="+PriceKey

	JSON = "application/json"

	MARSHAL_ERROR_RESPONSE = `{"status":0,"result":{},"error":"marshal error"}`
	INPUT_ERROR_RESPONSE   = `{"status":0,"result":{},"error":"input error"}`
	NO_DATA_ERROR_RESPONSE = `{"status":0,"result":{},"error":"no data"}`

	RPC_URL = "http://127.0.0.1:8545"

	ERROR_DB = "db error"
)

