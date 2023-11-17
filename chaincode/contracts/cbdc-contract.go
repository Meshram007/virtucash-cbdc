package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define key names for options
const nameKey = "name"
const symbolKey = "symbol"
const decimalsKey = "decimals"
const totalSupplyKey = "totalSupply"
const ledgerKey = "CBDC"

// Define objectType names for prefix
const allowancePrefix = "allowance"

// Define key names for options

// CBDCContract provides functions for transferring tokens between accounts
type CBDCContract struct {
	contractapi.Contract
	IsPaused bool
}

// MintTransaction represents a minting transaction
type MintTransaction struct {
	TxID          string `json:"txId"`
	Timestamp     string `json:"timestamp"`
	MinterAccount string `json:"minterAccount"`
	Amount        string `json:"amount"`
}

type Bond struct {
	AssetType       string `json:"assetType"`
	BondName        string `json:"bondName"`
	SecretPhrase    string `json:"secretPhrase"`
	BondValueInCBDC string `json:"bondValueInCBDC"`
	BondID          string `json:"bondID"`
}

// event provides an organized struct for emitting events
type event struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int    `json:"value"`
}

const collectionName string = "CoinCollection"

/** Governance Functions  */

// PauseTokenTransfers pauses token transfers. Only authorized entities can call this function.
func (s *CBDCContract) PauseTokenTransfers(ctx contractapi.TransactionContextInterface) (bool, error) {
	// // Governance check: Only authorized entities can pause token transfers
	// if !s.isAuthorizedEntity(ctx) {
	// 	return false, fmt.Errorf("Unauthorized entity. Pausing token transfers is restricted.")
	// }

	s.IsPaused = true
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("Error getting MSPID: %v", err)
	}
	log.Printf("Token transfers are paused by %s", mspID)
	return true, nil
}

// UnpauseTokenTransfers unpauses token transfers. Only authorized entities can call this function.
func (s *CBDCContract) UnpauseTokenTransfers(ctx contractapi.TransactionContextInterface) (bool, error) {
	// // Governance check: Only authorized entities can unpause token transfers
	// if !s.isAuthorizedEntity(ctx) {
	// 	return false, fmt.Errorf("Unauthorized entity. Unpausing token transfers is restricted.")
	// }

	s.IsPaused = false
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("Error getting MSPID: %v", err)
	}
	log.Printf("Token transfers are unpaused by %s", mspID)
	return true, nil
}

// GetAllAccountsWithOrgs returns a list of all accounts along with their respective organizations on the network
func (s *CBDCContract) GetAllAccountsWithOrgs(ctx contractapi.TransactionContextInterface) (map[string]string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get all accounts: %v", err)
	}
	defer resultsIterator.Close()

	accountsWithOrgs := make(map[string]string)

	// Get the MSP ID of the submitting client identity
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get MSPID: %v", err)
	}

	// Iterate through the results and collect the account names with their respective organizations
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("error iterating over results: %v", err)
		}

		// Assuming the account name is the key in the world state
		accountName := queryResponse.Key
		accountsWithOrgs[accountName] = clientMSPID
	}

	return accountsWithOrgs, nil
}

/** private data collection functions */

// BondExists returns true when asset with given ID exists in private data collection
func (s *CBDCContract) BondExists(ctx contractapi.TransactionContextInterface, BondID string) (bool, error) {

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, BondID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateCBDCBond creates a new instance of Bond
func (s *CBDCContract) CreateCBDCBond(ctx contractapi.TransactionContextInterface, bondID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "CentralBankMSP" {
		exists, err := s.BondExists(ctx, bondID)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset %s already exists", bondID)
		}

		bond := new(Bond)

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", err
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("Please provide the private data of bondName, secretPhrase, bondValueInCBDC, dealerName")
		}

		bondName, exists := transientData["bondName"]
		if !exists {
			return "", fmt.Errorf("The bondName was not specified in transient data. Please try again")
		}
		bond.BondName = string(bondName)

		secretPhrase, exists := transientData["secretPhrase"]
		if !exists {
			return "", fmt.Errorf("The secretPhrase was not specified in transient data. Please try again")
		}
		bond.SecretPhrase = string(secretPhrase)

		bondValueInCBDC, exists := transientData["bondValueInCBDC"]
		if !exists {
			return "", fmt.Errorf("The bondValueInCBDC was not specified in transient data. Please try again")
		}
		bond.BondValueInCBDC = string(bondValueInCBDC)

		bond.AssetType = "Bond"
		bond.BondID = bondID

		bytes, _ := json.Marshal(bond)
		err = ctx.GetStub().PutPrivateData(collectionName, bondID, bytes)
		if err != nil {
			return "", fmt.Errorf("could not able to write the data")
		}
		return fmt.Sprintf("Bond with id %v added successfully", bondID), nil
	} else {
		return fmt.Sprintf("Bond cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}

// GetCBDCBond retrieves an instance of Bond from the private data collection
func (s *CBDCContract) GetCBDCBond(ctx contractapi.TransactionContextInterface, bondID string) (*Bond, error) {
	exists, err := s.BondExists(ctx, bondID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", bondID)
	}

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, bondID)
	if err != nil {
		return nil, err
	}
	bond := new(Bond)

	err = json.Unmarshal(bytes, bond)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal private data collection data to type Bond")
	}

	return bond, nil

}

// RemoveCBDCBond deletes an instance of Bond from the private data collection
func (s *CBDCContract) RemoveCBDCBond(ctx contractapi.TransactionContextInterface, bondID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientOrgID == "CentralBankMSP" {
		exists, err := s.BondExists(ctx, bondID)

		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the asset %s does not exist", bondID)
		}

		return ctx.GetStub().DelPrivateData(collectionName, bondID)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the bond", clientOrgID)
	}
}

// Mint creates new tokens and adds them to minter's account balance
// This function triggers a Transfer event
func (s *CBDCContract) Mint(ctx contractapi.TransactionContextInterface, amount int) error {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Governance check: Ensure token transfers are not paused
	if s.IsPaused {
		return fmt.Errorf("Token transfers are paused. Minting tokens is restricted.")
	}

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "CentralBankMSP" {
		return fmt.Errorf("client is not authorized to mint new tokens")
	}

	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	if amount <= 0 {
		return fmt.Errorf("mint amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance int

	// If minter current balance doesn't yet exist, we'll create it with a current balance of 0
	if currentBalanceBytes == nil {
		currentBalance = 0
	} else {
		currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
	}

	updatedBalance, err := add(currentBalance, amount)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return err
	}

	// Update the totalSupply
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	var totalSupply int

	// If no tokens have been minted, initialize the totalSupply
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	// Add the mint amount to the total supply and update the state
	totalSupply, err = add(totalSupply, amount)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return err
	}

	strAmount := strconv.FormatInt(int64(amount), 10)

	// Get transaction ID
	txID := ctx.GetStub().GetTxID()

	// Get transaction timestamp
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	// Format timestamp
	timestamp := txTimestamp.AsTime()
	formattedTime := timestamp.Format(time.RFC1123)

	// Record the minting transaction in the ledger
	mintTransaction := MintTransaction{
		TxID:          txID,
		Timestamp:     formattedTime,
		MinterAccount: minter,
		Amount:        strAmount,
	}

	// Marshal the mint transaction to JSON
	mintTransactionJSON, err := json.Marshal(mintTransaction)
	if err != nil {
		return fmt.Errorf("failed to marshal mint transaction: %v", err)
	}

	// Append the mint transaction to the ledger
	err = ctx.GetStub().PutState(ledgerKey, mintTransactionJSON)
	if err != nil {
		return fmt.Errorf("failed to record mint transaction: %v", err)
	}

	// Emit the Transfer event
	transferEvent := event{"0x0", minter, amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

	return nil
}

// Burn redeems tokens the minter's account balance
// This function triggers a Transfer event
func (s *CBDCContract) Burn(ctx contractapi.TransactionContextInterface, amount int) error {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Governance check: Ensure token transfers are not paused
	if s.IsPaused {
		return fmt.Errorf("Token transfers are paused. Minting tokens is restricted.")
	}

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "CentralBankMSP" {
		return fmt.Errorf("client is not authorized to mint new tokens")
	}

	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	if amount <= 0 {
		return errors.New("burn amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance int

	// Check if minter current balance exists
	if currentBalanceBytes == nil {
		return errors.New("The balance does not exist")
	}

	currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	updatedBalance, err := sub(currentBalance, amount)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return err
	}

	// Update the totalSupply
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	// If no tokens have been minted, throw error
	if totalSupplyBytes == nil {
		return errors.New("totalSupply does not exist")
	}

	totalSupply, _ := strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// Subtract the burn amount to the total supply and update the state
	totalSupply, err = sub(totalSupply, amount)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return err
	}

	strAmount := strconv.FormatInt(int64(amount), 10)

	// Get transaction ID
	txID := ctx.GetStub().GetTxID()

	// Get transaction timestamp
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	// Format timestamp
	timestamp := txTimestamp.AsTime()
	formattedTime := timestamp.Format(time.RFC1123)

	// Record the minting transaction in the ledger
	mintTransaction := MintTransaction{
		TxID:          txID,
		Timestamp:     formattedTime,
		MinterAccount: minter,
		Amount:        strAmount,
	}

	// Marshal the mint transaction to JSON
	mintTransactionJSON, err := json.Marshal(mintTransaction)
	if err != nil {
		return fmt.Errorf("failed to marshal mint transaction: %v", err)
	}

	// Append the mint transaction to the ledger
	err = ctx.GetStub().PutState(ledgerKey, mintTransactionJSON)
	if err != nil {
		return fmt.Errorf("failed to record mint transaction: %v", err)
	}

	// Emit the Transfer event
	transferEvent := event{minter, "0x0", amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

	return nil
}

// Transfer transfers tokens from client account to recipient account
// recipient account must be a valid clientID as returned by the ClientID() function
// This function triggers a Transfer event
func (s *CBDCContract) Transfer(ctx contractapi.TransactionContextInterface, recipient string, amount int) error {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	err = transferHelper(ctx, clientID, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}

	// Emit the Transfer event
	transferEvent := event{clientID, recipient, amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	return nil
}

// BalanceOf returns the balance of the given account
func (s *CBDCContract) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	balanceBytes, err := ctx.GetStub().GetState(account)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", account)
	}

	balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

// ClientAccountBalance returns the balance of the requesting client's account
func (s *CBDCContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id: %v", err)
	}

	balanceBytes, err := ctx.GetStub().GetState(clientID)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", clientID)
	}

	balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

// ClientAccountID returns the id of the requesting client's account
// In this implementation, the client account ID is the clientId itself
// Users can use this function to get their own account id, which they can then give to others as the payment address
func (s *CBDCContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Get ID of submitting client identity
	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}

	return clientAccountID, nil
}

// TotalSupply returns the total token supply
func (s *CBDCContract) TotalSupply(ctx contractapi.TransactionContextInterface) (int, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Retrieve total supply of tokens from state of smart contract
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	var totalSupply int

	// If no tokens have been minted, return 0
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	log.Printf("TotalSupply: %d tokens", totalSupply)

	return totalSupply, nil
}

// Approve allows the spender to withdraw from the calling client's token account
// The spender can withdraw multiple times if necessary, up to the value amount
// This function triggers an Approval event
func (s *CBDCContract) Approve(ctx contractapi.TransactionContextInterface, spender string, value int) error {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Get ID of submitting client identity
	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Update the state of the smart contract by adding the allowanceKey and value
	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(value)))
	if err != nil {
		return fmt.Errorf("failed to update state of smart contract for key %s: %v", allowanceKey, err)
	}

	// Emit the Approval event
	approvalEvent := event{owner, spender, value}
	approvalEventJSON, err := json.Marshal(approvalEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Approval", approvalEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("client %s approved a withdrawal allowance of %d for spender %s", owner, value, spender)

	return nil
}

// Allowance returns the amount still available for the spender to withdraw from the owner
func (s *CBDCContract) Allowance(ctx contractapi.TransactionContextInterface, owner string, spender string) (int, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return 0, fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Read the allowance amount from the world state
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read allowance for %s from world state: %v", allowanceKey, err)
	}

	var allowance int

	// If no current allowance, set allowance to 0
	if allowanceBytes == nil {
		allowance = 0
	} else {
		allowance, err = strconv.Atoi(string(allowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	log.Printf("The allowance left for spender %s to withdraw from owner %s: %d", spender, owner, allowance)

	return allowance, nil
}

// TransferFrom transfers the value amount from the "from" address to the "to" address
// This function triggers a Transfer event
func (s *CBDCContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	// Get ID of submitting client identity
	spender, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{from, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Retrieve the allowance of the spender
	currentAllowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve the allowance for %s from world state: %v", allowanceKey, err)
	}

	var currentAllowance int
	currentAllowance, _ = strconv.Atoi(string(currentAllowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// Check if transferred value is less than allowance
	if currentAllowance < value {
		return fmt.Errorf("spender does not have enough allowance for transfer")
	}

	// Initiate the transfer
	err = transferHelper(ctx, from, to, value)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}

	// Decrease the allowance
	updatedAllowance, err := sub(currentAllowance, value)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(updatedAllowance)))
	if err != nil {
		return err
	}

	// Emit the Transfer event
	transferEvent := event{from, to, value}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("spender %s allowance updated from %d to %d", spender, currentAllowance, updatedAllowance)

	return nil
}

// Name returns a descriptive name for fungible tokens in this contract
// returns {String} Returns the name of the token

func (s *CBDCContract) Name(ctx contractapi.TransactionContextInterface) (string, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	bytes, err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return "", fmt.Errorf("failed to get Name bytes: %s", err)
	}

	return string(bytes), nil
}

// Symbol returns an abbreviated name for fungible tokens in this contract.
// returns {String} Returns the symbol of the token

func (s *CBDCContract) Symbol(ctx contractapi.TransactionContextInterface) (string, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return "", fmt.Errorf("Contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	bytes, err := ctx.GetStub().GetState(symbolKey)
	if err != nil {
		return "", fmt.Errorf("failed to get Symbol: %v", err)
	}

	return string(bytes), nil
}

// Set information for a token and intialize contract.
// param {String} name The name of the token
// param {String} symbol The symbol of the token
// param {String} decimals The decimals used for the token operations
func (s *CBDCContract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string, decimals string) (bool, error) {

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to intitialize contract
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get MSPID: %v", err)
	}

	// if clientMSPID != "CentralBankMSP" {
	if clientMSPID != "CentralBankMSP" {
		return false, fmt.Errorf("client is not authorized to initialize contract")
	}

	// Check contract options are not already set, client is not authorized to change them once intitialized
	bytes, err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	err = ctx.GetStub().PutState(nameKey, []byte(name))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	err = ctx.GetStub().PutState(symbolKey, []byte(symbol))
	if err != nil {
		return false, fmt.Errorf("failed to set symbol: %v", err)
	}

	err = ctx.GetStub().PutState(decimalsKey, []byte(decimals))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	return true, nil
}

/** Rich Queries Functions */

func (s *CBDCContract) GetCBDCHistory(ctx contractapi.TransactionContextInterface) ([]*MintTransaction, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(ledgerKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("resultsIterator", resultsIterator)

	defer resultsIterator.Close() // check for details why do we need to close it

	var records []*MintTransaction
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var mintingHistory MintTransaction
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &mintingHistory)
			if err != nil {
				return nil, err
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := MintTransaction{
			TxID:          response.TxId,
			Timestamp:     formattedTime,
			MinterAccount: mintingHistory.MinterAccount,
			Amount:        mintingHistory.Amount,
		}
		records = append(records, &record)
	}

	return records, nil
}

func (s *CBDCContract) GetAllCBDCBondsHistory(ctx contractapi.TransactionContextInterface) ([]*Bond, error) {
	queryString := `{"selector":{"assetType":"Bond"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	fmt.Println("Result of database test", resultsIterator, err)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return BondResultIteratorFunction(resultsIterator)
}

func BondResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Bond, error) {
	var bonds []*Bond
	for resultsIterator.HasNext() {
		fmt.Println("Result of database test 2", resultsIterator)
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var order Bond
		err = json.Unmarshal(queryResult.Value, &order)
		if err != nil {
			return nil, err
		}
		bonds = append(bonds, &order)
	}

	return bonds, nil
}

func (s *CBDCContract) GetCombinedHistory(ctx contractapi.TransactionContextInterface) (map[string]interface{}, error) {
	cbdcHistory, err := s.GetCBDCHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CBDC history: %v", err)
	}

	cbdcBondHistory, err := s.GetAllCBDCBondsHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get CBDC Bond history: %v", err)
	}

	result := make(map[string]interface{})
	result["CBDCHistory"] = cbdcHistory
	result["CBDCBondHistory"] = cbdcBondHistory

	return result, nil
}

/** Helper Functions */

// transferHelper is a helper function that transfers tokens from the "from" address to the "to" address
// Dependant functions include Transfer and TransferFrom
func transferHelper(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {

	if from == to {
		return fmt.Errorf("cannot transfer to and from same client account")
	}

	if value < 0 { // transfer of 0 is allowed in ERC-20, so just validate against negative amounts
		return fmt.Errorf("transfer amount cannot be negative")
	}

	fromCurrentBalanceBytes, err := ctx.GetStub().GetState(from)
	if err != nil {
		return fmt.Errorf("failed to read client account %s from world state: %v", from, err)
	}

	if fromCurrentBalanceBytes == nil {
		return fmt.Errorf("client account %s has no balance", from)
	}

	fromCurrentBalance, _ := strconv.Atoi(string(fromCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	if fromCurrentBalance < value {
		return fmt.Errorf("client account %s has insufficient funds", from)
	}

	toCurrentBalanceBytes, err := ctx.GetStub().GetState(to)
	if err != nil {
		return fmt.Errorf("failed to read recipient account %s from world state: %v", to, err)
	}

	var toCurrentBalance int
	// If recipient current balance doesn't yet exist, we'll create it with a current balance of 0
	if toCurrentBalanceBytes == nil {
		toCurrentBalance = 0
	} else {
		toCurrentBalance, _ = strconv.Atoi(string(toCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
	}

	fromUpdatedBalance, err := sub(fromCurrentBalance, value)
	if err != nil {
		return err
	}

	toUpdatedBalance, err := add(toCurrentBalance, value)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(from, []byte(strconv.Itoa(fromUpdatedBalance)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(to, []byte(strconv.Itoa(toUpdatedBalance)))
	if err != nil {
		return err
	}

	log.Printf("client %s balance updated from %d to %d", from, fromCurrentBalance, fromUpdatedBalance)
	log.Printf("recipient %s balance updated from %d to %d", to, toCurrentBalance, toUpdatedBalance)

	return nil
}

// add two number checking for overflow
func add(b int, q int) (int, error) {

	// Check overflow
	var sum int
	sum = q + b

	if (sum < q || sum < b) == (b >= 0 && q >= 0) {
		return 0, fmt.Errorf("Math: addition overflow occurred %d + %d", b, q)
	}

	return sum, nil
}

// Checks that contract options have been already initialized
func checkInitialized(ctx contractapi.TransactionContextInterface) (bool, error) {
	tokenName, err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get token name: %v", err)
	}

	if tokenName == nil {
		return false, nil
	}

	return true, nil
}

// sub two number checking for overflow
func sub(b int, q int) (int, error) {

	// sub two number checking
	if q <= 0 {
		return 0, fmt.Errorf("Error: the subtraction number is %d, it should be greater than 0", q)
	}
	if b < q {
		return 0, fmt.Errorf("Error: the number %d is not enough to be subtracted by %d", b, q)
	}
	var diff int
	diff = b - q

	return diff, nil
}
