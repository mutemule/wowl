package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mutemule/wowl/combatLog"
)

var combatLogFileName string
var debug bool

func init() {
	flag.StringVar(&combatLogFileName, "combatlog", "C:/Program Files (x86)/World of Warcraft/Logs/WoWCombatLog.txt", "The WoW combat log to parse")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
}

func main() {
	flag.Parse()

	parsedCombatLogFileName := "WoWCombatLogParsed.txt"
	info, encounters, err := combatLog.Parse(combatLogFileName)
	if err != nil {
		log.Fatalf("Failed to open combat log file: %s\n", err)
	}

	if debug {
		fmt.Printf("DEBUG: Combat log file to parse: %s\n", combatLogFileName)
	}

	fh, err := os.Create(parsedCombatLogFileName)
	if err != nil {
		log.Printf("Failed to create prased log file: %s", err)
	}
	w := bufio.NewWriter(fh)
	defer fh.Close()

	_, err = w.WriteString(info.Header + "\n")
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
