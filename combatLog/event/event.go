package event

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"
)

// Split takes a single combat log event and returns the datestampe along with a slice of event fields
func Split(s string) (dateStamp time.Time, events []string, err error) {
	dateEvent := strings.SplitN(s, "  ", 2)
	dateTime := dateEvent[0]

	dateStamp, _ = parseDate(dateTime)

	r := csv.NewReader(strings.NewReader(dateEvent[1]))
	r.LazyQuotes = true
	events, err = r.Read()

	return dateStamp, events, err
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
