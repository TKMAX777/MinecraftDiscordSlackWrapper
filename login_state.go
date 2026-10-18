package main

import (
	"fmt"
	"sync"
	"time"
)

// Register user logins State
type JoinState struct {
	sync.RWMutex
	State map[string]UserState
}

type UserState struct {
	JoinedAt time.Time
}

func NewJoinState() JoinState {
	return JoinState{State: make(map[string]UserState)}
}

func (j *JoinState) Join(username string) {
	j.Lock()
	if _, exists := j.State[username]; !exists {
		j.State[username] = UserState{JoinedAt: time.Now()}
	}
	j.Unlock()
}

func (j *JoinState) Snapshot() map[string]UserState {
	j.RLock()
	defer j.RUnlock()
	state := make(map[string]UserState, len(j.State))
	for name, user := range j.State {
		state[name] = user
	}
	return state
}

func (j *JoinState) Leave(username string) {
	j.Lock()
	delete(j.State, username)
	j.Unlock()
}

func formatLoginDuration(joinedAt, now time.Time) string {
	if joinedAt.IsZero() {
		return "--"
	}
	minutes := int64(now.Sub(joinedAt) / time.Minute)
	if minutes < 0 {
		minutes = 0
	}
	if minutes >= 24*60 {
		return fmt.Sprintf("%dd %02dh %02dm", minutes/(24*60), (minutes/60)%24, minutes%60)
	}
	if minutes >= 60 {
		return fmt.Sprintf("%dh %02dm", minutes/60, minutes%60)
	}
	return fmt.Sprintf("%dm", minutes)
}
