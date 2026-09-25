#!/usr/bin/env bash
set -e

SIM_FILE="/tmp/openfinality-sim.json"

if [ ! -f "$SIM_FILE" ]; then
  cat << 'EOF' > "$SIM_FILE"
{
  "balances": {
    "NigeriaBank": 1000000,
    "MoroccoBank": 5000000
  },
  "executions": {}
}
EOF
fi

ACTION=$1
INTENT_ID=$2

case $ACTION in
  "reserve")
    echo "Reserved capacity for $INTENT_ID"
    ;;
  "execute")
    if jq -e ".executions[\"$INTENT_ID\"]" "$SIM_FILE" > /dev/null 2>&1; then
      echo "Already executed (idempotent)"
    else
      tmp=$(mktemp)
      jq ".balances.NigeriaBank -= 5000 | .executions[\"$INTENT_ID\"] = \"SETTLED\"" "$SIM_FILE" > "$tmp" && mv "$tmp" "$SIM_FILE"
      echo "Execution completed"
    fi
    ;;
  "status")
    status=$(jq -r ".executions[\"$INTENT_ID\"] // \"UNKNOWN\"" "$SIM_FILE")
    echo "$status"
    ;;
  *)
    echo "Usage: simulator.sh [reserve|execute|status] <id>"
    exit 1
    ;;
esac
