package combatLog

import (
	"bufio"
	"log"
	"time"

	"../combat"
	"./commonEvent"
	"./v4Event"
)

// ParseEvent takes a single combat log event and returns the datestampe along with a slice of event fields
func ParseEvent(s string) (dateStamp time.Time, events []string, err error) {
	dateStamp, events, err = commonEvent.ParseEvent(s)

	return dateStamp, events, err
}

// ParseHeader takes the slice of header events and returns a struct representing prased values
func ParseHeader(headerFields []string) (combatLogInfo combat.Info, err error) {
	combatLogInfo, err = commonEvent.ParseHeader(headerFields)

	return combatLogInfo, nil
}

func Parse(combatLogInfo combat.Info, s *bufio.Scanner) (encounters []combat.Encounter, err error) {
	switch combatLogInfo.Version {
	case 4:
		encounters, err = v4Event.Parsev4CombatLog(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	return encounters, err
}
