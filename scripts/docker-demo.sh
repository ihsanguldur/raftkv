#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

PORTS=(8081 8082 8083 8084 8085)
NAMES=(node1 node2 node3 node4 node5)

echo "building and starting cluster..."
docker compose up -d --build

echo "waiting for leader election..."
leader=""
for i in $(seq 1 30); do
    for j in "${!NAMES[@]}"; do
        if docker compose logs "${NAMES[$j]}" 2>/dev/null | grep -q "became leader"; then
            leader="${NAMES[$j]}"
            leader_port="${PORTS[$j]}"
            break 2
        fi
    done
    sleep 1
done

if [ -z "$leader" ]; then
    echo "no leader elected in time"
    docker compose logs
    exit 1
fi

echo "leader: $leader (host port $leader_port)"

echo "--- PUT via node1 (proxied to the leader if node1 isn't it) ---"
curl -sw '\nHTTP:%{http_code}\n' -X PUT localhost:8081/kv/demo -H 'Content-Type: application/json' -d '{"value":"raft"}'

sleep 1
echo "--- GET from every node ---"
for j in "${!NAMES[@]}"; do
    echo -n "${NAMES[$j]}: "
    curl -s "localhost:${PORTS[$j]}/kv/demo" || echo "(unreachable)"
    echo
done

declare -A baseline
for j in "${!NAMES[@]}"; do
    baseline["${NAMES[$j]}"]=$(docker compose logs "${NAMES[$j]}" 2>/dev/null | grep -c "became leader" || true)
done

echo "--- killing leader ($leader) ---"
docker compose kill "$leader"

echo "waiting for re-election..."
new_leader=""
for i in $(seq 1 30); do
    for j in "${!NAMES[@]}"; do
        name="${NAMES[$j]}"
        if [ "$name" = "$leader" ]; then
            continue
        fi
        count=$(docker compose logs "$name" 2>/dev/null | grep -c "became leader" || true)
        if [ "$count" -gt "${baseline[$name]}" ]; then
            new_leader="$name"
            break 2
        fi
    done
    sleep 1
done

echo "new leader observed: ${new_leader:-none found, check logs}"

echo "--- GET from every surviving node (data should still be consistent) ---"
for j in "${!NAMES[@]}"; do
    if [ "${NAMES[$j]}" = "$leader" ]; then
        continue
    fi
    echo -n "${NAMES[$j]}: "
    curl -s "localhost:${PORTS[$j]}/kv/demo" || echo "(unreachable)"
    echo
done

echo "demo done. 'docker compose down -v' to tear down."