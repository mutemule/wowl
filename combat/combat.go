package combat

import (
	"time"
)

// Info represents the metadata about the combat log
type Info struct {
	Time            time.Time `json:"time"`
	Version         int       `json:"version"`
	AdvancedLogging bool      `json:"advancedlogging"`
	Header          string    `json:"header"`
}

// Fight is a generic representation for all battles, encounters, challenges, etc.
type Fight struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Start        time.Time   `json:"start"`
	End          time.Time   `json:"end"`
	DifficultyID int         `json:"difficultyID"`
	Difficulty   string      `json:"difficulty"`
	RaidSize     int         `json:"raidSize"`
	Kill         bool        `json:"kill"`
	Deaths       []UnitDeath `json:"deaths"`
	Events       []string    `json:"events"`
}

// UnitDeath records which units died and when
type UnitDeath struct {
	Time time.Time `json:"time"`
	Name string    `json:"name"`
}

// Difficulty maps the numeric encounter difficulty to plain english
var Difficulty = map[int]string{
	0:  "None",
	1:  "Regular Dungeon",
	2:  "Heroic Dungeon",
	3:  "10-player Raid",
	4:  "25-player Raid",
	5:  "10-player Heroic Raid",
	6:  "25-player Heroic Raid",
	7:  "Legacy LFR",
	8:  "Challenge Mode",
	9:  "40-player Raid",
	10: "Unknown",
	11: "Heroic Scenario",
	12: "Regular Scenario",
	13: "Unknown",
	14: "Regular Raid",
	15: "Heroic Raid",
	16: "Mythic Raid",
	17: "LFR",
	18: "Event",
	19: "Event",
	20: "Event Scenario",
	21: "Unknown",
	22: "Unknown",
	23: "Mythic Dungeon",
	24: "Timewalking",
	25: "PvP Scenario",
	26: "Unknown",
	27: "Unknown",
	28: "Unknown",
	29: "PvEvP Scenario",
	30: "Event",
	31: "Unknown",
	32: "PvP Scenario",
	33: "Timewalking",
	34: "PvP",
}
