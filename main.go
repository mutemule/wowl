package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mutemule/wowl/combatLog"
)

func main() {
	combatLogFileName := "WoWCombatLog.txt"
	// combatLogFileName := "C:/Program Files (x86)/World of Warcraft/Logs/warcraftlogsarchive/WoWCombatLog-archive-2018-01-22T06-33-42.964Z.txt"
	parsedCombatLogFileName := "WoWCombatLogParsed.txt"

	info, encounters, err := combatLog.Parse(combatLogFileName)

	// if debug {
	//   buffer := new(bytes.Buffer)
	//   encoder := json.NewEncoder(buffer)
	//   encoder.SetIndent("", "\t")

	//   err = encoder.Encode(encounters)
	//   if err != nil {
	//  	log.Fatal(err)
	//   }
	//   fmt.Println(buffer.String())
	// }

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
