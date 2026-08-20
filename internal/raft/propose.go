package raft

import (
	"errors"
	"time"
)

var ErrNotLeader = errors.New("not leader")
var ErrCommitTimeout = errors.New("commit timeout")

const commitTimeout = 2 * time.Second

func (n *Node) Propose(cmd Command) error {
	n.mu.Lock()

	if n.state != Leader {
		n.mu.Unlock()
		return ErrNotLeader
	}

	entry := LogEntry{
		Term:    n.currentTerm,
		Index:   n.lastLogIndex() + 1,
		Command: cmd,
	}
	n.appendEntry(entry)
	n.persist()

	ch := make(chan struct{})
	n.notifyCh[entry.Index] = ch
	n.mu.Unlock()

	select {
	case <-ch:
		return nil
	case <-time.After(commitTimeout):
		n.mu.Lock()
		delete(n.notifyCh, entry.Index)
		n.mu.Unlock()
		return ErrCommitTimeout
	}
}
