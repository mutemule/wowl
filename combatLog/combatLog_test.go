package combatLog

import (
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

func TestEventDateStampParsing(t *testing.T) {
	dateString := "1/10 15:02:15.348"
	expectedDate := time.Date(
		time.Now().Year(), 1, 10, 15, 02, 15, 348000000, time.UTC)

	// dateLayout := "1/2 15:04:05.000 2006"
	// expectedDateStamp, _ := time.Parse(dateLayout, dateStampString+" "+strconv.Itoa(time.Now().Year()))

	header := dateString + "  COMBAT_LOG_VERSION,4,ADVANCED_LOG_ENABLED,1"
	returnedDate, _, _ := ParseEvent(header)

	if expectedDate != returnedDate {
		t.Errorf("Date stamp parsing failed: expected '%s', but got '%s'\n", expectedDate, returnedDate)
	}
}
