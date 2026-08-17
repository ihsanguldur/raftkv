package raft

import "time"

func (n *Node) HandleRequestVote(args RequestVoteArgs) RequestVoteReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term < n.currentTerm {
		return RequestVoteReply{Term: n.currentTerm, VoteGranted: false}
	}

	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.votedFor = ""
		n.state = Follower
	}

	if n.votedFor == "" || n.votedFor == args.CandidateID {
		n.votedFor = args.CandidateID
		n.electionReset = time.Now()
		return RequestVoteReply{Term: n.currentTerm, VoteGranted: true}
	}

	return RequestVoteReply{Term: n.currentTerm, VoteGranted: false}
}

func (n *Node) HandleAppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term < n.currentTerm {
		return AppendEntriesReply{Term: n.currentTerm, Success: false}
	}

	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.votedFor = ""
	}

	n.state = Follower
	n.electionReset = time.Now()

	return AppendEntriesReply{Term: n.currentTerm, Success: true}
}
