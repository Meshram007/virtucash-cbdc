package main

import "fmt"

func main() {

	// use this functions to evaluate and submit txns
	// try calling these functions

	// result := submitTxnFn(
	// 	"centralbank",
	// 	"cbdcchannel",
	// 	"virtucash-cbdc",
	// 	"CBDCContract",
	// 	"invoke",
	// 	make(map[string][]byte),
	// 	"Mint",
	// 	"13",
	// )

	privateData := map[string][]byte{
		"bondName":        []byte("CryptBond"),
		"secretPhrase":    []byte("Red are new black"),
		"bondValueInCBDC": []byte("30000000"),
	}

	// result := submitTxnFn(
	// 	"centralbank",
	// 	"cbdcchannel",
	// 	"virtucash-cbdc",
	// 	"CBDCContract",
	// 	"invoke",
	// 	"Mint",
	// 	"11",
	// )

	// result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "TotalSupply")
	//result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetCBDCHistory")
	// result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllCBDCBondsHistory")

	// result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "BalanceOf", "eDUwOTo6Q049Y2VudHJhbGJhbmthZG1pbixPVT1hZG1pbixPPUh5cGVybGVkZ2VyLFNUPU5vcnRoIENhcm9saW5hLEM9VVM6OkNOPWZhYnJpYy1jYS1zZXJ2ZXIsT1U9RmFicmljLE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==")
	// result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "Transfer", "eDUwOTo6Q049Y2VudHJhbGJhbmthZG1pbixPVT1hZG1pbixPPUh5cGVybGVkZ2VyLFNUPU5vcnRoIENhcm9saW5hLEM9VVM6OkNOPWZhYnJpYy1jYS1zZXJ2ZXIsT1U9RmFicmljLE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==","10")

	result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "private", privateData, "CreateCBDCBond", "BOND-001")
	// result := submitTxnFn("commercialbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllCBDCBondsHistory")
	// result := submitTxnFn("consumer", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllCBDCBondsHistory")

	// result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllAssets")

	//result := submitTxnFn("consumer", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "ClientAccountBalance")

	fmt.Println(result)
}
