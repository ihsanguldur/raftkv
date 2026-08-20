package raft

func (n *Node) lastLogIndex() int {
	if len(n.log) == 0 {
		return 0
	}

	return n.log[len(n.log)-1].Index
}

func (n *Node) lastLogTerm() int {
	if len(n.log) == 0 {
		return 0
	}

	return n.log[len(n.log)-1].Term
}

func (n *Node) termAtIndex(index int) int {
	if index == 0 || index > len(n.log) {
		return 0
	}
	return n.log[index-1].Term
}

func (n *Node) truncateFrom(index int) {
	if index <= 0 || index > len(n.log) {
		return
	}
	n.log = n.log[:index-1]
}

func (n *Node) appendEntry(entry LogEntry) {
	n.log = append(n.log, entry)
}

func (n *Node) entriesFrom(index int) []LogEntry {
	if index <= 0 || index > len(n.log) {
		return nil
	}
	src := n.log[index-1:]
	dst := make([]LogEntry, len(src))
	copy(dst, src)
	return dst
}
