#!/usr/bin/env bash
set -e

export PATH=${PWD}/../../bin:${PWD}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/config

echo "Fetching channel block from orderer..."
echo "Waiting for orderer to elect leader and fetch block..."
# Copy orderer TLS CA to peers
docker cp ${PWD}/organizations/orderer/orderers/orderer1/tls/ca.crt peer0.nigeriaorg.openfinality.com:/tmp/orderer_ca.crt
docker cp ${PWD}/organizations/orderer/orderers/orderer1/tls/ca.crt peer0.ghanaorg.openfinality.com:/tmp/orderer_ca.crt
docker cp ${PWD}/organizations/orderer/orderers/orderer1/tls/ca.crt peer0.moroccoorg.openfinality.com:/tmp/orderer_ca.crt

MAX_RETRY=10
for ((i=1; i<=MAX_RETRY; i++)); do
  if docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer channel fetch 0 /tmp/afro-settlement.block -o orderer1.ordererorg.openfinality.com:7050 -c afro-settlement --tls --cafile /tmp/orderer_ca.crt; then
    echo "Successfully fetched block"
    break
  fi
  echo "Retrying ($i/$MAX_RETRY) in 5s..."
  sleep 5
done

# Copy block to other peers using docker cp
docker cp peer0.nigeriaorg.openfinality.com:/tmp/afro-settlement.block ./afro-settlement.block
docker cp ./afro-settlement.block peer0.ghanaorg.openfinality.com:/tmp/afro-settlement.block
docker cp ./afro-settlement.block peer0.moroccoorg.openfinality.com:/tmp/afro-settlement.block

echo "Joining NigeriaOrg peer0 to channel..."
docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.nigeriaorg.openfinality.com peer channel join -b /tmp/afro-settlement.block

echo "Joining GhanaOrg peer0 to channel..."
docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.ghanaorg.openfinality.com peer channel join -b /tmp/afro-settlement.block

echo "Joining MoroccoOrg peer0 to channel..."
docker exec -e CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/admin-msp peer0.moroccoorg.openfinality.com peer channel join -b /tmp/afro-settlement.block
