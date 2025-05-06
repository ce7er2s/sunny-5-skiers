package competitor

import (
	"fmt"
	"slices"
	"time"
)

type Competitor struct {
	Status       CompetitorStatusType
	CompetitorID int
	StartTime    time.Time
	EndTime      time.Time
}

func NewCompetitor(id int, startTime time.Time, endTime time.Time) Competitor {
	return Competitor{
		Status:       STATUS_REGISTERED,
		CompetitorID: id,
		StartTime:    startTime,
		EndTime:      endTime,
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
