package dispatcher

import (
	"time"
)

type Config struct {
	Laps          int    `json:"laps"`        // Количество кругов для основной дистанции
	LapLen        int    `json:"lapLen"`      // Длина каждого основного круга в метрах
	PenaltyLen    int    `json:"penaltyLen"`  // Длина каждого штрафного круга в метрах
	FiringLines   int    `json:"firingLines"` // Количество огневых рубежей (в задании "per lap", но обычно это общее число или число рубежей до финиша)
	SrcStartTime  string `json:"start"`
	SrcStartDelta string `json:"startDelta"`
	StartTime     time.Time
	StartDelta    time.Duration
}
