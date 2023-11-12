#!/bin/bash

function createProducer() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/producer.pharma.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/producer.pharma.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-producer --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-producer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-producer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-producer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-producer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/producer.pharma.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy producer's CA cert to producer's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/producer.pharma.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/producer/ca-cert.pem" "${PWD}/organizations/peerOrganizations/producer.pharma.com/msp/tlscacerts/ca.crt"

  # Copy producer's CA cert to producer's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/producer.pharma.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/producer/ca-cert.pem" "${PWD}/organizations/peerOrganizations/producer.pharma.com/tlsca/tlsca.producer.pharma.com-cert.pem"

  # Copy producer's CA cert to producer's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/producer.pharma.com/ca"
  cp "${PWD}/organizations/fabric-ca/producer/ca-cert.pem" "${PWD}/organizations/peerOrganizations/producer.pharma.com/ca/ca.producer.pharma.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-producer --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname ca-producer --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-producer --id.name produceradmin --id.secret produceradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-producer -M "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/producer.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-producer -M "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls" --enrollment.profile tls --csr.hosts peer0.producer.pharma.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/producer.pharma.com/peers/peer0.producer.pharma.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-producer -M "${PWD}/organizations/peerOrganizations/producer.pharma.com/users/User1@producer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/producer.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/producer.pharma.com/users/User1@producer.pharma.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://produceradmin:produceradminpw@localhost:7054 --caname ca-producer -M "${PWD}/organizations/peerOrganizations/producer.pharma.com/users/Admin@producer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/producer/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/producer.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/producer.pharma.com/users/Admin@producer.pharma.com/msp/config.yaml"
}

function createSupplier() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/supplier.pharma.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/supplier.pharma.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca-supplier --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-supplier.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/supplier.pharma.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy supplier's CA cert to supplier's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/supplier.pharma.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem" "${PWD}/organizations/peerOrganizations/supplier.pharma.com/msp/tlscacerts/ca.crt"

  # Copy supplier's CA cert to supplier's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/supplier.pharma.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem" "${PWD}/organizations/peerOrganizations/supplier.pharma.com/tlsca/tlsca.supplier.pharma.com-cert.pem"

  # Copy supplier's CA cert to supplier's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/supplier.pharma.com/ca"
  cp "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem" "${PWD}/organizations/peerOrganizations/supplier.pharma.com/ca/ca.supplier.pharma.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-supplier --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname ca-supplier --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-supplier --id.name supplieradmin --id.secret supplieradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-supplier -M "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-supplier -M "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls" --enrollment.profile tls --csr.hosts peer0.supplier.pharma.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/supplier.pharma.com/peers/peer0.supplier.pharma.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-supplier -M "${PWD}/organizations/peerOrganizations/supplier.pharma.com/users/User1@supplier.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier.pharma.com/users/User1@supplier.pharma.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://supplieradmin:supplieradminpw@localhost:8054 --caname ca-supplier -M "${PWD}/organizations/peerOrganizations/supplier.pharma.com/users/Admin@supplier.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/supplier/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/supplier.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/supplier.pharma.com/users/Admin@supplier.pharma.com/msp/config.yaml"
}

function createWholesaler() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/wholesaler.pharma.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:11054 --caname ca-wholesaler --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-wholesaler.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-wholesaler.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-wholesaler.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-11054-ca-wholesaler.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy wholesaler's CA cert to wholesaler's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem" "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/msp/tlscacerts/ca.crt"

  # Copy wholesaler's CA cert to wholesaler's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem" "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/tlsca/tlsca.wholesaler.pharma.com-cert.pem"

  # Copy wholesaler's CA cert to wholesaler's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/ca"
  cp "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem" "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/ca/ca.wholesaler.pharma.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-wholesaler --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname ca-wholesaler --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-wholesaler --id.name wholesaleradmin --id.secret wholesaleradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-wholesaler -M "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:11054 --caname ca-wholesaler -M "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls" --enrollment.profile tls --csr.hosts peer0.wholesaler.pharma.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/peers/peer0.wholesaler.pharma.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:11054 --caname ca-wholesaler -M "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/users/User1@wholesaler.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/users/User1@wholesaler.pharma.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://wholesaleradmin:wholesaleradminpw@localhost:11054 --caname ca-wholesaler -M "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/users/Admin@wholesaler.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/wholesaler/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/wholesaler.pharma.com/users/Admin@wholesaler.pharma.com/msp/config.yaml"
}

function createRetailer() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/retailer.pharma.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/retailer.pharma.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:12054 --caname ca-retailer --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-retailer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-retailer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-retailer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-12054-ca-retailer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/retailer.pharma.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy retailer's CA cert to retailer's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/peerOrganizations/retailer.pharma.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem" "${PWD}/organizations/peerOrganizations/retailer.pharma.com/msp/tlscacerts/ca.crt"

  # Copy retailer's CA cert to retailer's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/retailer.pharma.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem" "${PWD}/organizations/peerOrganizations/retailer.pharma.com/tlsca/tlsca.retailer.pharma.com-cert.pem"

  # Copy retailer's CA cert to retailer's /ca directory (for use by clients)
  mkdir -p "${PWD}/organizations/peerOrganizations/retailer.pharma.com/ca"
  cp "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem" "${PWD}/organizations/peerOrganizations/retailer.pharma.com/ca/ca.retailer.pharma.com-cert.pem"

  echo "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-retailer --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering user"
  set -x
  fabric-ca-client register --caname ca-retailer --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-retailer --id.name wholesaleradmin --id.secret wholesaleradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:12054 --caname ca-retailer -M "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/retailer.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/msp/config.yaml"

  echo "Generating the peer0-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:12054 --caname ca-retailer -M "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls" --enrollment.profile tls --csr.hosts peer0.retailer.pharma.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
  cp "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/retailer.pharma.com/peers/peer0.retailer.pharma.com/tls/server.key"

  echo "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:12054 --caname ca-retailer -M "${PWD}/organizations/peerOrganizations/retailer.pharma.com/users/User1@retailer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/retailer.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/retailer.pharma.com/users/User1@retailer.pharma.com/msp/config.yaml"

  echo "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://wholesaleradmin:wholesaleradminpw@localhost:12054 --caname ca-retailer -M "${PWD}/organizations/peerOrganizations/retailer.pharma.com/users/Admin@retailer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/retailer/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/retailer.pharma.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/retailer.pharma.com/users/Admin@retailer.pharma.com/msp/config.yaml"
}

function createOrderer() {
  echo "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/pharma.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/pharma.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/pharma.com/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy orderer org's CA cert to orderer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/pharma.com/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem"

  # Copy orderer org's CA cert to orderer org's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/ordererOrganizations/pharma.com/tlsca"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/pharma.com/tlsca/tlsca.pharma.com-cert.pem"

  echo "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/pharma.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/config.yaml"

  echo "Generating the orderer-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls" --enrollment.profile tls --csr.hosts orderer.pharma.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/server.key"

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/pharma.com/orderers/orderer.pharma.com/msp/tlscacerts/tlsca.pharma.com-cert.pem"

  echo "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/pharma.com/users/Admin@pharma.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/pharma.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/pharma.com/users/Admin@pharma.com/msp/config.yaml"
}

createProducer
createSupplier
createWholesaler
createRetailer
createOrderer
