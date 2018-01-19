package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"./combatLog"
	"./combatLog/generic"
)

func main() {
	combatLogFileName := "WoWCombatLog.txt"
	parsedCombatLogFileName := "WoWCombatLogParsed.txt"

	var encounters []generic.Encounter

	combatLogFile, err := os.Open(combatLogFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer combatLogFile.Close()

	scanner := bufio.NewScanner(combatLogFile)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	combatLogTime, combatLogHeaderFields, err := combatLog.ParseEvent(scanner.Text())
	if err != nil {
		log.Printf("Failed to read the combat log header:")
		log.Fatal(err)
	}

	// Obtain the combat log header
	combatLogInfo, err := combatLog.ParseHeader(combatLogHeaderFields)
	if err != nil {
		log.Printf("Failed to parse the combat log header '%s':", combatLogHeaderFields)
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

	encounters, err = combatLog.Parse(combatLogInfo, scanner)

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
		fmt.Printf("%s %s: %s (%s) (%d deaths)\n", encounter.Difficulty, encounter.Name, encounterResult, encounterLength, len(encounter.Deaths))
		// for _, death := range encounter.Deaths {
		// 	relativeDeathTime := death.Time.Sub(encounter.Start).Round(1 * time.Second)
		// 	fmt.Printf("  %s died at %s\n", death.Name, relativeDeathTime)
		// }

		for _, event := range encounter.Events {
			_, err = w.WriteString(event + "\n")
			if err != nil {
				log.Printf("Failed to write combat event: %s", err)
			}
		}
	}
}

// XXX: Add some mechanism to detect a reset
func killOrWipe(k bool) string {
	if k {
		return "Kill"
	}

	return "Wipe"
}
