package main

import (
	"bufio"
	"strconv"
)

func parsev4CombatLog(s *bufio.Scanner) (encounters []Encounter, err error) {
	var currentEncounter *Encounter

	for s.Scan() {
		rawCombatEvent := s.Text()
		combatEventTime, combatRecords, err := parseCombatLogEvent(rawCombatEvent)
		if err != nil {
			return encounters, err
		}

		// XXX: ew
		if currentEncounter == nil {
			currentEncounter = new(Encounter)
		}

		switch combatRecords[0] {
		case "ENCOUNTER_START":
			encounters = append(encounters, *new(Encounter))
			currentEncounter = &encounters[len(encounters)-1]

			encounterID, err := strconv.Atoi(combatRecords[1])
			if err != nil {
				return encounters, err
			}

			difficultyID, err := strconv.Atoi(combatRecords[3])
			if err != nil {
				return encounters, err
			}

			raidSize, err := strconv.Atoi(combatRecords[4])
			if err != nil {
				return encounters, err
			}

			currentEncounter.ID = encounterID
			currentEncounter.Name = combatRecords[2]
			currentEncounter.Start = combatEventTime
			currentEncounter.Difficulty = difficultyID
			currentEncounter.RaidSize = raidSize
			currentEncounter.Kill = false
			currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
		case "ENCOUNTER_END":
			currentEncounter.End = combatEventTime
			currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
			currentEncounter.Kill, err = strconv.ParseBool(combatRecords[5])
			if err != nil {
				return encounters, err
			}

			currentEncounter = new(Encounter)
		}

		if currentEncounter.ID != 0 {
			currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
		}
	}

	return encounters, err
}
