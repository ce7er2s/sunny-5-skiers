package event;

import (
	"time"
)

type Event struct {
	Timestamp    time.Time
	ID           int
	CompetitorID int
	ExtraParams  string

	SourceString string			// исходная из входящих
}
