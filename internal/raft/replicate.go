package raft

import (
	"log"
	"time"
)

const heartbeatInterval = 100 * time.Millisecond

func (n *Node) becomeLeader() {
	n.state = Leader
	term := n.currentTerm
	log.Printf("[%s] became leader for term %d", n.id, term)

	n.nextIndex = make(map[string]int)
	n.matchIndex = make(map[string]int)
	for _, peer := range n.peers {
		n.nextIndex[peer] = n.lastLogIndex() + 1
		n.matchIndex[peer] = 0
	}

	go n.runReplication(term)
}

func (n *Node) runReplication(term int) {
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
			go n.replicateToPeer(peer, id, term)
		}
	}
}

func (n *Node) replicateToPeer(peer, id string, term int) {
	n.mu.Lock()
	if n.state != Leader || n.currentTerm != term {
		n.mu.Unlock()
		return
	}
	nextIdx := n.nextIndex[peer]
	prevLogIndex := nextIdx - 1
	prevLogTerm := n.termAtIndex(prevLogIndex)
	entries := n.entriesFrom(nextIdx)
	leaderCommit := n.commitIndex
	addr := n.addr
	n.mu.Unlock()

	reply, err := n.sendAppendEntries(peer, AppendEntriesArgs{
		Term:         term,
		LeaderID:     id,
		LeaderAddr:   addr,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: leaderCommit,
	})
	if err != nil {
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if reply.Term > n.currentTerm {
		n.currentTerm = reply.Term
		n.state = Follower
		n.votedFor = ""
		n.persist()
		return
	}

	if n.state != Leader || n.currentTerm != term {
		return
	}

	if reply.Success {
		n.matchIndex[peer] = prevLogIndex + len(entries)
		n.nextIndex[peer] = n.matchIndex[peer] + 1
		n.advanceCommitIndex()
		return
	}

	if n.nextIndex[peer] > 1 {
		n.nextIndex[peer]--
	}
}

func (n *Node) advanceCommitIndex() {
	for N := n.lastLogIndex(); N > n.commitIndex; N-- {
		if n.termAtIndex(N) != n.currentTerm {
			continue
		}

		count := 1
		for _, peer := range n.peers {
			if n.matchIndex[peer] >= N {
				count++
			}
		}

		if count > (len(n.peers)+1)/2 {
			n.commitIndex = N
			return
		}
	}
}
