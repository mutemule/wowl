package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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

func main() {
	combatLogFileName := "C:/Program Files (x86)/World of Warcraft/Logs/WoWCombatLog.txt"

	var encounters []Encounter

	combatLogFile, err := os.Open(combatLogFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer combatLogFile.Close()

	scanner := bufio.NewScanner(combatLogFile)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	combatLogTime, combatLogHeaderFields, _ := parseCombatLogEvent(scanner.Text())

	// Obtain the combat log header
	combatLogInfo, err := parseCombatLogHeader(combatLogHeaderFields)
	if err != nil {
		log.Fatal(err)
	}
	combatLogInfo.Time = combatLogTime

	// Validate combat log version and configuration
	if combatLogInfo.Version != 4 {
		log.Fatalf("Unsupported combat log version: %d", combatLogInfo.Version)
	}

	if combatLogInfo.AdvancedLogging == false {
		log.Print("You need to enable advanced combat logging for full log usage.")
	}

	switch combatLogInfo.Version {
	case 4:
		encounters, err = parsev4CombatLog(scanner)
		if err != nil {
			log.Fatal(err)
		}
	}

	// buffer := new(bytes.Buffer)
	// encoder := json.NewEncoder(buffer)
	// encoder.SetIndent("", "\t")

	// err = encoder.Encode(encounters)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(buffer.String())

	fmt.Println("")
	for _, encounter := range encounters {
		encounterLength := encounter.End.Sub(encounter.Start).Round(1 * time.Second)
		encounterResult := killOrWipe(encounter.Kill)
		fmt.Printf("%s of %s lasted %s.\n", encounterResult, encounter.Name, encounterLength)
	}
}

func killOrWipe(k bool) string {
	if k {
		return "Kill"
	}

	return "Wipe"
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
