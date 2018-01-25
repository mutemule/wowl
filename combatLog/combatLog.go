package combatLog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mutemule/wowl/combat"
	"github.com/mutemule/wowl/combatLog/event"
	"github.com/mutemule/wowl/combatLog/v4"
)

// Parse will parse the full combat log and return the appropriate metadata and encounters
func Parse(combatLogFile string) (info combat.Info, encounters []combat.Encounter, err error) {
	combatLogFileHandle, err := os.Open(combatLogFile)
	if err != nil {
		return info, encounters, err
	}
	defer combatLogFileHandle.Close()

	scanner := bufio.NewScanner(combatLogFileHandle)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	combatTime, logHeaderFields, err := event.Split(scanner.Text())
	if err != nil {
		return info, encounters, err
	}

	// Obtain the combat log header
	combatInfo, err := parseHeader(logHeaderFields)
	if err != nil {
		log.Printf("Failed to parse the combat log header '%s':", logHeaderFields)
		log.Fatal(err)
	}
	combatInfo.Time = combatTime

	// Validate combat log version and configuration, should be its own function
	if combatInfo.Version != 4 {
		log.Fatalf("Unsupported combat log version: %d", combatInfo.Version)
	}

	// The logs are only really useful if advanced logging is enabled
	if combatInfo.AdvancedLogging == false {
		log.Print("You need to enable advanced combat logging for full log usage.")
	}

	switch combatInfo.Version {
	case 4:
		// XXX: stop passing the scanner around and just parse individual events
		// This will be more than a little tricky, but should be doable
		encounters, err = v4.Parse(scanner)
		if err != nil {
			return info, encounters, err
		}
	}

	return info, encounters, err
}

// parseHeader takes the slice of header events and returns a struct representing prased values
func parseHeader(headerFields []string) (combatInfo combat.Info, err error) {
	versionField := headerFields[0]
	if versionField != "COMBAT_LOG_VERSION" {
		err = fmt.Errorf("combatLog: Expected to find COMBAT_LOG_VERSION, found %s instead", versionField)
		return combatInfo, err
	}
	version, err := strconv.Atoi(headerFields[1])

	advancedLoggingField := headerFields[2]
	if advancedLoggingField != "ADVANCED_LOG_ENABLED" {
		err = fmt.Errorf("combatLog: Expected to find ADVANCED_LOG_ENABLED, found %s instead", advancedLoggingField)
		return combatInfo, err
	}
	advancedLogging, err := strconv.ParseBool(headerFields[3])

	combatInfo.Version = version
	combatInfo.AdvancedLogging = advancedLogging

	return combatInfo, nil
}
