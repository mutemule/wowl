package combatLog

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mutemule/wowl/combat"
)

// A simple integration test just to make sure things are working properly
func TestParsingPlaintextLog(t *testing.T) {
	logfile, _ := filepath.Abs("../test-data/WoWCombatLog-encounters.txt")
	_, _, err := Parse(logfile)
	if err != nil {
		t.Errorf("Failed to parse '%s': %v\n", logfile, err)
	}
}

// A simple integration test just to make sure things are working properly, for gzip'd logs
func TestParsingGzipLog(t *testing.T) {
	logfile, _ := filepath.Abs("../test-data/WoWCombatLog-encounters.txt.gz")
	_, _, err := Parse(logfile)
	if err != nil {
		t.Errorf("Failed to parse '%s': %v\n", logfile, err)
	}
}

func TestParsingValidHeader(t *testing.T) {
	expectedHeader := combat.Info{
		Time:            time.Now().UTC(),
		Version:         4,
		AdvancedLogging: true,
		Header:          "",
	}
	header := make([]string, 4)
	header[0] = "COMBAT_LOG_VERSION"
	header[1] = strconv.Itoa(expectedHeader.Version)
	header[2] = "ADVANCED_LOG_ENABLED"
	header[3] = "1"

	returnedHeader, _ := parseHeader(header)

	if expectedHeader.Version != returnedHeader.Version {
		t.Errorf("Header parsing failed: expected verion %d, but got version %d\n", expectedHeader.Version, returnedHeader.Version)
	}

	if expectedHeader.AdvancedLogging != returnedHeader.AdvancedLogging {
		t.Errorf("Header parsing failed: advanced logging should be %v, but got %v\n", expectedHeader.AdvancedLogging, returnedHeader.AdvancedLogging)
	}
}

func TestParsingNonexistentLogfile(t *testing.T) {
	_, _, err := Parse("nonexistant.txt")

	if err == nil {
		t.Errorf("Somehow parsing a nonexistent logfile succeeded.")
	}

	e, ok := err.(*os.PathError)

	if !ok {
		t.Errorf("We were expecting a PathError, but got '%v' instead.", e)
	}
}

//XXX: This should check to ensure the correct thing is failing
func TestParsingHeaderWithInvalidVersionFieldString(t *testing.T) {
	expectedHeader := combat.Info{
		Time:            time.Now().UTC(),
		Version:         4,
		AdvancedLogging: true,
		Header:          "",
	}
	header := make([]string, 4)
	header[0] = "COMBAT_LOG_VERSION_BROKEN"
	header[1] = strconv.Itoa(expectedHeader.Version)
	header[2] = "ADVANCED_LOG_ENABLED"
	header[3] = "1"

	_, err := parseHeader(header)

	if err == nil {
		t.Error("Expected an error when evaluated an invalid log version field, but got success.")
	}
}

//XXX: This should check to ensure the correct thing is failing
func TestParsingHeaderWithInvalidAdvancedLoggingString(t *testing.T) {
	expectedHeader := combat.Info{
		Time:            time.Now().UTC(),
		Version:         4,
		AdvancedLogging: true,
		Header:          "",
	}
	header := make([]string, 4)
	header[0] = "COMBAT_LOG_VERSION"
	header[1] = strconv.Itoa(expectedHeader.Version)
	header[2] = "ADVANCED_LOG_ENABLED_BROKEN"
	header[3] = "1"

	_, err := parseHeader(header)

	if err == nil {
		t.Error("Expected an error when evaluated an invalid log version field, but got success.")
	}
}

func TestParsingLogWithoutAdvancedLogging(t *testing.T) {
	expectedHeader := combat.Info{
		Time:            time.Now().UTC(),
		Version:         4,
		AdvancedLogging: true,
		Header:          "",
	}
	header := make([]string, 4)
	header[0] = "COMBAT_LOG_VERSION"
	header[1] = strconv.Itoa(expectedHeader.Version)
	header[2] = "ADVANCED_LOG_ENABLED"
	header[3] = "0"

	returnedHeader, err := parseHeader(header)

	if err == nil {
		t.Errorf("We are accepting a combat log without advanced logging: %v", returnedHeader)
	}

	if err.Error() != "advanced logging is not enabled" {
		t.Errorf("Invalid error returned when advanced logging is not enabled: %v", err)
	}
}
