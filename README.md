# raftkv

A distributed key-value store built on top of a **Raft consensus** implementation written from scratch.

The "store" side is deliberately kept trivial (just a `map`) — the real point is seeing how leader election, log replication, and crash recovery actually work in code, not just in theory. Nodes talk to each other over plain HTTP/JSON, so Raft messages (`RequestVote`, `AppendEntries`) can even be triggered manually with `curl`.

## Architecture

Internal layers of a single node:

```mermaid
flowchart LR
    Client(["HTTP Client"])

    subgraph Node["node (cmd/node)"]
        direction TB
        API["internal/api\nclient-facing HTTP"]
        RAFT["internal/raft\nRaft state machine"]
        KV["internal/kv\nkey-value store"]

        API --> RAFT
        RAFT -->|"apply committed entries"| KV
    end

    Client -->|"GET / PUT / DELETE"| API
    RAFT <-->|"RequestVote / AppendEntries\n(HTTP/JSON)"| Peers(["other nodes"])
```

Flow of a write request across a 5-node cluster:

```mermaid
flowchart TD
    C(["Client"]) -->|"PUT key=value"| L

    subgraph Cluster["5-node cluster"]
        L["Leader"]
        F1["Follower"]
        F2["Follower"]
        F3["Follower"]
        F4["Follower"]

        L -->|"AppendEntries"| F1
        L -->|"AppendEntries"| F2
        L -->|"AppendEntries"| F3
        L -->|"AppendEntries"| F4
    end

    F1 -.->|"ack"| L
    F2 -.->|"ack"| L
    L -->|"majority (3/5) acked → commit"| C
```

## Directory layout

```
cmd/node        — node process entrypoint
internal/raft   — Raft state machine (election, log replication, persistence)
internal/kv     — key-value store (the data structure the state machine sits on top of)
internal/api    — client-facing HTTP API
```