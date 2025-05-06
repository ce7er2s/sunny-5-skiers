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
