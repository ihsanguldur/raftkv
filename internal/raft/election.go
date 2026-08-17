package raft

import (
	"math/rand"
	"sync"
	"time"
)

const (
	electionTimeoutMin = 300 * time.Millisecond
	electionTimeoutMax = 600 * time.Millisecond
)

func randomElectionTimeout() time.Duration {
	return electionTimeoutMin + time.Duration(rand.Int63n(int64(electionTimeoutMax-electionTimeoutMin)))
}

func (n *Node) runElectionTimer() {
	timeout := randomElectionTimeout()

	for {
		time.Sleep(10 * time.Millisecond)

		n.mu.Lock()
		elapsed := time.Since(n.electionReset)
		state := n.state
		n.mu.Unlock()

		if state == Leader {
			return
		}

		if elapsed >= timeout {
			n.startElection()
			return
		}
	}
}

func (n *Node) startElection() {
	n.mu.Lock()
	n.state = Candidate
	n.currentTerm++
	term := n.currentTerm
	n.votedFor = n.id
	n.electionReset = time.Now()
	peers := n.peers
	n.mu.Unlock()

	var votes int32 = 1
	var wg sync.WaitGroup

	for _, peer := range peers {
		wg.Add(1)
		go func(peer string) {
			defer wg.Done()

			reply, err := n.sendRequestVote(peer, RequestVoteArgs{Term: term, CandidateID: n.id})
			if err != nil {
				return
			}

			n.mu.Lock()
			defer n.mu.Unlock()

			if reply.Term > n.currentTerm {
				n.currentTerm = reply.Term
				n.state = Follower
				n.votedFor = ""
				return
			}

			if reply.VoteGranted && n.state == Candidate && n.currentTerm == term {
				votes++
				if int(votes) > (len(peers)+1)/2 {
					n.becomeLeader()
				}
			}
		}(peer)
	}

	wg.Wait()

	n.mu.Lock()
	isLeader := n.state == Leader
	n.mu.Unlock()

	if !isLeader {
		go n.runElectionTimer()
	}
}

func (n *Node) Start() {
	go n.runElectionTimer()
}
