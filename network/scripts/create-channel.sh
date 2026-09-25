#!/usr/bin/env bash
set -e

export PATH=${PWD}/../../bin:${PWD}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/config

echo "Generating genesis block for afro-settlement..."
configtxgen -profile AfroSettlementChannel -outputBlock ./genesis.block -channelID afro-settlement

export ORDERER_CA=${PWD}/organizations/orderer/orderers/orderer1/tls/ca.crt
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/orderer/orderers/orderer1/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/orderer/orderers/orderer1/tls/server.key

for i in 1 2 3 4; do
  port=""
  if [ "$i" = "1" ]; then port="7053"; fi
  if [ "$i" = "2" ]; then port="8053"; fi
  if [ "$i" = "3" ]; then port="9053"; fi
  if [ "$i" = "4" ]; then port="10053"; fi
  
  echo "Joining orderer${i} to channel..."
  osnadmin channel join --channelID afro-settlement --config-block ./genesis.block -o localhost:${port} --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
done
