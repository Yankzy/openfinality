#!/usr/bin/env bash
set -e

export PATH=${PWD}/../../bin:${PWD}/bin:$PATH
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/fabric-ca/clients

echo "Enrolling OrdererOrg..."
mkdir -p organizations/orderer
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/orderer
fabric-ca-client enroll -u https://admin:adminpw@localhost:10054 --caname ca-orderer --tls.certfiles ${PWD}/organizations/fabric-ca/orderer/ca-cert.pem
mkdir -p organizations/orderer/msp/tlscacerts
cp organizations/fabric-ca/orderer/ca-cert.pem organizations/orderer/msp/tlscacerts/tlsca.ordererorg-cert.pem
# Node OUs
echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > organizations/orderer/msp/config.yaml

# Enroll orderers
for i in 1 2 3 4; do
  fabric-ca-client register --caname ca-orderer --id.name orderer${i} --id.secret orderer${i}pw --id.type orderer --tls.certfiles ${PWD}/organizations/fabric-ca/orderer/ca-cert.pem
  fabric-ca-client enroll -u https://orderer${i}:orderer${i}pw@localhost:10054 --caname ca-orderer -M ${PWD}/organizations/orderer/orderers/orderer${i}/msp --csr.hosts orderer${i}.ordererorg.openfinality.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/orderer/ca-cert.pem
  cp organizations/orderer/msp/config.yaml organizations/orderer/orderers/orderer${i}/msp/config.yaml
  fabric-ca-client enroll -u https://orderer${i}:orderer${i}pw@localhost:10054 --caname ca-orderer -M ${PWD}/organizations/orderer/orderers/orderer${i}/tls --enrollment.profile tls --csr.hosts orderer${i}.ordererorg.openfinality.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/orderer/ca-cert.pem
  cp organizations/orderer/orderers/orderer${i}/tls/tlscacerts/* organizations/orderer/orderers/orderer${i}/tls/ca.crt
  mkdir -p organizations/orderer/orderers/orderer${i}/msp/tlscacerts
  cp organizations/orderer/orderers/orderer${i}/tls/ca.crt organizations/orderer/orderers/orderer${i}/msp/tlscacerts/tlsca.ordererorg-cert.pem
  cp organizations/orderer/orderers/orderer${i}/tls/signcerts/* organizations/orderer/orderers/orderer${i}/tls/server.crt
  cp organizations/orderer/orderers/orderer${i}/tls/keystore/* organizations/orderer/orderers/orderer${i}/tls/server.key
done
fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/orderer/ca-cert.pem
fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:10054 --caname ca-orderer -M ${PWD}/organizations/orderer/users/Admin@ordererorg.openfinality.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/orderer/ca-cert.pem
cp organizations/orderer/msp/config.yaml organizations/orderer/users/Admin@ordererorg.openfinality.com/msp/config.yaml

for org in nigeria ghana morocco; do
  port=""
  if [ "$org" = "nigeria" ]; then port="7054"; fi
  if [ "$org" = "ghana" ]; then port="8054"; fi
  if [ "$org" = "morocco" ]; then port="9054"; fi
  
  echo "Enrolling ${org}Org..."
  mkdir -p organizations/${org}
  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/${org}
  fabric-ca-client enroll -u https://admin:adminpw@localhost:${port} --caname ca-${org} --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  mkdir -p organizations/${org}/msp/tlscacerts
  cp organizations/fabric-ca/${org}/ca-cert.pem organizations/${org}/msp/tlscacerts/tlsca.${org}org-cert.pem
  
  echo "NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-${port}-ca-${org}.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-${port}-ca-${org}.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-${port}-ca-${org}.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-${port}-ca-${org}.pem
    OrganizationalUnitIdentifier: orderer" > organizations/${org}/msp/config.yaml

  fabric-ca-client register --caname ca-${org} --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:${port} --caname ca-${org} -M ${PWD}/organizations/${org}/peers/peer0/msp --csr.hosts peer0.${org}org.openfinality.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  cp organizations/${org}/msp/config.yaml organizations/${org}/peers/peer0/msp/config.yaml
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:${port} --caname ca-${org} -M ${PWD}/organizations/${org}/peers/peer0/tls --enrollment.profile tls --csr.hosts peer0.${org}org.openfinality.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  cp organizations/${org}/peers/peer0/tls/tlscacerts/* organizations/${org}/peers/peer0/tls/ca.crt
  cp organizations/${org}/peers/peer0/tls/signcerts/* organizations/${org}/peers/peer0/tls/server.crt
  cp organizations/${org}/peers/peer0/tls/keystore/* organizations/${org}/peers/peer0/tls/server.key
  
  fabric-ca-client register --caname ca-${org} --id.name ${org}admin --id.secret ${org}adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  fabric-ca-client enroll -u https://${org}admin:${org}adminpw@localhost:${port} --caname ca-${org} -M ${PWD}/organizations/${org}/users/Admin@${org}org.openfinality.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  cp organizations/${org}/msp/config.yaml organizations/${org}/users/Admin@${org}org.openfinality.com/msp/config.yaml

  echo "Enrolling Settlement Adapter for ${org}..."
  fabric-ca-client register --caname ca-${org} --id.name settlement_adapter_${org} --id.secret adapterpw --id.type client --id.attrs "role=settlement-adapter:ecert" --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  fabric-ca-client enroll -u https://settlement_adapter_${org}:adapterpw@localhost:${port} --caname ca-${org} -M ${PWD}/organizations/${org}/users/Adapter@${org}org.openfinality.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/${org}/ca-cert.pem
  cp organizations/${org}/msp/config.yaml organizations/${org}/users/Adapter@${org}org.openfinality.com/msp/config.yaml
done
