package competitor

type CompetitorStatusType int

const (
	STATUS_UNKNOWN CompetitorStatusType = iota
	STATUS_REGISTERED
	STATUS_ON_START_LINE
	STATUS_ON_MAIN_LAP
	STATUS_ON_FIRING_RANGE
	STATUS_ON_PENALTY_LAP
	STATUS_FINISHED
	STATUS_NOT_STARTED
	STATUS_NOT_FINISHED
)

var competitorStatusToString = map[CompetitorStatusType]string{
	STATUS_UNKNOWN:         "STATUS_UNKNOWN",
	STATUS_REGISTERED:      "STATUS_REGISTERED",
	STATUS_ON_START_LINE:   "STATUS_ON_START_LINE",
	STATUS_ON_MAIN_LAP:     "STATUS_ON_MAIN_LAP",
	STATUS_ON_FIRING_RANGE: "STATUS_ON_FIRING_RANGE",
	STATUS_ON_PENALTY_LAP:  "STATUS_ON_PENALTY_LAP",
	STATUS_FINISHED:        "STATUS_FINISHED",
	STATUS_NOT_STARTED:     "STATUS_NOT_STARTED",
	STATUS_NOT_FINISHED:    "STATUS_NOT_FINISHED",
}

var competitorFSM = map[CompetitorStatusType][]CompetitorStatusType{
	STATUS_UNKNOWN: {
		STATUS_REGISTERED,
	},

	STATUS_REGISTERED: {
		STATUS_ON_START_LINE,
		STATUS_NOT_STARTED,
	},

	STATUS_ON_START_LINE: {
		STATUS_ON_MAIN_LAP,
		STATUS_NOT_STARTED,
	},

	STATUS_ON_MAIN_LAP: {
		STATUS_ON_FIRING_RANGE,
		STATUS_NOT_FINISHED,
		STATUS_FINISHED,
	},

	STATUS_ON_FIRING_RANGE: {
		STATUS_ON_PENALTY_LAP,
		STATUS_ON_MAIN_LAP,
		STATUS_NOT_FINISHED,
	},

	STATUS_ON_PENALTY_LAP: {
		STATUS_ON_FIRING_RANGE,
		STATUS_ON_MAIN_LAP,
		STATUS_ON_PENALTY_LAP,
		STATUS_NOT_FINISHED,
	},

	STATUS_FINISHED:     {},
	STATUS_NOT_STARTED:  {},
	STATUS_NOT_FINISHED: {},
}

var CompetitorStatusToReportStatus = map[CompetitorStatusType]string{
	STATUS_NOT_FINISHED: "NotFinished",
	STATUS_NOT_STARTED:  "NotStarted",
}
