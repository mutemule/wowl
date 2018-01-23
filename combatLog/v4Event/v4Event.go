package v4Event

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/mutemule/wowl/combat"
	"github.com/mutemule/wowl/combatLog/commonEvent"
)

// Parsev4CombatLog XXX: this needs to be broken down a bit more
func Parsev4CombatLog(s *bufio.Scanner) (encounters []combat.Encounter, err error) {
	var currentEncounter *combat.Encounter

	for s.Scan() {
		rawCombatEvent := s.Text()
		combatEventTime, combatRecords, err := commonEvent.ParseEvent(rawCombatEvent)
		if err != nil {
			log.Printf("Failed to parse line '%s':\n", rawCombatEvent)
			return encounters, err
		}

		switch combatRecords[0] {
		case "ENCOUNTER_START":
			encounters = append(encounters, *new(combat.Encounter))
			currentEncounter = &encounters[len(encounters)-1]

			*currentEncounter, err = parseEncounterStart(rawCombatEvent)

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

				currentEncounter = new(combat.Encounter)
			}

		case "UNIT_DIED":
			if currentEncounter != nil && currentEncounter.ID != 0 {
				unitUUID := combatRecords[5]
				unitName := combatRecords[6]

				if strings.HasPrefix(unitUUID, "Player-") {
					playerDeath := combat.UnitDeath{
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

func parseEncounterStart(combatEvent string) (encounter combat.Encounter, err error) {
	time, records, err := commonEvent.ParseEvent(combatEvent)

	encounter.Start = time
	encounter.Kill = false
	encounter.Events = append(encounter.Events, combatEvent)

	encounter.ID, err = strconv.Atoi(records[1])
	if err != nil {
		return encounter, err
	}

	encounter.Name = records[2]

	encounter.DifficultyID, err = strconv.Atoi(records[3])
	if err != nil {
		return encounter, err
	}
	encounter.Difficulty = combat.Difficulty[encounter.DifficultyID]

	encounter.RaidSize, err = strconv.Atoi(records[4])
	if err != nil {
		return encounter, err
	}

	return encounter, err
}
