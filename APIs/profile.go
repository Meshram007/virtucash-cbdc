package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"centralbank": {
		CertPath:     "../central-bank-network/organizations/peerOrganizations/centralbank.coin.com/users/User1@centralbank.coin.com/msp/signcerts/cert.pem",
		KeyDirectory: "../central-bank-network/organizations/peerOrganizations/centralbank.coin.com/users/User1@centralbank.coin.com/msp/keystore/",
		TLSCertPath:  "../central-bank-network/organizations/peerOrganizations/centralbank.coin.com/peers/peer0.centralbank.coin.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.centralbank.coin.com",
		MSPID:        "CentralBankMSP",
	},

	"commercialbank": {
		CertPath:     "../central-bank-network/organizations/peerOrganizations/commercialbank.coin.com/users/User1@commercialbank.coin.com/msp/signcerts/cert.pem",
		KeyDirectory: "../central-bank-network/organizations/peerOrganizations/commercialbank.coin.com/users/User1@commercialbank.coin.com/msp/keystore/",
		TLSCertPath:  "../central-bank-network/organizations/peerOrganizations/commercialbank.coin.com/peers/peer0.commercialbank.coin.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.commercialbank.coin.com",
		MSPID:        "CommercialBankMSP",
	},

	"consumer": {
		CertPath:     "../central-bank-network/organizations/peerOrganizations/consumer.coin.com/users/User1@consumer.coin.com/msp/signcerts/cert.pem",
		KeyDirectory: "../central-bank-network/organizations/peerOrganizations/consumer.coin.com/users/User1@consumer.coin.com/msp/keystore/",
		TLSCertPath:  "../central-bank-network/organizations/peerOrganizations/consumer.coin.com/peers/peer0.consumer.coin.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.consumer.coin.com",
		MSPID:        "ConsumerMSP",
	},
}
