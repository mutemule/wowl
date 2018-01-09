package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CombatLogInfo struct {
	Version         int  `json:"version"`
	AdvancedLogging bool `json:"advancedlogging"`
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

func main() {
	var encounters []Encounter

	combatLogFileName := "C:/Program Files (x86)/World of Warcraft/Logs/WoWCombatLog.txt"
	combatLogFile, err := os.Open(combatLogFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer combatLogFile.Close()

	scanner := bufio.NewScanner(combatLogFile)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	dateTime, combatHeader, _ := parseCombatLogEvent(scanner.Text())

	// Obtain the combat log header
	combatLogInfo, err := parseCombatLogInfo(combatHeader)
	if err != nil {
		log.Fatal(err)
	}

	// Validate combat log version and configuration
	if combatLogInfo.Version != 4 {
		log.Fatalf("Unsupported combat log version: %d", combatLogInfo.Version)
	}

	if combatLogInfo.AdvancedLogging == false {
		log.Print("You need to enable advanced combat logging for full log usage.")
	}
	fmt.Printf("Found valid combat log at %s\n", dateTime)

	var encounter Encounter
	for scanner.Scan() {
		combatEventTime, combatEvent, _ := parseCombatLogEvent(scanner.Text())

		r := csv.NewReader(strings.NewReader(combatEvent))
		r.LazyQuotes = true
		for {
			combatRecords, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(combatEvent)
				log.Fatal(err)
			}

			if combatRecords[0] == "ENCOUNTER_START" {
				encounterID, err := strconv.Atoi(combatRecords[1])
				if err != nil {
					log.Fatal(err)
				}

				difficultyID, err := strconv.Atoi(combatRecords[3])
				if err != nil {
					log.Fatal(err)
				}

				raidSize, err := strconv.Atoi(combatRecords[4])
				if err != nil {
					log.Fatal(err)
				}

				encounter.ID = encounterID
				encounter.Name = combatRecords[2]
				encounter.Start = combatEventTime
				encounter.Difficulty = difficultyID
				encounter.RaidSize = raidSize
				encounter.Kill = false
				encounter.Events = append(encounter.Events, combatEvent)
			} else if combatRecords[0] == "ENCOUNTER_END" {
				encounter.End = combatEventTime
				encounter.Events = append(encounter.Events, combatEvent)
				encounters = append(encounters, encounter)

				fmt.Println(encounter)
			} else if combatRecords[0] == "UNIT_DIED" {
				encounter.Kill = true
			} else if encounter.Name != "" {
				encounter.Events = append(encounter.Events, combatEvent)
			}
		}
	}

	fmt.Println("All done!")
}

func parseCombatLogEvent(s string) (dateStamp time.Time, event string, err error) {
	dateEvent := strings.SplitN(s, "  ", 2)
	dateTime := dateEvent[0]
	event = dateEvent[1]

	currentYear := time.Now().Year()

	layout := "1/2 15:04:05.000 2006"
	dateStamp, err = time.Parse(layout, dateTime+" "+strconv.Itoa(currentYear))

	return dateStamp, event, err
}

func parseCombatLogInfo(s string) (combatLogInfo CombatLogInfo, err error) {
	r := csv.NewReader(strings.NewReader(s))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return combatLogInfo, err
		}

		combatLogVersionString := record[0]
		if combatLogVersionString != "COMBAT_LOG_VERSION" {
			// XXX: this should return an error instead
			log.Fatalf("Expected combat log header, got '%s' instead.", s)
		}
		combatLogVersion, err := strconv.Atoi(record[1])

		advancedLoggingString := record[2]
		if advancedLoggingString != "ADVANCED_LOG_ENABLED" {
			// XXX: this should return an error instead
			log.Fatalf("Expected advanced logging indicator, got '%s' instead.", advancedLoggingString)
		}
		advancedLogging, err := strconv.ParseBool(record[3])

		combatLogInfo.Version = combatLogVersion
		combatLogInfo.AdvancedLogging = advancedLogging
	}

	return combatLogInfo, nil
}
