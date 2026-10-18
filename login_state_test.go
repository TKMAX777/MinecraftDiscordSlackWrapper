package main

import (
	"testing"
	"time"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
)

func TestLoginSessions(t *testing.T) {
	j := NewJoinState()
	before := time.Now()
	j.Join("Alex")
	first := j.Snapshot()["Alex"].JoinedAt
	if first.Before(before) || first.After(time.Now()) {
		t.Fatal("join time not recorded")
	}
	j.Join("Alex")
	if !j.Snapshot()["Alex"].JoinedAt.Equal(first) {
		t.Fatal("duplicate join reset login time")
	}
	copy := j.Snapshot()
	delete(copy, "Alex")
	if len(j.Snapshot()) != 1 {
		t.Fatal("snapshot mutated live state")
	}
	j.Leave("Alex")
	if len(j.Snapshot()) != 0 {
		t.Fatal("logout did not clear session")
	}
	j.Join("Alex")
	if j.Snapshot()["Alex"].JoinedAt.Before(first) {
		t.Fatal("rejoin used older session")
	}
}

func TestLoginDuration(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		elapsed time.Duration
		want    string
	}{
		{0, "0m"}, {59 * time.Second, "0m"}, {time.Minute, "1m"},
		{59 * time.Minute, "59m"}, {60 * time.Minute, "1h 00m"},
		{125 * time.Minute, "2h 05m"}, {25*time.Hour + 3*time.Minute, "1d 01h 03m"},
		{-time.Minute, "0m"},
	} {
		if got := formatLoginDuration(now.Add(-tc.elapsed), now); got != tc.want {
			t.Fatalf("%v: got %q want %q", tc.elapsed, got, tc.want)
		}
	}
	if formatLoginDuration(time.Time{}, now) != "--" {
		t.Fatal("unknown login time must not show a fabricated duration")
	}
}

func TestRosterRefreshOnTick(t *testing.T) {
	messages := make(chan minecraft.Message)
	ticks := make(chan time.Time)
	joined := make(chan struct{}, 1)
	refreshed := make(chan UserState, 1)
	done := make(chan struct{})
	state := NewJoinState()
	go func() {
		defer close(done)
		relayMinecraftMessages(messages, ticks, &state, []MessageSender{func(m minecraft.Message) error { joined <- struct{}{}; return nil }}, []func() error{func() error { refreshed <- state.Snapshot()["Alex"]; return nil }})
	}()
	messages <- minecraft.Message{Type: minecraft.MessageTypeJoin, User: "Alex"}
	select {
	case <-joined:
	case <-time.After(5 * time.Second):
		t.Fatal("join not processed")
	}
	start := state.Snapshot()["Alex"].JoinedAt
	ticks <- time.Now()
	select {
	case got := <-refreshed:
		if !got.JoinedAt.Equal(start) {
			t.Fatal("refresh lost original join time")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("tick did not refresh roster")
	}
	close(messages)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("relay did not stop")
	}
}
