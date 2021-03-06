package event

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Split takes a single combat log event and returns the datestampe along with a slice of event fields
func Split(s string) (dateStamp time.Time, events []string, err error) {
	if len(s) == 0 {
		err = fmt.Errorf("unable to parse empty event")
		return dateStamp, events, err
	}
	dateEvent := strings.SplitN(s, "  ", 2)
	dateTime := dateEvent[0]

	dateStamp, err = parseDate(dateTime)
	if err != nil {
		return dateStamp, events, err
	}

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
	if combatEventDate.Year() != currentYear {
		err = fmt.Errorf("event: Failed to parse event date of '%s'", s)
	}

	if combatEventDate.After(currentDate) {
		previousYear := currentYear - 1
		combatEventDate, err = time.Parse(layout, s+" "+strconv.Itoa(previousYear))
		if combatEventDate.Year() != previousYear {
			err = fmt.Errorf("event: Failed to parse old event date of '%s'", s)
		}
	}

	return combatEventDate, err
}
