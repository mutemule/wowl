package v4

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mutemule/wowl/combat"
	"github.com/mutemule/wowl/combatLog/event"
)

// Parse picks up after the combat log header to continue parsing the combat log
// XXX: this needs to be broken down a bit more
func Parse(reader *bufio.Reader) (fights []combat.Fight, err error) {
	var currentFight *combat.Fight

	for {
		rawCombatEvent, err := reader.ReadString('\n')
		if err == io.EOF {
			return fights, nil
		}
		if err != nil {
			log.Fatalf("Unhandled error: %+v", err)
		}

		combatEventTime, combatRecords, err := event.Split(rawCombatEvent)
		if err != nil {
			log.Printf("Failed to parse line '%s':\n", rawCombatEvent)
			return fights, err
		}

		switch combatRecords[0] {
		case "ENCOUNTER_START":
			// err = startEncounter(combatEventTime, combatRecords, currentFight)
			currentFight, err := handleEncounter(reader, rawCombatEvent, "ENCOUNTER_END")
			fights = append(fights, currentFight)
			if err != nil {
				return fights, err
			}

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
}

func handleEncounter(reader *bufio.Reader, initialEvent string, terminatingEvent string) (fight combat.Fight, err error) {
	fight.Players = make(map[string]bool)
	initialEventTime, initialEventRecords, err := event.Split(initialEvent)
	if err != nil {
		return fight, err
	}
	fight.Events = append(fight.Events, initialEvent)
	fight.Start = initialEventTime
	fight.Name = initialEventRecords[2]
	fight.Kill = false

	for fight.End == *new(time.Time) {
		rawCombatEvent, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Print("Ran out of log in the middle of an encounter. Returning what we've got...")
			return fight, nil
		}
		if err != nil {
			log.Fatalf("Error reading log: %+v", err)
		}

		fight.Events = append(fight.Events, rawCombatEvent)
		combatEventTime, combatRecords, err := event.Split(rawCombatEvent)
		if err != nil {
			return fight, err
		}

		switch combatRecords[0] {
		case "ENCOUNTER_END":
			if terminatingEvent == "ENCOUNTER_END" {
				fight.Kill, err = strconv.ParseBool(combatRecords[5])
				fight.End = combatEventTime
			}

		case "UNIT_DIED":
			unitUUID := combatRecords[5]
			unitName := combatRecords[6]

			if strings.HasPrefix(unitUUID, "Player-") {
				playerDeath := combat.UnitDeath{
					Name: unitName,
					Time: combatEventTime,
				}

				fight.Deaths = append(fight.Deaths, playerDeath)
			}
		}

		// Record the player names in here
		if (strings.HasPrefix(combatRecords[1], "Player-")) && (combatRecords[0] != "COMBATANT_INFO") {
			playerName := combatRecords[2]
			fight.Players[playerName] = true
		}

		// XXX: This is horrible hax that we only do because we can't seek in bufio,
		// only in os.File.
		// The datestamp is 17 bytes, followed by two spaces, for 19 bytes of garbage.
		// We want to read in the subsequent 20 bytes and look for fixed strings that indicate
		// some kind of logging snafu
		readAhead, err := reader.Peek(39)
		if err != nil {
			// If we didn't get a full timestamp in the peek, there's no point in continuing
			// This technically leaves a bug if we read 17 or 18 bytes, but that should be handled elsewhere
			if len(readAhead) < 17 {
				log.Print("Re-started logging in the middle of an encounter; the game probably crashed. Doing our best to handle this...")
				fight.End = combatEventTime
			} else {
				readAheadEvent := string(readAhead[19:])
				if strings.HasPrefix(readAheadEvent, "COMBAT_LOG_VERSION") {
					log.Print("Re-started logging in the middle of an encounter; the game probably crashed. Doing our best to handle this...")
					fight.End = combatEventTime
				}

				if strings.HasPrefix(readAheadEvent, "ENCOUNTER_START") {
					log.Print("Attempted to start a new encounter while we're still in our existing one.")
					fight.End = combatEventTime
				}
			}
		}
	}

	return fight, err
}
