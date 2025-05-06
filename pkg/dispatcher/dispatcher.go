package dispatcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/ce7er2s/sunny_5_skiers/pkg/competitor"
	"github.com/ce7er2s/sunny_5_skiers/pkg/event"
)

func Dispatch(EventSource io.Reader, cfg Config) {
	fmt.Println("DEBUG: config info")
	fmt.Println(cfg)

	bufEventSource := bufio.NewReader(EventSource)
	var competitorsMap map[int]int = make(map[int]int)
	var competitors []competitor.Competitor

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

		// событие регистрации участника
		if evt.EventID == event.EVENT_ID_COMPETITOR_REGISTERED {
			competitors = append(competitors, competitor.NewCompetitor(evt.CompetitorID, cfg.StartTime, cfg.StartTime.Add(cfg.StartDelta)))
			competitorsMap[evt.CompetitorID] = len(competitors) - 1
			cfg.StartTime = cfg.StartTime.Add(cfg.StartDelta)
		} else {
			// eww stinky
			//HandlerMap[evt.EventID](&(competitors[competitorsMap[evt.CompetitorID]]), &evt)
		}

		// ^ переместить при реализации всей HandlerMap
		if evt.EventID == event.EVENT_ID_COMPETITOR_REGISTERED ||
			evt.EventID == event.EVENT_ID_START_TIME_SET_BY_DRAW ||
			evt.EventID == event.EVENT_ID_COMPETITOR_ON_START_LINE ||
			evt.EventID == event.EVENT_ID_COMPETITOR_STARTED {
			err = HandlerMap[evt.EventID](&(competitors[competitorsMap[evt.CompetitorID]]), &evt)
			if err != nil {
				log.Printf("Skip line \"%s\" because of error: %s", line, err.Error())
				continue
			}
		}

		evt_json, err := json.Marshal(evt)
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't convert Event to JSON: %s", err.Error()))
		}
		fmt.Println(string(evt_json))
	}

	fmt.Println("DEBUG: competitors info")
	for k, v := range competitors {
		fmt.Println(k, v)
	}
}
