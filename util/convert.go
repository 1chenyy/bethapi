package util

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func HexToDec(hex string)int64{
	b:=big.Int{}
	_,result:=b.SetString(hex,0)
	if !result {
		return -1
	}
	return b.Int64()
}

func HexToIntString(hex string)string{
	b:=big.Int{}
	_,result:=b.SetString(hex,0)
	if !result {
		return ""
	}
	return b.String()
}

func HexToString(hexs string)string{
	b,err:=hex.DecodeString(RecoveryHexString(hexs))
	result := ""
	if err != nil {
		result = ""
	}else {
		result = string(b)
	}
	return result
}

func DecToHex(i int64)string{
	return strconv.FormatInt(i,16)
}

func DecToString(i int64)string{
	return strconv.FormatInt(i,10)
}

func StringToInt64(s string)int64{
	b:=big.Int{}
	_,result:=b.SetString(s,0)
	if !result {
		return -1
	}
	return b.Int64()
}

func Int64ToBytes(i int64)[]byte{
	var big big.Int
	big.SetInt64(i)
	return big.Bytes()
}

func BytesToInt64(b []byte)int64{
	var big big.Int
	big.SetBytes(b)
	return big.Int64()
}

func RecoveryHexString(hex string)string{
	if strings.Contains(hex, "0x") {
		return strings.Replace(hex,"0x","",1)
	}
	return hex
}

var base = float64(1e18)

func CalculateFee(a,b string)string{
	x:=HexToDec(a)
	y:=HexToDec(b)
	if x==-1 || y == -1 {
		return ""
	}
	return fmt.Sprintf("%.7f",float64(x)*float64(y)/base)
}

func HexWeiToEth(a string)string{
	x:=HexToDec(a)
	if x < 0 {
		return ""
	}
	return fmt.Sprintf("%.7f",float64(x)/base)
}

func IntWeiToEth(a string)string{
	x:=StringToInt64(a)
	if x < 0 {
		return ""
	}
	return fmt.Sprintf("%.7f",float64(x)/base)
}


