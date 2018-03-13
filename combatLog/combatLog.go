package combatLog

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/mutemule/wowl/combat"
	"github.com/mutemule/wowl/combatLog/event"
	"github.com/mutemule/wowl/combatLog/v4"
)

// Parse will parse the full combat log and return the appropriate metadata and encounters
func Parse(fileName string) (info combat.Info, encounters []combat.Encounter, err error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return info, encounters, err
	}
	defer fd.Close()

	scanner, err := getScanner(fd)
	if err != nil {
		log.Printf("Error: Failed to open the combat log file: %s\n", err)
		return info, encounters, err
	}

	scanner.Split(bufio.ScanLines)
	scanner.Scan()

	// Obtain the combat log header and log start
	combatTime, logHeaderFields, err := event.Split(scanner.Text())
	if err != nil {
		return info, encounters, err
	}

	// Parse the combat log header
	combatInfo, err := parseHeader(logHeaderFields)
	if err != nil {
		log.Printf("Failed to parse the combat log header '%s':", logHeaderFields)
		log.Fatal(err)
	}
	combatInfo.Time = combatTime

	// The logs are only really useful if advanced logging is enabled
	if combatInfo.AdvancedLogging == false {
		log.Print("You need to enable advanced combat logging for full log usage.")
	}

	switch combatInfo.Version {
	default:
		log.Fatalf("Unsupported combat log version '%d'", combatInfo.Version)
	case 4:
		encounters, err = v4.Parse(scanner)
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

// getScanner takes an open file and returns an appropriate buffered scanner object for that file
// This allows us to easily add support for compressed files
func getScanner(fd *os.File) (scanner *bufio.Scanner, err error) {
	bReader := bufio.NewReader(fd)
	firstTwoBytes, err := bReader.Peek(2)

	if firstTwoBytes[0] == 31 && firstTwoBytes[1] == 139 {
		gzipReader, err := gzip.NewReader(bReader)
		if err != nil {
			return scanner, err
		}
		defer gzipReader.Close()

		// We explode the gzip file because the gzip reader sometimes returns partial lines
		// So rather than track line states ourselves, it's easier to just use a temp file
		// to handle all the log reading
		// A little harder on disk and maybe memory, but considerably easier for us.
		uncompressedFile, err := ioutil.TempFile(os.TempDir(), "wowl-WoWCombatLog")
		defer os.Remove(uncompressedFile.Name())
		_, err = io.Copy(uncompressedFile, gzipReader)
		uncompressedFile.Seek(0, 0)

		scanner = bufio.NewScanner(uncompressedFile)
	} else {
		scanner = bufio.NewScanner(bReader)
	}

	return scanner, err
}
