package v4

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"../generic"
)

// Parsev4CombatLog XXX: this needs to be broken down a bit more
func Parsev4CombatLog(s *bufio.Scanner) (encounters []generic.Encounter, err error) {
	var currentEncounter *generic.Encounter

	for s.Scan() {
		rawCombatEvent := s.Text()
		combatEventTime, combatRecords, err := generic.ParseEvent(rawCombatEvent)
		if err != nil {
			log.Printf("Failed to parse line '%s':\n", rawCombatEvent)
			return encounters, err
		}

		switch combatRecords[0] {
		case "ENCOUNTER_START":
			encounters = append(encounters, *new(generic.Encounter))
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
			currentEncounter.DifficultyID = difficultyID
			currentEncounter.Difficulty = generic.Difficulty[difficultyID]
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

				currentEncounter = new(generic.Encounter)
			}

		case "UNIT_DIED":
			if currentEncounter != nil && currentEncounter.ID != 0 {
				unitUUID := combatRecords[5]
				unitName := combatRecords[6]

				if strings.HasPrefix(unitUUID, "Player-") {
					playerDeath := generic.UnitDeath{
						Name: unitName,
						Time: combatEventTime,
					}

					currentEncounter.Deaths = append(currentEncounter.Deaths, playerDeath)
				}
			}
		}

		if currentEncounter != nil && currentEncounter.ID != 0 {
			currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
		}
	}

	return encounters, err
}
