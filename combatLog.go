package main

import (
	"encoding/csv"
	"log"
	"strconv"
	"strings"
	"time"
)

type CombatLogInfo struct {
	Time            time.Time `json:"time"`
	Version         int       `json:"version"`
	AdvancedLogging bool      `json:"advancedlogging"`
	Header          string    `json:"header"`
}

type Encounter struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Difficulty int       `json:"difficulty"`
	RaidSize   int       `json:"raidSize"`
	Kill       bool      `json:"kill"`
	Events     []string  `json:"events"`
}

var difficulty = map[int]string{
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
	14: "Raid",
}

func parseCombatLogEvent(s string) (dateStamp time.Time, events []string, err error) {
	dateEvent := strings.SplitN(s, "  ", 2)
	dateTime := dateEvent[0]

	currentYear := time.Now().Year()

	layout := "1/2 15:04:05.000 2006"
	dateStamp, err = time.Parse(layout, dateTime+" "+strconv.Itoa(currentYear))

	r := csv.NewReader(strings.NewReader(dateEvent[1]))
	events, err = r.Read()

	return dateStamp, events, err
}

func parseCombatLogHeader(headerFields []string) (combatLogInfo CombatLogInfo, err error) {
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
