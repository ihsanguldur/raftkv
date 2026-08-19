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

	log         []LogEntry
	commitIndex int
	lastApplied int

	nextIndex  map[string]int
	matchIndex map[string]int

	notifyCh map[int]chan struct{}

	applyFn func(Command)

	httpClient *http.Client
}

func NewNode(id string, peers []string, applyFn func(Command)) *Node {
	return &Node{
		id:          id,
		peers:       peers,
		currentTerm: 0,
		votedFor:    "",
		state:       Follower,
		log:         []LogEntry{},
		commitIndex: 0,
		lastApplied: 0,
		notifyCh:    make(map[int]chan struct{}),
		applyFn:     applyFn,
		httpClient: &http.Client{
			Timeout: 200 * time.Millisecond,
		},
	}
}
