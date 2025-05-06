package dispatcher

import (
	"fmt"

	"github.com/ce7er2s/sunny_5_skiers/pkg/competitor"
	"github.com/ce7er2s/sunny_5_skiers/pkg/event"
)

type EventHandleFunc func(c *competitor.Competitor, e *event.Event) error

// decorator :O
func checkCompetitorID(handler EventHandleFunc) EventHandleFunc {
	return func(c *competitor.Competitor, e *event.Event) error {
		if c.CompetitorID != e.CompetitorID {
			return fmt.Errorf("Can't handle event: CompetitorID doesn't match with given competitor")
		}
		return handler(c, e)
	}
}

// EventID = 1 (EVENT_ID_COMPETITOR_REGISTERED)
func handleRegistration(c *competitor.Competitor, _ *event.Event) error {
	c.Status = "Registered"
	return nil
}
