package main

import "sync"

// Register user logins State
type JoinState struct {
	sync.RWMutex
	State map[string]UserState
}

// for future functions
type UserState struct{}

func NewJoinState() JoinState {
	return JoinState{State: make(map[string]UserState)}
}

func (j *JoinState) Join(username string) {
	j.Lock()
	j.State[username] = UserState{}
	j.Unlock()
}

func (j *JoinState) Leave(username string) {
	j.Lock()
	delete(j.State, username)
	j.Unlock()
}
