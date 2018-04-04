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
func Parse(s *bufio.Scanner) (fights []combat.Fight, err error) {
	var currentFight *combat.Fight

	for s.Scan() {
		rawCombatEvent := s.Text()
		combatEventTime, combatRecords, err := event.Split(rawCombatEvent)
		if err != nil {
			log.Printf("Failed to parse line '%s':\n", rawCombatEvent)
			return fights, err
		}

		switch combatRecords[0] {
		case "ENCOUNTER_START":
			// XXX: We should probably be using a constructor here
			fights = append(fights, *new(combat.Fight))
			currentFight = &fights[len(fights)-1]
			currentFight.Players = make(map[string]bool)

			err = startEncounter(combatEventTime, combatRecords, currentFight)

		case "ENCOUNTER_END":
			err = endEncounter(combatEventTime, combatRecords, currentFight)
			if err != nil {
				return fights, err
			}
			currentFight.Events = append(currentFight.Events, rawCombatEvent)
			currentFight = new(combat.Fight)

		case "UNIT_DIED":
			if currentFight != nil && currentFight.ID != 0 {
				unitUUID := combatRecords[5]
				unitName := combatRecords[6]

				if strings.HasPrefix(unitUUID, "Player-") {
					playerDeath := combat.UnitDeath{
						Name: unitName,
						Time: combatEventTime,
					}

					currentFight.Deaths = append(currentFight.Deaths, playerDeath)
				}
			}
		}

		if currentFight != nil && currentFight.ID != 0 {
			if (strings.HasPrefix(combatRecords[1], "Player-")) && (combatRecords[0] != "COMBATANT_INFO") {
				playerName := combatRecords[2]
				currentFight.Players[playerName] = true
			}
			currentFight.Events = append(currentFight.Events, rawCombatEvent)
		}
	}

	return fights, err
}

func startEncounter(time time.Time, records []string, encounter *combat.Fight) (err error) {
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

func endEncounter(time time.Time, records []string, encounter *combat.Fight) (err error) {
	if encounter == nil || encounter.ID == 0 {
		err = fmt.Errorf("Found an ENCOUNTER_END event without a corresponding ENCOUNTER_START event, ignoring: %s", records)
		return err
	}

	encounter.End = time
	encounter.Kill, err = strconv.ParseBool(records[5])

	return nil
}
