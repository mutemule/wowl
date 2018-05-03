package v4

// func TestParsingValidCombatStart(t *testing.T) {
// 	expectedEncounterName := "Portal Keeper Hasabel"
// 	expectedRaidSize := 25
// 	combatStartEvent := "1/21 20:43:48.614  ENCOUNTER_START,2064,\"Portal Keeper Hasabel\",17,25"

// 	encounter := *new(combat.Fight)
// 	combatEventTime, combatRecords, _ := event.Split(combatStartEvent)

// 	err := startEncounter(combatEventTime, combatRecords, &encounter)
// 	if err != nil {
// 		t.Errorf("Failed to parse a valid combat start event: %s", err)
// 	}

// 	if expectedEncounterName != encounter.Name {
// 		t.Errorf("Incorrect encounter name identified: expected '%s', got '%s'.", expectedEncounterName, encounter.Name)
// 	}

// 	if expectedRaidSize != encounter.RaidSize {
// 		t.Errorf("Incorrect encounter size identified: expeected %d, got %d.", expectedRaidSize, encounter.RaidSize)
// 	}
// }

// func TestParsingValidCombatEnd(t *testing.T) {
// 	combatStartEvent := "1/30 15:36:24.208  ENCOUNTER_START,2076,\"Garothi Worldbreaker\",17,25"
// 	combatEndEvent := "1/30 15:41:11.311  ENCOUNTER_END,2076,\"Garothi Worldbreaker\",17,25,1"
// 	expectedEncounterLength := time.Duration(287) * time.Second

// 	startTime, _, _ := event.Split(combatStartEvent)

// 	encounter := *new(combat.Fight)
// 	encounter.ID = 2076
// 	encounter.Name = "Garothi Worldbreaker"
// 	encounter.Start = startTime

// 	endTime, endEvents, _ := event.Split(combatEndEvent)

// 	err := endEncounter(endTime, endEvents, &encounter)
// 	if err != nil {
// 		t.Errorf("Failed to parse a valid combat end event: %s", err)
// 	}

// 	if !encounter.Kill {
// 		t.Errorf("Failed to register a kill on combat end.")
// 	}

// 	encounterLength := encounter.End.Sub(encounter.Start).Round(1 * time.Second)
// 	if expectedEncounterLength != encounterLength {
// 		t.Errorf("Incorrect encounter length: expected %v, got %v.", expectedEncounterLength, encounterLength)
// 	}
// }
