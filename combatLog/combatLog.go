package combatLog

import (
	"bufio"
	"log"
	"strconv"
	"time"

	"./generic"
	"./v4"
)

// ParseEvent takes a single combat log event and returns the datestampe along with a slice of event fields
func ParseEvent(s string) (dateStamp time.Time, events []string, err error) {
	dateStamp, events, err = generic.ParseEvent(s)

	return dateStamp, events, err
}

// ParseHeader takes the slice of header events and returns a struct representing prased values
func ParseHeader(headerFields []string) (combatLogInfo generic.Info, err error) {
	combatLogVersionField := headerFields[0]
	if combatLogVersionField != "COMBAT_LOG_VERSION" {
		// XXX: this should return an error instead
		log.Fatalf("Expected combat log header, got '%s' instead.", headerFields)
	}
	combatLogVersion, err := strconv.Atoi(headerFields[1])

	advancedLoggingField := headerFields[2]
	if advancedLoggingField != "ADVANCED_LOG_ENABLED" {
		// XXX: this should return an error instead
		log.Fatalf("Expected advanced logging indicator, got '%s' instead.", advancedLoggingField)
	}
	advancedLogging, err := strconv.ParseBool(headerFields[3])

	combatLogInfo.Version = combatLogVersion
	combatLogInfo.AdvancedLogging = advancedLogging

	return combatLogInfo, nil
}

func Parse(combatLogInfo generic.Info, s *bufio.Scanner) (encounters []generic.Encounter, err error) {
	switch combatLogInfo.Version {
	case 4:
		encounters, err = v4.Parsev4CombatLog(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	return encounters, err
}
