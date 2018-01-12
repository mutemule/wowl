package main

import "testing"

func TestKill(t *testing.T) {
	expectedKill := "Kill"
	returnedKill := killOrWipe(true)

	if expectedKill != returnedKill {
		t.Errorf("Expecting %s, but got %s instead.\n", expectedKill, returnedKill)
	}
}

func TestWipe(t *testing.T) {
	expectedWipe := "Wipe"
	returnedWipe := killOrWipe(false)

	if expectedWipe != returnedWipe {
		t.Errorf("Expecting %s, but got %s instead.\n", expectedWipe, returnedWipe)
	}
}
