package dispatcher

import (
	"time"
)

type Config struct {
	Laps          int    `json:"laps"`
	LapLen        int    `json:"lapLen"`
	PenaltyLen    int    `json:"penaltyLen"`
	FiringLines   int    `json:"firingLines"`
	SrcStartTime  string `json:"start"`
	SrcStartDelta string `json:"startDelta"`
	StartTime     time.Time
	StartDelta    time.Duration
}
