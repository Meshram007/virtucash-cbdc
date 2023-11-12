mkdir -p pharmaceutical/Chaincode

cd pharmaceutical/Chaincode/

go mod init github.com/Assignment_3/pharmaceutical/chaincode


touch contracts/medicine-contract.go


./network.sh up createChannel -c pharmachannel -ca -s couchdb


cd addOrg3

./addOrg3.sh up -c pharmachannel -ca -s couchdb


cd ..


./network.sh deployCC -ccn pharmaceutical -ccp ../../pharmaceutical/chaincode/ -ccl go -c pharmachannel -ccv 1.0 -ccs 1


General Environment variables:-
export FABRIC_CFG_PATH=$PWD/../config/
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_TLS_ENABLED=true


Environment variables for Org1:-
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C pharmachannel -n pharmaceutical --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"AddMedicine","Args":["M-01", "Paracetamol", "20", "22/07/2023", "22/07/2025" ]}' 

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C pharmachannel -n pharmaceutical --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"AddMedicine","Args":["M-02", "Dianapr", "10", "22/07/2022", "22/07/2026" ]}' 

peer chaincode query -C pharmachannel -n pharmaceutical -c '{"function":"GetCarHistory","Args":["M-02"]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C pharmachannel -n pharmaceutical --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"function":"DeleteMedicine","Args":["M-02"]}' 

peer chaincode query -C pharmachannel -n pharmaceutical -c '{"function":"GetCarHistory","Args":["M-02"]}'

peer chaincode query -C pharmachannel -n pharmaceutical -c '{"function":"ReadMedicine","Args":["M-01"]}'

peer chaincode query -C pharmachannel -n pharmaceutical -c '{"Args":["ReadMedicine", "M-01"]}'



export NEW_CC_PACKAGE_ID=kbaautomobile_2.0:90cd042a04b9b835a05bca4447663e8f62c2899dcd8a54eee186d6406ec1a05d

===========================================================================================================================================
/////Open a command terminal with in Pharma-network folder, let's call this terminal as host terminal

cd Pharma-network/

############## host terminal ##############

------------Register the ca admin for each organization—----------------

//Build the docker-compose-ca.yaml in the docker folder

docker-compose -f docker/docker-compose-ca.yaml up -d

sudo chmod -R 777 organizations/

------------Register and enroll the users for each organization—-----------

//Build the registerEnroll.sh script file

chmod +x registerEnroll.sh

./registerEnroll.sh

—-------------Build the infrastructure—-----------------

//Build the docker-compose-3org.yaml in the docker folder

docker-compose -f docker/docker-compose-3org.yaml up -d

-------------Generate the genesis block—-------------------------------

//Build the configtx.yaml file in the config folder

export FABRIC_CFG_PATH=./config

export CHANNEL_NAME=pharmachannel

configtxgen -profile ThreeOrgsChannel -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME


------ Create the application channel------

export ORDERER_CA=./organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY

osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY


/////Open another terminal with in pharma-network folder, let's call this terminal as peer0_Producer terminal.

############## peer0_Producer terminal ##############

// Build the core.yaml in peercfg folder


export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=pharmachannel
export CORE_PEER_LOCALMSPID=ProducerMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/producer.pharma.com/users/Admin@producer.pharma.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt
export ORG4_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/ca.crt

—---------------Join peer to the channel—-------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list

/////Open another terminal with in Pharma-network folder, let's call this terminal as peer0_Supplier terminal.

############## peer0_Supplier terminal ##############

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=pharmachannel 
export CORE_PEER_LOCALMSPID=SupplierMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier.pharma.com/users/Admin@supplier.pharma.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt
export ORG4_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/ca.crt

—---------------Join peer to the channel—-------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list

/////Open another terminal with in Pharma-network folder, let's call this terminal as peer0_Wholesaler terminal.

############## peer0_Wholesaler terminal ##############

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=pharmachannel 
export CORE_PEER_LOCALMSPID=WholesalerMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/users/Admin@wholesaler.pharma.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt
export ORG4_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/ca.crt

—---------------Join peer to the channel—-------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list


/////Open another terminal with in Pharma-network folder, let's call this terminal as peer0_Retailer terminal.

############## peer0_Retailer terminal ##############

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=pharmachannel 
export CORE_PEER_LOCALMSPID=RetailerMSP 
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:12051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/retailer.pharma.com/users/Admin@retailer.pharma.com/msp
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt
export ORG4_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/ca.crt

—---------------Join peer to the channel—-------------

peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block

peer channel list

—-------------anchor peer update—-----------

############## peer0_Producer terminal ##############

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.ManufacturerMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.producer.pharma.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA

############## peer0_Supplier terminal ##############

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.SupplierMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.supplier.pharma.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA

############## peer0_Wholesaler terminal ##############

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.WholesalerMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.wholesaler.pharma.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA

peer channel getinfo -c $CHANNEL_NAME


############## peer0_Retailer terminal ##############

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.WholesalerMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.wholesaler.pharma.com","port": 12051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA

peer channel getinfo -c $CHANNEL_NAME

—-----------------Chaincode lifecycle—-------------------

//Build the chaincode (Change MSPId and collection file)

/// Make sure that pharmaceutical chaincode is available in Chaincode folder which is in the same location of Pharma-network. 

############## peer0_Producer terminal ##############

peer lifecycle chaincode package kbapharma.tar.gz --path ../Chaincode/ --lang golang --label kbaauto_1.0

peer lifecycle chaincode install kbapharma.tar.gz

peer lifecycle chaincode queryinstalled

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbapharma.tar.gz)

############## peer0_Supplier terminal ##############

peer lifecycle chaincode install kbapharma.tar.gz

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbapharma.tar.gz)


############## peer0_Wholesaler terminal ##############

peer lifecycle chaincode install kbapharma.tar.gz

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbapharma.tar.gz)

############## peer0_Retailer terminal ##############

peer lifecycle chaincode install kbapharma.tar.gz

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid kbapharma.tar.gz)



############## peer0_Producer terminal ##############


peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --channelID $CHANNEL_NAME --name pharmaceutical --version 1.0 --collections-config ../Chaincode/collection-pharma.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent

############## peer0_Supplier terminal ##############

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --channelID $CHANNEL_NAME --name pharmaceutical --version 1.0 --collections-config ../Chaincode/collection-pharma.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent

############## peer0_Wholesaler terminal ##############

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --channelID $CHANNEL_NAME --name pharmaceutical --version 1.0 --collections-config ../Chaincode/collection-pharma.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent

############## peer0_Retailer terminal ##############

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --channelID $CHANNEL_NAME --name pharmaceutical --version 1.0 --collections-config ../Chaincode/collection-pharma.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent



############## peer0_Producer terminal ##############


peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name pharmaceutical --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-pharma.json --tls --cafile $ORDERER_CA --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --channelID $CHANNEL_NAME --name pharmaceutical --version 1.0 --sequence 1 --collections-config ../Chaincode/collection-pharma.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT

peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name pharmaceutical --cafile $ORDERER_CA

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n pharmaceutical --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"function":"AddMedicine","Args":["M-01", "Paracetamol", "20", "22/07/2023", "22/07/2025" ]}'

peer chaincode query -C $CHANNEL_NAME -n pharmaceutical -c '{"function":"GetCarHistory","Args":["M-02"]}'


--------Invoke Private Transaction----------

############## peer0_Supplier terminal ##############

export MAKE=$(echo -n "Tata" | base64 | tr -d \\n)

export MODEL=$(echo -n "Tiago" | base64 | tr -d \\n)

export COLOR=$(echo -n "White" | base64 | tr -d \\n)

export DEALER_NAME=$(echo -n "XXX" | base64 | tr -d \\n)

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n pharmaceutical --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"Args":["OrderContract:CreateOrder","ORD201"]}' --transient "{\"make\":\"$MAKE\",\"model\":\"$MODEL\",\"color\":\"$COLOR\",\"dealerName\":\"$DEALER_NAME\"}"

peer chaincode query -C $CHANNEL_NAME -n pharmaceutical -c '{"Args":["OrderContract:ReadOrder","ORD201"]}'


--------- Stop the Pharma-network --------------

############## host terminal ##############

docker-compose -f docker/docker-compose-3org.yaml down

docker-compose -f docker/docker-compose-ca.yaml down

docker rm -f $(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')

docker volume rm $(docker volume ls -q)

sudo rm -rf organizations/

sudo rm -rf channel-artifacts/

sudo rm kbapharma.tar.gz

docker ps -a

// if there still exists the containers then execute the following commands.

docker rm $(docker container ls -q) --force

docker container prune

docker system prune

docker volume prune

docker network prune

///Run using startAutomobileNetwork.sh script

//Build startAutomobileNetwork.sh script file

chmod +x startAutomobileNetwork.sh

./startAutomobileNetwork.sh

//To submit transaction as ManufacturerMSP

export CHANNEL_NAME=pharmachannel
export FABRIC_CFG_PATH=./peercfg
export CORE_PEER_LOCALMSPID=ManufacturerMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/producer.pharma.com/users/Admin@producer.pharma.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.pharma.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n pharmaceutical --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"function":"CreateCar","Args":["Car-101", "Tata", "Nexon", "White", "Factory-1", "22/07/2023"]}'

peer chaincode query -C $CHANNEL_NAME -n pharmaceutical -c '{"Args":["GetAllCars"]}'

//To stop the network using script file

//Build stopAutomobileNetwork.sh script file

chmod +x stopAutomobileNetwork.sh

./stopAutomobileNetwork.sh






