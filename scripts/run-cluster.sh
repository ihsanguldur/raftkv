#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

NODES=(node1 node2 node3 node4 node5)
PORTS=(8081 8082 8083 8084 8085)

mkdir -p bin logs

echo "building..."
go build -o bin/node ./cmd/node

pids=()

cleanup() {
    echo "stopping cluster..."
    for pid in "${pids[@]}"; do
        kill "$pid" 2>/dev/null || true
    done
}
trap cleanup EXIT INT TERM

for i in "${!NODES[@]}"; do
    id="${NODES[$i]}"
    port="${PORTS[$i]}"

    peers=""
    for j in "${!NODES[@]}"; do
        if [ "$j" != "$i" ]; then
            peers="${peers}localhost:${PORTS[$j]},"
        fi
    done
    peers="${peers%,}"

    ./bin/node --id "$id" --port "$port" --peers "$peers" > "logs/${id}.log" 2>&1 &
    pids+=($!)
    echo "started $id on :$port (pid $!)"
done

echo "cluster running. tail -f logs/*.log to watch leader election"
echo "press Ctrl+C to stop"

wait