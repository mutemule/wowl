package main

import (
	"bufio"
	"flag"
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
	info, fights, err := combatLog.Parse(combatLogFileName)
	if err != nil {
		log.Fatalf("Failed to open combat log file: %s\n", err)
	}

	if debug {
		log.Printf("DEBUG: Combat log file to parse: %s\n", combatLogFileName)
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

	log.Printf("Found %d fights total:\n", len(fights))
	for _, fight := range fights {
		fightLength := fight.End.Sub(fight.Start).Round(1 * time.Second)
		fightResult := killOrWipe(fight.Kill)
		log.Printf("%s %s: %s (%s) (%d deaths)\n", fight.Difficulty, fight.Name, fightResult, fightLength, len(fight.Deaths))
		// for _, death := range fight.Deaths {
		// 	relativeDeathTime := death.Time.Sub(fight.Start).Round(1 * time.Second)
		// 	fmt.Printf("  %s died at %s\n", death.Name, relativeDeathTime)
		// }

		for _, event := range fight.Events {
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
