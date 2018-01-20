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

// Encounter represents all the details about a given encounter
type Encounter struct {
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
	1:  "5-player",
	2:  "5-player Heroic",
	3:  "10-player Raid",
	4:  "25-player Raid",
	5:  "10-player Heroic Raid",
	6:  "25-player Heroic Raid",
	7:  "LFR",
	8:  "Challenge Mode",
	9:  "40-player Raid",
	11: "Heroic Scenario",
	12: "Scenario",
	14: "Regular",
	15: "Heroic",
	16: "Mythic",
	17: "LFR",
}
