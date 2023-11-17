# Project: CBDC (Central Bank Digital Coin)

# Step-by-Step Instructions 

**Intall The Below Dependancies First**
```
https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html
```

**Clone The Reposistory**

```
git clone https://github.com/Meshram007/virtucash-cbdc.git
```

```
cd virtucash-cbdc
```

```
cd chaincode
```

**Import and install a Go modules and files**
```
go mod tidy
```
```
go mod vendor
```

**Note**: If go version is of format 1.21.3 then change it to 1.21

```
cd ../../
```

**central-bank-network**

**Start the network**

```
cd central-bank-network
```
**Run the shell script file to start network**

```
./startCoinNetwork.sh
```


**Try Below commands**

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"function":"Transfer","Args":["eDUwOTo6Q049Y29tbWVyY2lhbGJhbmthZG1pbixPVT1hZG1pbixPPUh5cGVybGVkZ2VyLFNUPU5vcnRoIENhcm9saW5hLEM9VVM6OkNOPWZhYnJpYy1jYS1zZXJ2ZXIsT1U9RmFicmljLE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==","10"]}' 

```


```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["BalanceOf", "eDUwOTo6Q049Y2VudHJhbGJhbmthZG1pbixPVT1hZG1pbixPPUh5cGVybGVkZ2VyLFNUPU5vcnRoIENhcm9saW5hLEM9VVM6OkNOPWZhYnJpYy1jYS1zZXJ2ZXIsT1U9RmFicmljLE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw=="]}'
```


```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"function":"Burn","Args":["10"]}'
```

**Implementation of PDC:- Create asset 2 i.e CBDC Bond**
```
export BONDNAME=$(echo -n "RupeeBond" | base64 | tr -d \\n)
export SECRETPHRASE=$(echo -n "Oranges are new black" | base64 | tr -d \\n)
export BONDVALUEINCBDC=$(echo -n "20000000" | base64 | tr -d \\n)
```

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"Args":["CBDCContract:CreateCBDCBond","BOND-01"]}' --transient "{\"bondName\":\"$BONDNAME\",\"secretPhrase\":\"$SECRETPHRASE\",\"bondValueInCBDC\":\"$BONDVALUEINCBDC\"}"
```

```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["CBDCContract:getCBDCBond","BOND-01"]}'
```

**Environment variables for Commercial Bank in new terminal**

```
export CHANNEL_NAME=cbdcchannel
export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=CommercialBankMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/users/Admin@commercialbank coin.com/msp
export CORE_PEER_ADDRESS=localhost:9051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/coin.com/orderers/orderer.coin.com/msp/tlscacerts/tlsca.coin.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt 
```
**Query - Bond-001**
```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["CBDCContract:getCBDCBond","BOND-01"]}'
```

**Environment variables for Consumer in new terminal**

```
export CHANNEL_NAME=cbdcchannel
export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=ConsumerMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/consumer.coin.com/users/Admin@consumer.coin.com/msp
export CORE_PEER_ADDRESS=localhost:11051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/coin.com/orderers/orderer.coin.com/msp/tlscacerts/tlsca.coin.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt 
```

```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["CBDCContract:ReadCBDCBond","BOND-01"]}'
```

**Rich Queries**

**Environment variables for Central Bank**
```
export CHANNEL_NAME=cbdcchannel
export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=CentralBankMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/centralbank.coin.com/users/Admin@centralbank.coin.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/coin.com/orderers/orderer.coin.com/msp/tlscacerts/tlsca.coin.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt
```
**Create CBDC Bond**

```
export BONDNAME=$(echo -n "CryptBond" | base64 | tr -d \\n)
export SECRETPHRASE=$(echo -n "Red are new black" | base64 | tr -d \\n)
export BONDVALUEINCBDC=$(echo -n "30000000" | base64 | tr -d \\n)
```

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"Args":["CBDCContract:CreateCBDCBond","BOND-02"]}' --transient "{\"bondName\":\"$BONDNAME\",\"secretPhrase\":\"$SECRETPHRASE\",\"bondValueInCBDC\":\"$BONDVALUEINCBDC\"}"
```

```
export BONDNAME=$(echo -n "BTCBond" | base64 | tr -d \\n)
export SECRETPHRASE=$(echo -n "White are new black" | base64 | tr -d \\n)
export BONDVALUEINCBDC=$(echo -n "4000000" | base64 | tr -d \\n)
```

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"Args":["CBDCContract:CreateCBDCBond","BOND-03"]}' --transient "{\"bondName\":\"$BONDNAME\",\"secretPhrase\":\"$SECRETPHRASE\",\"bondValueInCBDC\":\"$BONDVALUEINCBDC\"}"
```
**Rich Query**

```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["CBDCContract:GetCBDCHistory"]}' | jq
```

```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["CBDCContract:GetAllCBDCBondsHistory"]}' | jq
```

```
peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["CBDCContract:GetAllAssets"]}' | jq
```

**To view couchdb**

http://localhost:7984/_utils/

userid: `admin`

password: `adminpw`

**Tear down the network**

```
./stopCoinNetwork.sh
```

