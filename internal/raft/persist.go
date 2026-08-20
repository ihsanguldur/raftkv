package raft

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type persistedState struct {
	CurrentTerm int        `json:"current_term"`
	VotedFor    string     `json:"voted_for"`
	Log         []LogEntry `json:"log"`
}

func (n *Node) persistFilePath() string {
	return filepath.Join(n.dataDir, n.id+".json")
}

func (n *Node) persist() {
	state := persistedState{
		CurrentTerm: n.currentTerm,
		VotedFor:    n.votedFor,
		Log:         n.log,
	}

	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("[%s] persist marshal error: %v", n.id, err)
		return
	}

	if err := os.MkdirAll(n.dataDir, 0755); err != nil {
		log.Printf("[%s] persist mkdir error: %v", n.id, err)
		return
	}

	path := n.persistFilePath()
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		log.Printf("[%s] persist write error: %v", n.id, err)
		return
	}
	if err := os.Rename(tmp, path); err != nil {
		log.Printf("[%s] persist rename error: %v", n.id, err)
	}
}

func (n *Node) loadPersisted() {
	data, err := os.ReadFile(n.persistFilePath())
	if err != nil {
		return
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("[%s] persist load error: %v", n.id, err)
		return
	}

	n.currentTerm = state.CurrentTerm
	n.votedFor = state.VotedFor
	n.log = state.Log
}