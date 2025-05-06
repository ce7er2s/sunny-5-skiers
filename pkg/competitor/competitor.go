package competitor

import (
	"fmt"
	"slices"
	"time"
)

type Competitor struct {
	Status         CompetitorStatusType
	CompetitorID   int
	StartTime      time.Time
	EndTime        time.Time
	ShootingScore  [][]int
	Timings        [][2]time.Time
	LapCount       int
	FiringRange    int
	FiringLines    int
	PenaltyLaps    int
	PenaltyStart   time.Time
	PenaltyPeriod  float64
	ShotsAvailable int
	ShotsTaken     int
}

func NewCompetitor(id int, startTime time.Time, endTime time.Time, laps int, lines int) Competitor {
	shootingScore := make([][]int, laps)
	for i := range shootingScore {
		shootingScore[i] = make([]int, lines)
	}

	return Competitor{
		Status:         STATUS_REGISTERED,
		CompetitorID:   id,
		StartTime:      startTime,
		EndTime:        endTime,
		ShootingScore:  shootingScore,
		Timings:        make([][2]time.Time, laps),
		LapCount:       laps,
		FiringRange:    0,
		FiringLines:    lines,
		PenaltyLaps:    0,
		PenaltyPeriod:  0.0,
		ShotsAvailable: 0,
		ShotsTaken:     0,
	}
}

func (c *Competitor) SetStatus(status CompetitorStatusType) error {
	if slices.Contains(competitorFSM[c.Status], status) {
		c.Status = status
		return nil
	} else {
		return fmt.Errorf("Can't change state from %s to %s", competitorStatusToString[c.Status], competitorStatusToString[status])
	}
}
