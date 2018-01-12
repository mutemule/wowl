package combatLog

import (
	"encoding/csv"
	"log"
	"strconv"
	"strings"
	"time"
)

// Info represents the metadata about the combat log
type Info struct {
	Time            time.Time `json:"time"`
	Version         int       `json:"version"`
	AdvancedLogging bool      `json:"advancedlogging"`
	Header          string    `json:"header"`
}

// Encounter represents all the details about a given encounter
type Encounter struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Start        time.Time   `json:"start"`
	End          time.Time   `json:"end"`
	DifficultyID int         `json:"difficultyID"`
	Difficulty   string      `json:"difficulty"`
	RaidSize     int         `json:"raidSize"`
	Kill         bool        `json:"kill"`
	Deaths       []UnitDeath `json:"deaths"`
	Events       []string    `json:"events"`
}

// UnitDeath records which units died and when
type UnitDeath struct {
	Time time.Time `json:"time"`
	Name string    `json:"name"`
}

// Difficulty maps the numeric encounter difficulty to plain english
var Difficulty = map[int]string{
	0:  "None",
	1:  "5-player",
	2:  "5-player Heroic",
	3:  "10-player Raid",
	4:  "25-player Raid",
	5:  "10-player Heroic Raid",
	6:  "25-player Heroic Raid",
	7:  "LFR",
	8:  "Challenge Mode",
	9:  "40-player Raid",
	11: "Heroic Scenario",
	12: "Scenario",
	14: "Regular",
	15: "Heroic",
	16: "Mythic",
	17: "LFR",
}

// parseDate takes a CombatLog-formatted datestamp and returns a full time.Time() struct
func parseDate(s string) (date time.Time, err error) {
	layout := "1/2 15:04:05.000 2006"
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()

	date, err = time.Parse(layout, s+" "+strconv.Itoa(currentYear))
	if err != nil {
		return date, err
	}

	if date.Month() > currentMonth {
		previousYear := currentYear - 1
		date, err = time.Parse(layout, s+" "+strconv.Itoa(previousYear))
	}

	return date, err
}

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
func ParseHeader(headerFields []string) (combatLogInfo Info, err error) {
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
