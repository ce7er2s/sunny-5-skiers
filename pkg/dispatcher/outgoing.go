package dispatcher

type OutgoingIDType int32

const (
	OUTGOING_NOT_STARTED OutgoingIDType = iota
	OUTGOING_NOT_FINISHED
)

var OutgoinfIDToString = map[OutgoingIDType]string{
	OUTGOING_NOT_STARTED:  "OUTGOING_NOT_STARTED",
	OUTGOING_NOT_FINISHED: "OUTGOING_NOT_FINISHED",
}

type OutgoingEvent struct {
	OutgoingID OutgoingIDType
}

func (o OutgoingEvent) Error() string {
	return OutgoinfIDToString[o.OutgoingID]
}

func NewOutgoingEvent(reason OutgoingIDType) OutgoingEvent {
	return OutgoingEvent{
		OutgoingID: reason,
	}
}
