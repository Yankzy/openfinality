#!/usr/bin/env bash
set -e

cd network

echo "Waiting for orderer to be ready..."
sleep 10

bash scripts/create-channel.sh
sleep 5
bash scripts/join-channel.sh

echo "Network is up and channel is created."

echo "Deploying chaincode..."
cd ..
bash network/scripts/deploy-chaincode.sh
echo "Chaincode is successfully deployed and initialized."
