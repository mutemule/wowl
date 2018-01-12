package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	combatLogFileName := "C:/Program Files (x86)/World of Warcraft/Logs/WoWCombatLog.txt"
	parsedCombatLogFileName := "C:/Program Files (x86)/World of Warcraft/Logs/WoWCombatLogParsed.txt"

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

	fh, err := os.Create(parsedCombatLogFileName)
	if err != nil {
		log.Printf("Failed to create prased log file: %s", err)
	}
	w := bufio.NewWriter(fh)
	defer fh.Close()

	_, err = w.WriteString(combatLogInfo.Header + "\n")
	if err != nil {
		log.Printf("Failed to write combat log header: %s", err)
	}

	for _, encounter := range encounters {
		encounterLength := encounter.End.Sub(encounter.Start).Round(1 * time.Second)
		encounterResult := killOrWipe(encounter.Kill)
		fmt.Printf("%s of %s lasted %s.\n", encounterResult, encounter.Name, encounterLength)

		for _, event := range encounter.Events {
			_, err = w.WriteString(event + "\n")
			if err != nil {
				log.Printf("Failed to write combat event: %s", err)
			}
		}
	}

}

func killOrWipe(k bool) string {
	if k {
		return "Kill"
	}

	return "Wipe"
}
