#!/usr/bin/env bash
set -e

cd network

echo "Tearing down nodes..."
docker compose -f compose/docker-compose-nodes.yaml down -v --remove-orphans || true
docker compose -f compose/docker-compose-ccaas.yaml down -v --remove-orphans || true

echo "Tearing down CAs..."
docker compose -f compose/docker-compose-base.yaml down -v --remove-orphans || true

rm -rf organizations/fabric-ca/
rm -rf organizations/orderer/msp
rm -rf organizations/orderer/orderers
rm -rf organizations/nigeria/msp organizations/nigeria/peers organizations/nigeria/users
rm -rf organizations/ghana/msp organizations/ghana/peers organizations/ghana/users
rm -rf organizations/morocco/msp organizations/morocco/peers organizations/morocco/users
rm -f genesis.block
