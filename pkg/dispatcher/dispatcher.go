package dispatcher

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/ce7er2s/sunny-5-skiers/pkg/competitor"
	"github.com/ce7er2s/sunny-5-skiers/pkg/event"
)

var timeLayout string = "15:04:05.000"

func Dispatch(EventSource io.Reader, cfg Config) {
	//fmt.Println("DEBUG: config info")
	//fmt.Println(cfg)

	bufEventSource := bufio.NewReader(EventSource)
	var competitorsMap map[int]int = make(map[int]int)
	var competitors []competitor.Competitor

	var competitorsWatch map[int]int = make(map[int]int)

	for {
		line, err := bufEventSource.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(fmt.Errorf("Fatal error while parsing a file: %s", err.Error()))
		}

		if len(line) == 0 {
			continue
		}

		line = line[:len(line)-1]
		evt, err := event.NewEvent(line)

		if err != nil {
			log.Printf("Skip line \"%s\" because of error: %s", line, err.Error())
			continue
		}

		// проверка не стартовавших участников
		for k, v := range competitorsWatch {
			if evt.Timestamp.After(competitors[v].EndTime) {
				competitors[v].SetStatus(competitor.STATUS_NOT_STARTED)
				fmt.Printf("%s -- Competitor(%d) disqualified.\n", competitors[v].EndTime.String(), competitors[v].CompetitorID)
				delete(competitorsWatch, k)
			}
		}

		// событие регистрации участника
		if evt.EventID == event.EVENT_ID_COMPETITOR_REGISTERED {
			competitors = append(competitors, competitor.NewCompetitor(evt.CompetitorID, cfg.StartTime, cfg.StartTime.Add(cfg.StartDelta), cfg.Laps, cfg.FiringLines))
			competitorsMap[evt.CompetitorID] = len(competitors) - 1
			cfg.StartTime = cfg.StartTime.Add(cfg.StartDelta)

			competitorsWatch[evt.CompetitorID] = competitorsMap[evt.CompetitorID]

		} else {
			err = HandlerMap[evt.EventID](&(competitors[competitorsMap[evt.CompetitorID]]), &evt)
			if err != nil {
				oevent, ok := err.(OutgoingEvent)
				if ok {
					if oevent.OutgoingID == OUTGOING_FINISHED {
						fmt.Printf("[%s] 33 %d\n", evt.Timestamp.Format(timeLayout), evt.CompetitorID)
					}
					if oevent.OutgoingID == OUTGOING_NOT_STARTED {
						fmt.Printf("[%s] 32 %d\n", competitors[competitorsMap[evt.CompetitorID]].EndTime.Format(timeLayout), evt.CompetitorID)
					}
				} else {
					log.Printf("Line \"%s\" error: %s", line, err.Error())
				}
			}
		}

		// запоминаем участников, которые не начинали
		if evt.EventID == event.EVENT_ID_START_TIME_SET_BY_DRAW || evt.EventID == event.EVENT_ID_COMPETITOR_ON_START_LINE || evt.EventID == event.EVENT_ID_COMPETITOR_REGISTERED {
			competitorsWatch[evt.CompetitorID] = competitorsMap[evt.CompetitorID]
		}

		// сносим, если они начали
		if evt.EventID == event.EVENT_ID_COMPETITOR_STARTED {
			delete(competitorsWatch, evt.CompetitorID)
		}

		// evt_json, err := json.Marshal(evt)
		// if err != nil {
		//	log.Fatal(fmt.Sprintf("Can't convert Event to JSON: %s", err.Error()))
		// }
		// fmt.Println(string(evt_json))
	}

	var reports []Report
	for _, v := range competitors {
		reports = append(reports, NewReport(&v, cfg))
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Time.Before(reports[j].Time)
	})

	for _, v := range reports {
		var status string

		if v.Status == competitor.STATUS_FINISHED {
			status = v.Time.Format(timeLayout)
		} else {
			status = competitor.CompetitorStatusToReportStatus[v.Status]
		}
		fmt.Printf("[%s] %d ", status, v.CompetitorID)

		// ewww stinky
		fmt.Print("[")
		for i, lap := range v.Laps {
			fmt.Print(lap.String())
			if i != len(v.Laps)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print("] ")

		// ewww stinky
		fmt.Print(Lap{v.PenaltyTime, v.PenaltySpeed, cfg.PenaltyLen * v.PenaltyLapsCount}.String())
		fmt.Printf(" %d/%d\n", v.ShotsHit, v.ShotsTaken)

	}
}
