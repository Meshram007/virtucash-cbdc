# Step-by-Step Instructions 

**Chaincode Folder Structure**

```
mkdir -p KBA-Automobile/Chaincode
```

```
cd KBA-Automobile/Chaincode/
```
**Initialize a Go module**
```
go mod init github.com/kba-chf/kba-automobile/chaincode
```

**Note**: If go version is of format 1.21.3 then change it to 1.21

```
touch main.go
```
```
mkdir contracts
```
```
touch contracts/car-contract.go
```
```
go mod tidy
```
```
cd ../../
```

**Automobile Network**

**Start the network**

```
cd fabric-samples/test-network
```


```
./network.sh up createChannel -c autochannel -ca -s couchdb
```


**Bring up org3**

```
cd addOrg3
```


```
./addOrg3.sh up -c autochannel -ca -s couchdb
```


```
cd ..
```

**Deploy the chaincode**
```
./network.sh deployCC -ccn KBA-Automobile -ccp ../../KBA-Automobile/Chaincode/ -ccl go -c autochannel -ccv 1.0 -ccs 1
```

**General Environment variables**

```
export FABRIC_CFG_PATH=$PWD/../config/
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_TLS_ENABLED=true
```

**Environment variables for Org1**

```
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```
**Invoke - CreateCar**
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C autochannel -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"CreateCar","Args":["Car-01", "Tata", "Nexon", "White", "Factory-1", "22/07/2023" ]}'
```
**Query - ReadCar**

```
peer chaincode query -C autochannel -n KBA-Automobile -c '{"function":"ReadCar","Args":["Car-01"]}'
```

```
peer chaincode query -C autochannel -n KBA-Automobile -c '{"Args":["ReadCar", "Car-01"]}'
```

```
peer chaincode query -C autochannel -n KBA-Automobile -c '{"function":"CarContract:ReadCar","Args":["Car-01"]}'
```
**Testing access control**

**Environment variables for Org2**
```
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```
**Invoke - CreateCar**

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C autochannel -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"CreateCar","Args":["Car-03", "Tata", "Nexon", "White", "Factory-1", "22/07/2023" ]}'
```
**Query - ReadCar**

```
peer chaincode query -C autochannel -n KBA-Automobile -c '{"Args":["ReadCar", "Car-01"]}'
```
**Upgrade chaincode with deleteCar**

```
./network.sh deployCC -ccn KBA-Automobile -ccp ../../KBA-Automobile/Chaincode/ -ccl go -c autochannel -ccv 2.0 -ccs 2
```
**Environment variables for Org1**

```
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

**Invoke - CreateCar**

```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C autochannel -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"CreateCar","Args":["Car-02", "Tata", "Tiago", "Blue", "Factory-1", "22/09/2023" ]}'
```
**Invoke - DeleteCar**
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C autochannel -n KBA-Automobile --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"DeleteCar","Args":["Car-02"]}'
```
**Query - ReadCar**
```
peer chaincode query -C autochannel -n KBA-Automobile -c '{"Args":["ReadCar", "Car-02"]}'
```

**To view couchdb**

http://localhost:7984/_utils/

userid: `admin`

password: `adminpw`

**To view blockchain**

`docker exec -it peer0.org1.example.com /bin/sh`

`ls /var/hyperledger/production/ledgersData/chains/chains/autochannel/`

`cat /var/hyperledger/production/ledgersData/chains/chains/autochannel/blockfile_000000`

`exit`

**To get blockchain information of a specified channel**

`peer channel getinfo -c autochannel`

**Tear down the network**

```
./network.sh down
```


```
docker volume rm $(docker volume ls -q)
```

