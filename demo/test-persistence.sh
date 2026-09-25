#!/bin/bash
set -e

echo "Starting Gateway API..."
go build -o bin/gateway cmd/gateway/main.go cmd/gateway/fabric.go
./bin/gateway &
GATEWAY_PID=$!
sleep 5

INTENT_ID="test-persistence-$(date +%s)"
echo "1. Ghana creates intent claiming Ghana initiator"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents \
  -H "X-AfroRail-Org: GhanaMSP" \
  -H "Content-Type: application/json" \
  -d "{\"id\": \"$INTENT_ID\", \"counterparty_msp_id\": \"MoroccoMSP\"}")
if ! echo $RES | grep -q "SUCCESS"; then echo "FAIL: Ghana create ($RES)"; exit 1; fi

echo "2. Morocco accepts Ghana->Morocco intent"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/accept \
  -H "X-AfroRail-Org: MoroccoMSP")
if ! echo $RES | grep -q "SUCCESS"; then echo "FAIL: Morocco accept failed ($RES)"; exit 1; fi

echo "3. Ghana creates reservation"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/reservations \
  -H "X-AfroRail-Org: GhanaMSP")
if ! echo $RES | grep -q "SUCCESS"; then echo "FAIL: Ghana reserve failed ($RES)"; exit 1; fi

echo "4. Ghana execute intent"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/execute \
  -H "X-AfroRail-Org: GhanaMSP")
if ! echo $RES | grep -q "SUCCESS"; then echo "FAIL: Ghana execute failed ($RES)"; exit 1; fi

echo "5. Authorized settlement adapter confirms"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/proof \
  -H "X-AfroRail-Org: Adapter")
if ! echo $RES | grep -q "SUCCESS"; then echo "FAIL: Adapter confirm failed ($RES)"; exit 1; fi

echo "Querying before restart:"
QUERY_RES=$(curl -s -X GET http://localhost:8080/v1/settlement-intents/$INTENT_ID \
  -H "X-AfroRail-Org: GhanaMSP")
echo "Result before restart: $QUERY_RES"

kill $GATEWAY_PID

echo "Restarting network..."
docker compose -f network/compose/docker-compose-nodes.yaml -f network/compose/docker-compose-ccaas.yaml stop

echo "Waiting for containers to fully stop..."
sleep 5

docker compose -f network/compose/docker-compose-nodes.yaml -f network/compose/docker-compose-ccaas.yaml start

echo "Waiting for network to restart and chaincode to re-establish..."
sleep 30

echo "Starting Gateway API again..."
./bin/gateway &
GATEWAY_PID2=$!
sleep 15

echo "Querying after restart:"
QUERY_RES2=$(curl -s -X GET http://localhost:8080/v1/settlement-intents/$INTENT_ID \
  -H "X-AfroRail-Org: GhanaMSP")
echo "Result after restart: $QUERY_RES2"

kill $GATEWAY_PID2

if echo $QUERY_RES2 | grep -q '"status":"FINAL"'; then
  echo "PERSISTENCE TEST PASSED!"
else
  echo "PERSISTENCE TEST FAILED!"
  exit 1
fi
