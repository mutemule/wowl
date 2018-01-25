package v4

import (
	"testing"
)

func TestParsingValidCombatStart(t *testing.T) {
	expectedEncounterName := "Portal Keeper Hasabel"
	expectedRaidSize := 25
	combatStartEvent := "1/21 20:43:48.614  ENCOUNTER_START,2064,\"Portal Keeper Hasabel\",17,25"

	encounter, err := startEncounter(combatStartEvent)
	if err != nil {
		t.Errorf("Failed to parse a valid combat start event: %s", err)
	}

	if expectedEncounterName != encounter.Name {
		t.Errorf("Incorrect encounter name identified: expected '%s', got '%s'.", expectedEncounterName, encounter.Name)
	}

	if expectedRaidSize != encounter.RaidSize {
		t.Errorf("Incorrect encounter size identified: expeected %d, got %d.", expectedRaidSize, encounter.RaidSize)
	}
}
