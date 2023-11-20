package main

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CBDCData struct {
	TxID          string `json:"txId"`
	Timestamp     string `json:"timestamp"`
	MinterAccount string `json:"minterAccount"`
	Amount        string `json:"amount"`
}

type Bond struct {
	BondName        string `json:"bondName"`
	SecretPhrase    string `json:"secretPhrase"`
	BondValueInCBDC string `json:"bondValueInCBDC"`
	BondID          string `json:"bondID"`
}

func main() {
	router := gin.Default()

	// Example route for Mint transaction
	router.POST("/api/mint/:amount", func(ctx *gin.Context) {
		amount := ctx.Param("amount")
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "Mint", amount)
		ctx.JSON(http.StatusOK, gin.H{"Minted CDBC:": result})
	})

	// Example route for Mint transaction
	router.POST("/api/burn/:amount", func(ctx *gin.Context) {
		amount := ctx.Param("amount")
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "Burn", amount)
		ctx.JSON(http.StatusOK, gin.H{"Burned CDBC:": result})
	})

	// Example route for TotalSupply query
	router.GET("/api/total-supply", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "TotalSupply")
		ctx.JSON(http.StatusOK, gin.H{"Total Supply": result})
	})

	router.POST("/api/transfer/:from/:to/:amount", func(ctx *gin.Context) {
		// Get the parameters from the URL
		from := ctx.Param("from")
		to := ctx.Param("to")
		amount := ctx.Param("amount")

		fmt.Println("testing", from, to, amount)

		// Call the Transfer function with the provided parameters
		result := submitTxnFn(from, "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "Transfer", to, amount)

		ctx.JSON(http.StatusOK, gin.H{"Transferred CBDC": result})
	})

	router.GET("/api/name", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "Name")

		// In the submitTxnFn function, after receiving the response
		fmt.Println("Raw response:", result)

		ctx.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.GET("/api/symbol", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "Symbol")
		ctx.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.POST("/api/pause-token", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "PauseTokenTransfers")
		ctx.JSON(http.StatusOK, gin.H{"Paused Token Functionalities Successfully": result})
	})

	router.POST("/api/unpause-token", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "UnpauseTokenTransfers")
		ctx.JSON(http.StatusOK, gin.H{"Unpaused Token Functionalities Successfully": result})
	})

	router.GET("/api/accounts", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllAccountsWithOrgs")

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(result), &data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse JSON response"})
			return
		}

		// Construct the desired response format
		response := make(map[string]interface{})
		for key, value := range data {
			response[key] = value
		}

		ctx.JSON(http.StatusOK, response)
		// ctx.JSON(http.StatusOK, gin.H{"Get all the accounts": result})
	})

	router.GET("/api/client-bank-balance", func(ctx *gin.Context) {
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "ClientAccountBalance")
		ctx.JSON(http.StatusOK, gin.H{"Client Bank Balance": result})
	})

	router.GET("/api/balanceof/:account", func(ctx *gin.Context) {
		// Get the account parameter from the URL
		account := ctx.Param("account")

		fmt.Println("Tetsing", account)

		// Call the BalanceOf function with the provided account
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "BalanceOf", account)

		// Parse the result to an integer
		balance, err := strconv.Atoi(result)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse balance: %v", err)})
			return
		}

		// Return the balance as JSON
		ctx.JSON(http.StatusOK, gin.H{"balance": balance})
	})

	router.GET("/api/cbdc-history", func(ctx *gin.Context) {
		// Call the GetCBDCHistory function
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetCBDCHistory")

		// Unmarshal the JSON response into a slice of CBDCData
		var cbdcData []CBDCData
		if err := json.Unmarshal([]byte(result), &cbdcData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse JSON response"})
			return
		}

		// Construct the desired response format
		response := make([]CBDCData, len(cbdcData))
		for i, data := range cbdcData {
			response[i] = CBDCData{
				TxID:          data.TxID,
				Timestamp:     data.Timestamp,
				MinterAccount: data.MinterAccount,
				Amount:        data.Amount,
			}
		}

		ctx.JSON(http.StatusOK, response)
	})

	router.POST("/api/bond", func(ctx *gin.Context) {
		var req Bond
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("bond  %s", req)

		privateData := map[string][]byte{
			"bondName":        []byte(req.BondName),
			"secretPhrase":    []byte(req.SecretPhrase),
			"bondValueInCBDC": []byte(req.BondValueInCBDC),
		}

		submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "private", privateData, "CreateCBDCBond", req.BondID)

		// submitTxnFn("dealer", "autochannel", "KBA-Automobile", "OrderContract", "private", privateData, "CreateOrder", req.OrderId)

		ctx.JSON(http.StatusOK, req)
	})

	router.POST("/api/remove-bond/:bondID", func(ctx *gin.Context) {
		bondID := ctx.Param("bondID")
	
		result := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "invoke", make(map[string][]byte), "RemoveCBDCBond", bondID)
	
		if result != "" {
			ctx.JSON(http.StatusOK, gin.H{"message": "CBDC bond removed successfully"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove CBDC bond"})
		}
	})

	router.GET("/api/all-cbdc-bonds-history", func(ctx *gin.Context) {
		result := submitTxnFn("commercialbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllCBDCBondsHistory")
	
		var data []*Bond
		if err := json.Unmarshal([]byte(result), &data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse JSON response"})
			return
		}
	
		ctx.JSON(http.StatusOK, data)
	})

	router.GET("/api/all-assets", func(ctx *gin.Context) {
		cbdcHistoryResult := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetCBDCHistory")
	
		var cbdcHistoryData []CBDCData
		if err := json.Unmarshal([]byte(cbdcHistoryResult), &cbdcHistoryData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse CBDC history JSON response"})
			return
		}
	
		cbdcBondHistoryResult := submitTxnFn("centralbank", "cbdcchannel", "virtucash-cbdc", "CBDCContract", "query", make(map[string][]byte), "GetAllCBDCBondsHistory")
	
		var cbdcBondHistoryData []*Bond
		if err := json.Unmarshal([]byte(cbdcBondHistoryResult), &cbdcBondHistoryData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse CBDC Bond history JSON response"})
			return
		}
	
		result := make(map[string]interface{})
		result["CBDCHistory"] = cbdcHistoryData
		result["CBDCBondHistory"] = cbdcBondHistoryData
	
		ctx.JSON(http.StatusOK, result)
	})

	router.Run("localhost:3000")
}
