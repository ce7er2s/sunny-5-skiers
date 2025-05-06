package dispatcher

import (
	"fmt"
	"time"

	"github.com/ce7er2s/sunny-5-skiers/pkg/competitor"
)

func (l Lap) String() string {
	var t string
	var s string

	if l.Time == 0 {
		t = ""
	} else {
		t = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Add(l.Time).Format(timeLayout)
	}

	if l.Time == 0 {
		s = ""
	} else {
		s = fmt.Sprintf("%.3f", float64(l.LapLen)/l.Time.Seconds())
	}

	return fmt.Sprintf("{%s,%s}", t, s)
}

type Lap struct {
	Time   time.Duration
	Speed  float64
	LapLen int
}

type Report struct {
	Time         time.Time
	Status       competitor.CompetitorStatusType
	CompetitorID int
	Laps         []Lap

	PenaltyTime  time.Duration
	PenaltySpeed float64
	PenaltyLaps  int
	ShotsTaken   int
	ShotsHit     int
}

func NewReport(c *competitor.Competitor, cfg Config) Report {
	var laps []Lap = make([]Lap, len(c.Timings))
	var elapsedTime time.Time = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)

	// ewww stinky
	for i, v := range c.Timings {
		laps[i].LapLen = cfg.LapLen
		i = len(c.Timings) - i - 1
		laps[i].Time = v[1].Sub(v[0])
		elapsedTime = elapsedTime.Add(laps[i].Time)
		laps[i].Speed = float64(cfg.LapLen) / laps[i].Time.Seconds()
	}

	var shotsTaken int = c.ShotsTaken
	var shotsHit int = 0

	for _, array := range c.ShootingScore {
		for _, v := range array {
			shotsHit += v
		}
	}

	var penaltyTime time.Duration = time.Duration(c.PenaltyPeriod)
	var penaltySpeed float64 = 0

	if penaltyTime.Seconds() != 0 {
		penaltySpeed = float64(c.PenaltyLaps*cfg.PenaltyLen) / penaltyTime.Seconds()
	}

	return Report{
		Time:         elapsedTime,
		CompetitorID: c.CompetitorID,
		Laps:         laps,
		Status:       c.Status,
		ShotsTaken:   shotsTaken,
		ShotsHit:     shotsHit,
		PenaltyTime:  penaltyTime,
		PenaltySpeed: penaltySpeed,
		PenaltyLaps:  c.PenaltyLaps,
	}
}
