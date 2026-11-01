#!/usr/bin/env bash
set -e

export PATH=${PWD}/../../bin:${PWD}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/config

ALREADY_COMMITTED=false
if docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer lifecycle chaincode querycommitted --channelID openfinality-settlement --name openfinality >/dev/null 2>&1; then
  ALREADY_COMMITTED=true
fi

if [ "$ALREADY_COMMITTED" = "false" ]; then
  echo "Packaging CCaaS for each org..."
  cd network/ccaas

  # Nigeria
  cp connection-nigeria.json connection.json
  tar -czf code.tar.gz connection.json
  tar -czf openfinality-nigeria.tar.gz metadata.json code.tar.gz

  # Ghana
  cp connection-ghana.json connection.json
  tar -czf code.tar.gz connection.json
  tar -czf openfinality-ghana.tar.gz metadata.json code.tar.gz

  # Morocco
  cp connection-morocco.json connection.json
  tar -czf code.tar.gz connection.json
  tar -czf openfinality-morocco.tar.gz metadata.json code.tar.gz

  cd ../..

  echo "Copying CCaaS package to peers..."
  docker cp network/ccaas/openfinality-nigeria.tar.gz peer0.nigeriaorg.openfinality.com:/tmp/openfinality.tar.gz
  docker cp network/ccaas/openfinality-ghana.tar.gz peer0.ghanaorg.openfinality.com:/tmp/openfinality.tar.gz
  docker cp network/ccaas/openfinality-morocco.tar.gz peer0.moroccoorg.openfinality.com:/tmp/openfinality.tar.gz

  echo "Installing chaincode on NigeriaOrg..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer lifecycle chaincode install /tmp/openfinality.tar.gz || true

  echo "Installing chaincode on GhanaOrg..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.ghanaorg.openfinality.com peer lifecycle chaincode install /tmp/openfinality.tar.gz || true

  echo "Installing chaincode on MoroccoOrg..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.moroccoorg.openfinality.com peer lifecycle chaincode install /tmp/openfinality.tar.gz || true
fi

echo "Querying installed chaincode..."
docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer lifecycle chaincode queryinstalled > log-nigeria.txt
CC_PACKAGE_ID_NIGERIA=$(grep -o "openfinality_1.0:[a-f0-9]*" log-nigeria.txt | head -n 1)

docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.ghanaorg.openfinality.com peer lifecycle chaincode queryinstalled > log-ghana.txt
CC_PACKAGE_ID_GHANA=$(grep -o "openfinality_1.0:[a-f0-9]*" log-ghana.txt | head -n 1)

docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.moroccoorg.openfinality.com peer lifecycle chaincode queryinstalled > log-morocco.txt
CC_PACKAGE_ID_MOROCCO=$(grep -o "openfinality_1.0:[a-f0-9]*" log-morocco.txt | head -n 1)

echo "Building chaincode docker image..."
docker build -t openfinality-chaincode:latest ./chaincode/openfinality

echo "Starting CCaaS services..."
export CC_PACKAGE_ID_NIGERIA
export CC_PACKAGE_ID_GHANA
export CC_PACKAGE_ID_MOROCCO
docker compose -f network/compose/docker-compose-ccaas.yaml up -d --force-recreate

if [ "$ALREADY_COMMITTED" = "true" ]; then
  echo "Chaincode openfinality already committed on channel openfinality-settlement, skipping approve and commit."
else
  echo "Approving for NigeriaOrg..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer lifecycle chaincode approveformyorg -o orderer1.ordererorg.openfinality.com:7050 --tls --cafile /etc/hyperledger/fabric/tls/orderer/server.crt --channelID openfinality-settlement --name openfinality --version 1.0 --sequence 1 --package-id $CC_PACKAGE_ID_NIGERIA

  echo "Approving for GhanaOrg..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.ghanaorg.openfinality.com peer lifecycle chaincode approveformyorg -o orderer1.ordererorg.openfinality.com:7050 --tls --cafile /etc/hyperledger/fabric/tls/orderer/server.crt --channelID openfinality-settlement --name openfinality --version 1.0 --sequence 1 --package-id $CC_PACKAGE_ID_GHANA

  echo "Approving for MoroccoOrg..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.moroccoorg.openfinality.com peer lifecycle chaincode approveformyorg -o orderer1.ordererorg.openfinality.com:7050 --tls --cafile /etc/hyperledger/fabric/tls/orderer/server.crt --channelID openfinality-settlement --name openfinality --version 1.0 --sequence 1 --package-id $CC_PACKAGE_ID_MOROCCO

  echo "Checking commit readiness..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer lifecycle chaincode checkcommitreadiness --channelID openfinality-settlement --name openfinality --version 1.0 --sequence 1 --tls --cafile /etc/hyperledger/fabric/tls/orderer/server.crt --output json

  echo "Committing chaincode..."
  docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer lifecycle chaincode commit -o orderer1.ordererorg.openfinality.com:7050 --tls --cafile /etc/hyperledger/fabric/tls/orderer/server.crt --channelID openfinality-settlement --name openfinality --version 1.0 --sequence 1 --peerAddresses peer0.nigeriaorg.openfinality.com:7051 --tlsRootCertFiles /etc/hyperledger/fabric/organizations/nigeria/peers/peer0/tls/ca.crt --peerAddresses peer0.ghanaorg.openfinality.com:9051 --tlsRootCertFiles /etc/hyperledger/fabric/organizations/ghana/peers/peer0/tls/ca.crt --peerAddresses peer0.moroccoorg.openfinality.com:11051 --tlsRootCertFiles /etc/hyperledger/fabric/organizations/morocco/peers/peer0/tls/ca.crt
fi
