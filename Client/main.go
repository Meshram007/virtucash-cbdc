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
	// 	"11",
	// )

	// privateData := map[string][]byte{
	// 	"make":       []byte("Maruti"),
	// 	"model":      []byte("Alto"),
	// 	"color":      []byte("Red"),
	// 	"dealerName": []byte("Popular"),
	// }


	result := submitTxnFn(
		"centralbank",
		"cbdcchannel",
		"virtucash-cbdc",
		"CBDCContract",
		"invoke",
		"Mint",
		"11",
	)

	// result := submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "private", privateData, "CreateOrder", "ORD-03")
	// result := submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "ReadOrder", "ORD-03")
	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetAllCars")
	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "OrderContract", "query", make(map[string][]byte), "GetAllOrders")
	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetMatchingOrders", "Car-06")
	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "MatchOrder", "Car-06", "ORD-01")
	// result := submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "RegisterCar", "Car-06", "Dani", "KL-01-CD-01")
	fmt.Println(result)
}
