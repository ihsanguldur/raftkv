package raft

import (
	"net/http"
	"sync"
	"time"
)

type Node struct {
	mu sync.Mutex

	id    string
	peers []string

	currentTerm   int
	votedFor      string
	state         State
	electionReset time.Time

	httpClient *http.Client
}

func NewNode(id string, peers []string) *Node {
	return &Node{
		id:          id,
		peers:       peers,
		currentTerm: 0,
		votedFor:    "",
		state:       Follower,
		httpClient: &http.Client{
			Timeout: 200 * time.Millisecond,
		},
	}
}
