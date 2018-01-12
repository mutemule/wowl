package combatLog

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// func TestKill(t *testing.T) {
// 	expectedKill := "Kill"
// 	kill := killOrWipe(true)

// 	if kill != expectedKill {
// 		t.Errorf("Expecting %s, but got %s instead.\n", expectedKill, kill)
// 	}
// }

func TestParsingEventDateStamp(t *testing.T) {
	dateString := "1/10 15:02:15.348"
	expectedDate := time.Date(time.Now().Year(), 1, 10, 15, 02, 15, 348000000, time.UTC)

	// dateLayout := "1/2 15:04:05.000 2006"
	// expectedDateStamp, _ := time.Parse(dateLayout, dateStampString+" "+strconv.Itoa(time.Now().Year()))

	header := dateString + "  COMBAT_LOG_VERSION,4,ADVANCED_LOG_ENABLED,1"
	returnedDate, _, _ := ParseEvent(header)

	if expectedDate != returnedDate {
		t.Errorf("Date stamp parsing failed: expected '%s', but got '%s'\n", expectedDate, returnedDate)
	}
}

func TestParsingOldDateStamp(t *testing.T) {
	currentDate := time.Now().UTC()
	expectedDate := time.Date(currentDate.Year()-1, 1, currentDate.Day()+7, 15, 02, 15, 348000000, time.UTC)

	dateString := fmt.Sprintf("%d/%d %02d:%02d:%02d.348",
		int(expectedDate.Month()),
		expectedDate.Day(),
		expectedDate.Hour(),
		expectedDate.Minute(),
		expectedDate.Second(),
	)

	header := dateString + "  COMBAT_LOG_VERSION,4,ADVANCED_LOG_ENABLED,1"
	returnedDate, _, _ := ParseEvent(header)

	if expectedDate != returnedDate {
		t.Errorf("Date stamp parsing failed: expected '%s', but got '%s'\n", expectedDate, returnedDate)
	}
}

func TestParsingBasicCombatEvent(t *testing.T) {
	event := "1/10 15:02:15.540  SPELL_DAMAGE,Player-3694-07121BB2,\"Kerliah-Lightbringer\",0x514,0x0,Creature-0-3137-1712-27282-127233-0001566D00,\"Flameweaver\",0xa48,0x0,194153,\"Lunar Strike\",0x40,Creature-0-3137-1712-27282-127233-0001566D00,0000000000000000,23252757,199327752,0,0,1,0,0,0,-3053.49,10534.20,111,248907,-1,64,0,0,0,nil,nil,nil"
	expectedNumberOfEvents := 34
	_, returnedEvents, _ := ParseEvent(event)

	if expectedNumberOfEvents != len(returnedEvents) {
		t.Errorf("Event parsing failed: expected %d fields, but got %d.\n", expectedNumberOfEvents, len(returnedEvents))
	}
}

func TestParsingEmbeddedQuoteInCombatEvent(t *testing.T) {
	event := "1/7 16:54:31.838  SPELL_CAST_SUCCESS,Player-47-0781E030,\"Islen-Eitrigg\",0x10512,0x0,Creature-0-3777-1116-32565-77310-00005293F9,\"Mad \\\"King\\\" Sporeon\",0xa48,0x0,204197,\"Purge the Wicked\",0x4,Player-47-0781E030,0000000000000000,41258,41258,0,966,0,51000,51000,1020,1723.37,-738.11,388"
	expectedNumberOfEvents := 25
	_, returnedEvents, _ := ParseEvent(event)

	if expectedNumberOfEvents != len(returnedEvents) {
		t.Errorf("Event parsing failed: expected %d fields, but got %d.\n", expectedNumberOfEvents, len(returnedEvents))
	}
}

func TestParsingValidHeader(t *testing.T) {
	expectedHeader := Info{
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

	returnedHeader, _ := ParseHeader(header)

	if expectedHeader.Version != returnedHeader.Version {
		t.Errorf("Header parsing failed: expected verion %d, but got version %d\n", expectedHeader.Version, returnedHeader.Version)
	}

	if expectedHeader.AdvancedLogging != returnedHeader.AdvancedLogging {
		t.Errorf("Header parsing failed: advanced logging should be %v, but got %v\n", expectedHeader.AdvancedLogging, returnedHeader.AdvancedLogging)
	}
}
