package combatLog

import (
	"strconv"
	"testing"
	"time"

	"../combat"
)

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
