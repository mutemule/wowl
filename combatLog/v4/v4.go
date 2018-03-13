package v4

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mutemule/wowl/combat"
	"github.com/mutemule/wowl/combatLog/event"
)

// Parse picks up after the combat log header to continue parsing the combat log
// XXX: this needs to be broken down a bit more
func Parse(s *bufio.Scanner) (encounters []combat.Encounter, err error) {
	var currentEncounter *combat.Encounter

	for s.Scan() {
		rawCombatEvent := s.Text()
		combatEventTime, combatRecords, err := event.Split(rawCombatEvent)
		if err != nil {
			log.Printf("Failed to parse line '%s':\n", rawCombatEvent)
			return encounters, err
		}

		switch combatRecords[0] {
		case "ENCOUNTER_START":
			encounters = append(encounters, *new(combat.Encounter))
			currentEncounter = &encounters[len(encounters)-1]

			err = startEncounter(combatEventTime, combatRecords, currentEncounter)

		case "ENCOUNTER_END":
			err = endEncounter(combatEventTime, combatRecords, currentEncounter)
			if err != nil {
				return encounters, err
			}
			currentEncounter.Events = append(currentEncounter.Events, rawCombatEvent)
			currentEncounter = new(combat.Encounter)

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

func startEncounter(time time.Time, records []string, encounter *combat.Encounter) (err error) {
	encounter.Start = time
	encounter.Name = records[2]
	encounter.Kill = false

	encounter.ID, err = strconv.Atoi(records[1])
	if err != nil {
		return err
	}

	encounter.DifficultyID, err = strconv.Atoi(records[3])
	if err != nil {
		return err
	}
	encounter.Difficulty = combat.Difficulty[encounter.DifficultyID]

	encounter.RaidSize, err = strconv.Atoi(records[4])
	if err != nil {
		return err
	}

	return nil
}

func endEncounter(time time.Time, records []string, encounter *combat.Encounter) (err error) {
	if encounter == nil || encounter.ID == 0 {
		err = fmt.Errorf("Found an ENCOUNTER_END event without a corresponding ENCOUNTER_START event, ignoring: %s", records)
		return err
	}

	encounter.End = time
	encounter.Kill, err = strconv.ParseBool(records[5])

	return nil
}
