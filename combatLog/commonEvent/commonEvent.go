package commonEvent

import (
	"encoding/csv"
	"log"
	"strconv"
	"strings"
	"time"

	"../../combat"
)

// ParseEvent takes a single combat log event and returns the datestampe along with a slice of event fields
func ParseEvent(s string) (dateStamp time.Time, events []string, err error) {
	dateEvent := strings.SplitN(s, "  ", 2)
	dateTime := dateEvent[0]

	dateStamp, _ = parseDate(dateTime)

	r := csv.NewReader(strings.NewReader(dateEvent[1]))
	r.LazyQuotes = true
	events, err = r.Read()

	return dateStamp, events, err
}

// ParseHeader takes the slice of header events and returns a struct representing prased values
func ParseHeader(headerFields []string) (combatLogInfo combat.Info, err error) {
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

// parseDate takes a CombatLog-formatted datestamp and returns a full time.Time() struct
func parseDate(s string) (combatEventDate time.Time, err error) {
	layout := "1/2 15:04:05.000 2006"
	currentDate := time.Now()
	currentYear := currentDate.Year()

	combatEventDate, err = time.Parse(layout, s+" "+strconv.Itoa(currentYear))
	if err != nil {
		return combatEventDate, err
	}

	if combatEventDate.After(currentDate) {
		previousYear := currentYear - 1
		combatEventDate, err = time.Parse(layout, s+" "+strconv.Itoa(previousYear))

		return combatEventDate, err
	}

	return combatEventDate, err
}
