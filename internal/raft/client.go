package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (n *Node) sendRequestVote(peer string, args RequestVoteArgs) (RequestVoteReply, error) {
	var reply RequestVoteReply

	body, err := json.Marshal(args)
	if err != nil {
		return reply, err
	}

	url := fmt.Sprintf("http://%s/raft/request-vote", peer)
	resp, err := n.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return reply, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (n *Node) sendAppendEntries(peer string, args AppendEntriesArgs) (AppendEntriesReply, error) {
	var reply AppendEntriesReply

	body, err := json.Marshal(args)
	if err != nil {
		return reply, err
	}

	url := fmt.Sprintf("http://%s/raft/append-entries", peer)
	resp, err := n.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return reply, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return reply, err
	}

	return reply, nil
}
