package raft

import "time"

func (n *Node) HandleRequestVote(args RequestVoteArgs) RequestVoteReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term < n.currentTerm {
		return RequestVoteReply{Term: n.currentTerm, VoteGranted: false}
	}

	changed := false
	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.votedFor = ""
		n.state = Follower
		changed = true
	}

	upToDate := args.LastLogTerm > n.lastLogTerm() ||
		(args.LastLogTerm == n.lastLogTerm() && args.LastLogIndex >= n.lastLogIndex())

	granted := false
	if (n.votedFor == "" || n.votedFor == args.CandidateID) && upToDate {
		n.votedFor = args.CandidateID
		n.electionReset = time.Now()
		changed = true
		granted = true
	}

	if changed {
		n.persist()
	}

	return RequestVoteReply{Term: n.currentTerm, VoteGranted: granted}
}

func (n *Node) HandleAppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term < n.currentTerm {
		return AppendEntriesReply{Term: n.currentTerm, Success: false}
	}

	n.leaderAddr = args.LeaderAddr

	changed := false
	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.votedFor = ""
		changed = true
	}

	if args.PrevLogIndex > 0 && n.termAtIndex(args.PrevLogIndex) != args.PrevLogTerm {
		if changed {
			n.persist()
		}
		return AppendEntriesReply{Term: n.currentTerm, Success: false}
	}

	for _, entry := range args.Entries {
		existingTerm := n.termAtIndex(entry.Index)
		if existingTerm != 0 && existingTerm != entry.Term {
			n.truncateFrom(entry.Index)
			changed = true
		}
		if entry.Index > n.lastLogIndex() {
			n.appendEntry(entry)
			changed = true
		}
	}

	if changed {
		n.persist()
	}

	if args.LeaderCommit > n.commitIndex {
		lastNew := args.PrevLogIndex + len(args.Entries)
		if args.LeaderCommit < lastNew {
			n.commitIndex = args.LeaderCommit
		} else {
			n.commitIndex = lastNew
		}
	}

	n.state = Follower
	n.electionReset = time.Now()

	return AppendEntriesReply{Term: n.currentTerm, Success: true}
}
