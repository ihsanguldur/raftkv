package raft

import "time"

func (n *Node) runApplyLoop() {
	for {
		time.Sleep(10 * time.Millisecond)

		n.mu.Lock()
		var toApply []LogEntry
		if n.commitIndex > n.lastApplied {
			toApply = n.entriesFrom(n.lastApplied + 1)
			count := n.commitIndex - n.lastApplied
			toApply = toApply[:count]
			n.lastApplied = n.commitIndex
		}
		n.mu.Unlock()

		for _, entry := range toApply {
			n.applyFn(entry.Command)

			n.mu.Lock()
			if ch, ok := n.notifyCh[entry.Index]; ok {
				close(ch)
				delete(n.notifyCh, entry.Index)
			}
			n.mu.Unlock()
		}
	}
}
