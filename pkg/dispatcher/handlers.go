package dispatcher

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ce7er2s/sunny-5-skiers/pkg/competitor"
	"github.com/ce7er2s/sunny-5-skiers/pkg/event"
)

var timeLayout string = "15:04:05.000"

type EventHandleFunc func(c *competitor.Competitor, e *event.Event) error

// EventID = 1 (EVENT_ID_COMPETITOR_REGISTERED)
// фактически обрабатывается в dispatcher
func handleRegistration(_ *competitor.Competitor, _ *event.Event) error {
	return nil
}

// EventID = 2 (EVENT_ID_START_TIME_SET_BY_DRAW)
func handleSetStartTime(c *competitor.Competitor, e *event.Event) error {
	timestamp, err := time.Parse(timeLayout, e.ExtraParams)
	if err != nil {
		return err
	}

	duration := c.EndTime.Sub(c.StartTime)

	c.StartTime = timestamp
	c.EndTime = timestamp.Add(duration)
	// eww stinky
	return nil
}

// EventID = 3 (EVENT_ID_COMPETITOR_ON_START_LINE)
func handleOnStartLine(c *competitor.Competitor, _ *event.Event) error {
	return c.SetStatus(competitor.STATUS_ON_START_LINE)
}

// EventID = 4 (EVENT_ID_COMPETITOR_STARTED)
func handleStart(c *competitor.Competitor, e *event.Event) error {
	// если уже дисквалифицировали
	// if c.Status == competitor.STATUS_NOT_STARTED {
	//     return nil
	// }

	if e.Timestamp.Before(c.EndTime) {
		c.Timings[c.LapCount-1][0] = e.Timestamp
	} else {
		c.SetStatus(competitor.STATUS_NOT_STARTED)
		return NewOutgoingEvent(OUTGOING_NOT_STARTED)
	}

	return c.SetStatus(competitor.STATUS_ON_MAIN_LAP)
}

// EventID = 5 (EVENT_ID_COMPETITOR_ON_FIRING_RANGE)
func handleOnFiringRange(c *competitor.Competitor, e *event.Event) error {
	i, err := strconv.Atoi(e.ExtraParams)
	if err != nil {
		return err
	}

	c.Timings[c.LapCount-1][1] = e.Timestamp
	c.FiringRange = i - 1
	return c.SetStatus(competitor.STATUS_ON_FIRING_RANGE)
}

// EventID = 6 (EVENT_ID_TARGET_HIT)
func handleTargetHit(c *competitor.Competitor, e *event.Event) error {
	i, err := strconv.Atoi(e.ExtraParams)
	if err != nil {
		return err
	}
	fmt.Println(i)
	c.ShootingScore[c.FiringRange][i-1] = 1
	return nil
}

// EventID = 7 (EVENT_ID_COMPETITOR_LEFT_FIRING_RANGE)
func handleLeftFiringRange(c *competitor.Competitor, e *event.Event) error {
	var penaltyLaps int = 0
	for _, i := range c.ShootingScore[c.FiringRange] {
		penaltyLaps += i
	}

	c.PenaltyLaps = c.FiringLines - penaltyLaps

	if c.PenaltyLaps > 0 {
		return c.SetStatus(competitor.STATUS_ON_PENALTY_LAP)
	} else {
		return c.SetStatus(competitor.STATUS_ON_MAIN_LAP)
	}
}

// EventID = 8 (EVENT_ID_COMPETITOR_ENTERED_PENALTY)
func handleEnterPentalty(c *competitor.Competitor, e *event.Event) error {
	c.PenaltyStart = e.Timestamp

	return nil
}

// EventID = 9 (EVENT_ID_COMPETITOR_LEFT_PENALTY)
func handleLeftPentalty(c *competitor.Competitor, e *event.Event) error {
	c.PenaltyPeriod += float64(e.Timestamp.Sub(c.PenaltyStart))
	c.PenaltyLaps -= 1

	if c.PenaltyLaps == 0 {
		return c.SetStatus(competitor.STATUS_ON_MAIN_LAP)
	}

	return nil
}

func handleEndMainLap(c *competitor.Competitor, e *event.Event) error {
	c.LapCount -= 1

	if c.LapCount == 0 {
		return c.SetStatus(competitor.STATUS_FINISHED)
	}

	return nil
}

func handleCantContinue(c *competitor.Competitor, e *event.Event) error {
	err := c.SetStatus(competitor.STATUS_NOT_FINISHED)
	if err != nil {
		return err
	}

	return NewOutgoingEvent(OUTGOING_NOT_FINISHED)
}

var HandlerMap = map[event.EventIDType]EventHandleFunc{
	event.EVENT_ID_COMPETITOR_REGISTERED:        handleRegistration,
	event.EVENT_ID_START_TIME_SET_BY_DRAW:       handleSetStartTime,
	event.EVENT_ID_COMPETITOR_ON_START_LINE:     handleOnStartLine,
	event.EVENT_ID_COMPETITOR_STARTED:           handleStart,
	event.EVENT_ID_COMPETITOR_ON_FIRING_RANGE:   handleOnFiringRange,
	event.EVENT_ID_TARGET_HIT:                   handleTargetHit,
	event.EVENT_ID_COMPETITOR_LEFT_FIRING_RANGE: handleLeftFiringRange,
	event.EVENT_ID_COMPETITOR_ENTERED_PENALTY:   handleEnterPentalty,
	event.EVENT_ID_COMPETITOR_LEFT_PENALTY:      handleLeftPentalty,
	event.EVENT_ID_COMPETITOR_ENDED_MAIN_LAP:    handleEndMainLap,
	event.EVENT_ID_COMPETITOR_CANNOT_CONTINUE:   handleCantContinue,
}
