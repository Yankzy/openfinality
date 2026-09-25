#!/usr/bin/env bash
set -e

cd network

echo "Cleaning up old data..."
rm -rf organizations/fabric-ca/
rm -rf organizations/orderer/msp
rm -rf organizations/orderer/orderers
rm -rf organizations/nigeria/msp
rm -rf organizations/nigeria/peers
rm -rf organizations/ghana/msp
rm -rf organizations/ghana/peers
rm -rf organizations/morocco/msp
rm -rf organizations/morocco/peers

echo "Starting CAs..."
mkdir -p organizations/fabric-ca/nigeria
mkdir -p organizations/fabric-ca/ghana
mkdir -p organizations/fabric-ca/morocco
mkdir -p organizations/fabric-ca/orderer
docker compose -f compose/docker-compose-base.yaml up -d

echo "Waiting for CAs to start..."
sleep 10

echo "Enrolling identities..."
bash scripts/enroll.sh

echo "Starting nodes..."
docker compose -f compose/docker-compose-nodes.yaml up -d

echo "Waiting for nodes to start..."
sleep 15
echo "Bootstrap complete."
