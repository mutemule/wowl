package event

import (
	"fmt"
	"testing"
	"time"
)

func TestParsingValidEventDateStamp(t *testing.T) {
	dateString := "1/10 15:02:15.348"
	expectedDate := time.Date(time.Now().Year(), 1, 10, 15, 02, 15, 348000000, time.UTC)

	header := dateString + "  COMBAT_LOG_VERSION,4,ADVANCED_LOG_ENABLED,1"
	returnedDate, _, _ := Split(header)

	if expectedDate != returnedDate {
		t.Errorf("Date stamp parsing failed: expected '%s', but got '%s'\n", expectedDate, returnedDate)
	}
}

func TestParsingOldDateStamp(t *testing.T) {
	currentDate := time.Now().UTC()
	expectedDate := time.Date(currentDate.Year()-1, currentDate.Month(), currentDate.Day()+7, 15, 02, 15, 348000000, time.UTC)

	dateString := fmt.Sprintf("%d/%d %02d:%02d:%02d.348",
		int(expectedDate.Month()),
		expectedDate.Day(),
		expectedDate.Hour(),
		expectedDate.Minute(),
		expectedDate.Second(),
	)

	header := dateString + "  COMBAT_LOG_VERSION,4,ADVANCED_LOG_ENABLED,1"
	returnedDate, _, _ := Split(header)

	if expectedDate != returnedDate {
		t.Errorf("Date stamp parsing failed: expected '%s', but got '%s'\n", expectedDate, returnedDate)
	}
}

func TestParsingBasicCombatEvent(t *testing.T) {
	event := "1/10 15:02:15.540  SPELL_DAMAGE,Player-3694-07121BB2,\"Kerliah-Lightbringer\",0x514,0x0,Creature-0-3137-1712-27282-127233-0001566D00,\"Flameweaver\",0xa48,0x0,194153,\"Lunar Strike\",0x40,Creature-0-3137-1712-27282-127233-0001566D00,0000000000000000,23252757,199327752,0,0,1,0,0,0,-3053.49,10534.20,111,248907,-1,64,0,0,0,nil,nil,nil"
	expectedNumberOfEvents := 34
	_, returnedEvents, _ := Split(event)

	if expectedNumberOfEvents != len(returnedEvents) {
		t.Errorf("Event parsing failed: expected %d fields, but got %d.\n", expectedNumberOfEvents, len(returnedEvents))
	}
}

func TestParsingEmbeddedQuoteInCombatEvent(t *testing.T) {
	event := "1/7 16:54:31.838  SPELL_CAST_SUCCESS,Player-47-0781E030,\"Islen-Eitrigg\",0x10512,0x0,Creature-0-3777-1116-32565-77310-00005293F9,\"Mad \\\"King\\\" Sporeon\",0xa48,0x0,204197,\"Purge the Wicked\",0x4,Player-47-0781E030,0000000000000000,41258,41258,0,966,0,51000,51000,1020,1723.37,-738.11,388"
	expectedNumberOfEvents := 25
	_, returnedEvents, err := Split(event)

	if err != nil {
		t.Errorf("Failed to parse event with embeded quotes: %s", err)
	}

	if expectedNumberOfEvents != len(returnedEvents) {
		t.Errorf("Event parsing failed: expected %d fields, but got %d.\n", expectedNumberOfEvents, len(returnedEvents))
	}
}

func TestParsingInvalidEventDateStamp(t *testing.T) {
	dateString := "1/10 15:02:15.348 1980"
	header := dateString + "  COMBAT_LOG_VERSION,4,ADVANCED_LOG_ENABLED,1"
	returnedDate, _, err := Split(header)

	if err == nil {
		t.Errorf("Successfully parsed event date with invalid year field: recieved date of '%s'", returnedDate)
	}
}
