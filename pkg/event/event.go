package event

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var eventRegexpString string = `^\[(?P<TimeStamp>\d{2}:\d{2}:\d{2}\.\d{3})\]\s+(?P<EventID>\d+)\s+(?P<CompetitorID>\d+)(?:\s+(?P<ExtraParams>[ \w:\.]+)?)?$`
var eventRegexp *regexp.Regexp = regexp.MustCompile(eventRegexpString)
var timeLayout string = "15:04:05.000"

// убрать json в конце -- используется для отладки
type Event struct {
	Timestamp    time.Time   `json:"TimeStamp"`
	EventID      EventIDType `json:"EventID"`
	CompetitorID int         `json:"CompetitorID"`
	ExtraParams  string      `json:"ExtraParams"`

	SourceString string `json:"SourceString"`
	// исходная из входящих
}

func NewEvent(line string) (Event, error) {
	if !eventRegexp.MatchString(line) {
		return Event{}, fmt.Errorf("String doesn't match event format.")
	}

	match := eventRegexp.FindStringSubmatch(line)
	fields := make(map[string]string)
	for i, name := range eventRegexp.SubexpNames() {
		fields[name] = match[i]
	}

	timestamp, err := time.Parse(timeLayout, fields["TimeStamp"])
	if err != nil {
		return Event{}, fmt.Errorf("Can't parse timestamp from \"%s\": %s", fields["TimeStamp"], err.Error())
	}

	id, err := strconv.Atoi(fields["EventID"])
	eventId := EventIDType(id)

	if err != nil {
		return Event{}, fmt.Errorf("Can't parse EventID from \"%s\": %s", fields["EventID"], err.Error())
	}

	competitorId, err := strconv.Atoi(fields["CompetitorID"])
	if err != nil {
		return Event{}, fmt.Errorf("Can't parse CompetitorID from \"%s\": %s", fields["CompetitorID"], err.Error())
	}

	return Event{
		Timestamp:    timestamp,
		EventID:      eventId,
		CompetitorID: competitorId,
		ExtraParams:  fields["ExtraParams"],
		SourceString: line,
	}, nil
}

func (e Event) String() string {
	timeStr := e.Timestamp.Format(timeLayout)
	var description string

	switch e.EventID {
	case EVENT_ID_COMPETITOR_REGISTERED: // 1
		description = fmt.Sprintf("The competitor(%d) registered", e.CompetitorID)
	case EVENT_ID_START_TIME_SET_BY_DRAW: // 2
		description = fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", e.CompetitorID, e.ExtraParams)
	case EVENT_ID_COMPETITOR_ON_START_LINE: // 3
		description = fmt.Sprintf("The competitor(%d) is on the start line", e.CompetitorID)
	case EVENT_ID_COMPETITOR_STARTED: // 4
		description = fmt.Sprintf("The competitor(%d) has started", e.CompetitorID)
	case EVENT_ID_COMPETITOR_ON_FIRING_RANGE: // 5
		description = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", e.CompetitorID, e.ExtraParams)
	case EVENT_ID_TARGET_HIT: // 6
		description = fmt.Sprintf("The target(%s) has been hit by competitor(%d)", e.ExtraParams, e.CompetitorID)
	case EVENT_ID_COMPETITOR_LEFT_FIRING_RANGE: // 7
		description = fmt.Sprintf("The competitor(%d) left the firing range", e.CompetitorID)
	case EVENT_ID_COMPETITOR_ENTERED_PENALTY: // 8
		description = fmt.Sprintf("The competitor(%d) entered the penalty laps", e.CompetitorID)
	case EVENT_ID_COMPETITOR_LEFT_PENALTY: // 9
		description = fmt.Sprintf("The competitor(%d) left the penalty laps", e.CompetitorID)
	case EVENT_ID_COMPETITOR_ENDED_MAIN_LAP: // 10
		description = fmt.Sprintf("The competitor(%d) ended the main lap", e.CompetitorID)
	case EVENT_ID_COMPETITOR_CANNOT_CONTINUE: // 11
		description = fmt.Sprintf("The competitor(%d) can`t continue: %s", e.CompetitorID, e.ExtraParams)
	}

	return fmt.Sprintf("[%s] %s", timeStr, description)
}
