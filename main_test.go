package main

import "testing"

func TestKill(t *testing.T) {
	expectedKill := "Kill"
	kill := killOrWipe(true)

	if kill != expectedKill {
		t.Errorf("Expecting %s, but got %s instead.\n", expectedKill, kill)
	}
}

func TestWipe(t *testing.T) {
	expectedWipe := "Wipe"
	wipe := killOrWipe(false)

	if wipe != expectedWipe {
		t.Errorf("Expecting %s, but got %s instead.\n", expectedWipe, wipe)
	}
}
