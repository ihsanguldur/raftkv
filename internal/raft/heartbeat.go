package raft

import (
	"log"
	"time"
)

const heartbeatInterval = 100 * time.Millisecond

func (n *Node) becomeLeader() {
	n.state = Leader
	term := n.currentTerm
	log.Printf("[%s] become leader for term %d", n.id, term)
	go n.runHeartbeat(term)
}

func (n *Node) runHeartbeat(term int) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		n.mu.Lock()
		if n.state != Leader || n.currentTerm != term {
			n.mu.Unlock()
			return
		}
		peers := n.peers
		id := n.id
		n.mu.Unlock()

		for _, peer := range peers {
			go func(peer string) {
				reply, err := n.sendAppendEntries(peer, AppendEntriesArgs{Term: term, LeaderID: id})
				if err != nil {
					return
				}

				n.mu.Lock()
				defer n.mu.Unlock()
				if reply.Term > n.currentTerm {
					n.currentTerm = reply.Term
					n.state = Follower
					n.votedFor = ""
				}
			}(peer)
		}
	}
}
