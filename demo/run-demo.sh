#!/usr/bin/env bash
set -e

echo "OpenFinality Developer Preview (API-driven)"
echo ""

# Build gateway if not exists
if [ ! -f "bin/gateway" ]; then
    make build
fi

# Start Gateway in background if not already running
if ! curl -s http://localhost:8080/healthz > /dev/null; then
    echo "Starting Gateway API..."
    bin/gateway &
    GATEWAY_PID=$!
    sleep 3
else
    GATEWAY_PID=""
fi

INTENT_ID="demo-intent-$(date +%s)"

echo "Settlement Intent (Initiator: NigeriaOrg)"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents \
  -H "X-AfroRail-Org: NigeriaMSP" \
  -H "Content-Type: application/json" \
  -d "{\"id\": \"$INTENT_ID\", \"counterparty_msp_id\": \"MoroccoMSP\"}")
echo $RES
sleep 2

echo "Counterparty (MoroccoOrg)"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/accept \
  -H "X-AfroRail-Org: MoroccoMSP")
echo $RES
sleep 2

echo "Capacity (MoroccoOrg)"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/reservations \
  -H "X-AfroRail-Org: MoroccoMSP")
echo $RES
sleep 2

echo "Settlement Instruction (NigeriaOrg)"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/execute \
  -H "X-AfroRail-Org: NigeriaMSP")
echo $RES
sleep 2

echo "External Settlement (Simulator)"
./demo/simulator.sh execute "$INTENT_ID"
echo "✓ simulator completed"

echo "Settlement Proof (Adapter)"
RES=$(curl -s -X POST http://localhost:8080/v1/settlement-intents/$INTENT_ID/proof \
  -H "X-AfroRail-Org: Adapter")
echo $RES
sleep 2

echo "Final State"
curl -s http://localhost:8080/v1/settlement-intents/$INTENT_ID | jq -r .status

if [ -n "$GATEWAY_PID" ]; then
    kill $GATEWAY_PID
fi
