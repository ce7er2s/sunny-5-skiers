package dispatcher

import (
	"fmt"
	"time"

	"github.com/ce7er2s/sunny_5_skiers/pkg/competitor"
	"github.com/ce7er2s/sunny_5_skiers/pkg/event"
)

var timeLayout string = "15:04:05.000"

type EventHandleFunc func(c *competitor.Competitor, e *event.Event) error

// проверка сходится ли id события и участника
// излишне?
func checkCompetitorID(handler EventHandleFunc) EventHandleFunc {
	return func(c *competitor.Competitor, e *event.Event) error {
		if c.CompetitorID != e.CompetitorID {
			return fmt.Errorf("Can't handle event: CompetitorID doesn't match with given competitor")
		}
		return handler(c, e)
	}
}

// EventID = 1 (EVENT_ID_COMPETITOR_REGISTERED)
// фактически обрабатывается в dispatcher
func handleRegistration(_ *competitor.Competitor, _ *event.Event) error {
	return nil
}

// EventID = 2 (EVENT_ID_START_TIME_SET_BY_DRAW)
func handleSetStartTime(c *competitor.Competitor, e *event.Event) error {
	timestamp, err := time.Parse(timeLayout, e.ExtraParams)
	if err != nil {
		return fmt.Errorf("Can't parse timestamp from extraParams: \"%s\": %s", e.ExtraParams, err.Error())
	}

	c.StartTime = timestamp
	return nil
}

// EventID = 3 (EVENT_ID_COMPETITOR_ON_START_LINE)
func handleOnStartLine(c *competitor.Competitor, _ *event.Event) error {
	return c.SetStatus(competitor.STATUS_ON_START_LINE)
}

// EventID = 4 (EVENT_ID_COMPETITOR_STARTED)
func handleStart(c *competitor.Competitor, e *event.Event) error {
	if e.Timestamp.Before(c.EndTime) {
		err := c.SetStatus(competitor.STATUS_ON_MAIN_LAP)
		if err != nil {
			return err
		}
		// реализовать таймеры для кругов
		//
	} else {
		return NewOutgoingEvent(OUTGOING_NOT_STARTED)
	}

	return nil
}

var HandlerMap = map[event.EventIDType]EventHandleFunc{
	event.EVENT_ID_COMPETITOR_REGISTERED:    handleRegistration,
	event.EVENT_ID_START_TIME_SET_BY_DRAW:   handleSetStartTime,
	event.EVENT_ID_COMPETITOR_ON_START_LINE: handleOnStartLine,
	event.EVENT_ID_COMPETITOR_STARTED:       handleStart,
}
