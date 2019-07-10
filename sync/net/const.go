package net

const (
	key_1 = "QMIGUVHKYB1XSTHXQ7UXCSJ8M9K7GBFIJE"
	key_2 = "MM7JQYHR4UTYQVPWKTF1QPXZYWD8B3947Y"
	key_3 = "DRVYICXS84A2328KAGGCQ3R41U4EJT8Z95"
	key_4 = "ICGS6JVIX3JJKURQZB2IWFR23SZSZCDVF9" //myregeth

	PriceKey = "2a30e5cc2e44d10602bc7993fb0b0567755a37f8858fd249133a9c6d28746ccd"

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

	RPC_URL = "http://10.166.33.85:8546"

	ERROR_DB = "db error"
)

