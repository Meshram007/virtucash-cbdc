#!/bin/bash

echo "------------Register the ca admin for each organization—----------------"

docker-compose -f docker/docker-compose-ca.yaml up -d
sleep 3

sudo chmod -R 777 organizations/

echo "------------Register and enroll the users for each organization—-----------"

chmod +x registerEnroll.sh

./registerEnroll.sh
sleep 3

echo "—-------------Build the infrastructure—-----------------"

docker-compose -f docker/docker-compose-3org.yaml up -d
sleep 3

echo "-------------Generate the genesis block—-------------------------------"

export FABRIC_CFG_PATH=${PWD}/config

export CHANNEL_NAME=cbdcchannel

configtxgen -profile ThreeOrgsChannel -outputBlock ${PWD}/channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
sleep 2

echo "------ Create the application channel------"

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/coin.com/orderers/orderer.coin.com/msp/tlscacerts/tlsca.coin.com-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/coin.com/orderers/orderer.coin.com/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/coin.com/orderers/orderer.coin.com/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ${PWD}/channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_LOCALMSPID=CentralBankMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/centralbank.coin.com/users/Admin@centralbank.coin.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt
export ORG3_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt
sleep 2

echo "—---------------Join CentralBank peer to the channel—-------------"

echo ${FABRIC_CFG_PATH}
sleep 2
peer channel join -b ${PWD}/channel-artifacts/${CHANNEL_NAME}.block
sleep 3

echo "-----channel List----"
peer channel list

echo "—-------------CentralBank anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json

cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.CentralBankMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.centralbank.coin.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA
sleep 1

echo "—---------------package chaincode—-------------"

peer lifecycle chaincode package cbdccoin.tar.gz --path ${PWD}/../chaincode/ --lang golang --label cbdccoin_1.0
sleep 1

echo "—---------------install chaincode in CentralBank peer—-------------"

peer lifecycle chaincode install cbdccoin.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled
sleep 2

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid cbdccoin.tar.gz)

echo "—---------------Approve chaincode in CentralBank peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --channelID $CHANNEL_NAME --name virtucash-cbdc --version 1.0 --collections-config ../chaincode/collection-coin.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2

export CORE_PEER_LOCALMSPID=CommercialBankMSP 
export CORE_PEER_ADDRESS=localhost:9051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/commercialbank.coin.com/users/Admin@commercialbank.coin.com/msp

echo "—---------------Join CommercialBank peer to the channel—-------------"

peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 2
peer channel list
sleep 2

echo "—-------------CommercialBank anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.CommercialBankMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.commercialbank.coin.com","port": 9051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA
sleep 1

echo "—---------------install chaincode in CommercialBank peer—-------------"

peer lifecycle chaincode install cbdccoin.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled

echo "—---------------Approve chaincode in CommercialBank peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --channelID $CHANNEL_NAME --name virtucash-cbdc --version 1.0 --collections-config ../chaincode/collection-coin.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 1

export CORE_PEER_LOCALMSPID=ConsumerMSP 
export CORE_PEER_ADDRESS=localhost:11051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/consumer.coin.com/users/Admin@consumer.coin.com/msp

echo "—---------------Join Consumer peer to the channel—-------------"

peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
sleep 1
peer channel list

echo "—-------------Consumer anchor peer update—-----------"

peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 1

cd channel-artifacts

configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json

jq '.channel_group.groups.Application.groups.ConsumerMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.consumer.coin.com","port": 11051}]},"version": "0"}}' config_copy.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

cd ..

peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA
sleep 1

peer channel getinfo -c $CHANNEL_NAME

echo "—---------------install chaincode in Consumer peer—-------------"

peer lifecycle chaincode install cbdccoin.tar.gz
sleep 3

peer lifecycle chaincode queryinstalled

echo "—---------------Approve chaincode in Consumer peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --channelID $CHANNEL_NAME --name virtucash-cbdc --version 1.0 --collections-config ../chaincode/collection-coin.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 1

echo "—---------------Commit chaincode in Retailer peer—-------------"

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name virtucash-cbdc --version 1.0 --sequence 1 --collections-config ../chaincode/collection-coin.json --tls --cafile $ORDERER_CA --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --channelID $CHANNEL_NAME --name virtucash-cbdc --version 1.0 --sequence 1 --collections-config ../chaincode/collection-coin.json --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT
sleep 2

peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name virtucash-cbdc --cafile $ORDERER_CA
sleep 2

echo "------------------Set path for CentralBankMSP to invoke functionalities ------------------------"

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
sleep 1


echo "------------------Invoke intialize function and query the functions ------------------------"

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"function":"Initialize","Args":["Central Bank Digital Coin", "CBDC", "18"]}'
sleep 3

peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["Name"]}'
sleep 1

peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["Symbol"]}'
sleep 1

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.coin.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n virtucash-cbdc --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT --peerAddresses localhost:11051 --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT -c '{"function":"Mint","Args":["1000"]}'
sleep 3

peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["TotalSupply"]}'
sleep 1

# peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["PauseTokenTransfers"]}'
# sleep 1

peer chaincode query -C cbdcchannel -n virtucash-cbdc -c '{"Args":["ClientAccountBalance"]}'
sleep 1