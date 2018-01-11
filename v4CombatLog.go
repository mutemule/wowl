package main

import (
	"bufio"
	"log"
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
			if currentEncounter == nil || currentEncounter.ID == 0 {
				log.Println("Found an ENCOUNTER_END event without a corresponding ENCOUNTER_START event, ignoring.")
				log.Println(rawCombatEvent)
			} else {
				currentEncounter.End = combatEventTime
				currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
				currentEncounter.Kill, err = strconv.ParseBool(combatRecords[5])
				if err != nil {
					return encounters, err
				}

				currentEncounter = new(Encounter)
			}
		}

		if currentEncounter != nil && currentEncounter.ID != 0 {
			currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
		}
	}

	return encounters, err
}
