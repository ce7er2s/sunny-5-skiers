package dispatcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/ce7er2s/sunny_5_skiers/pkg/event"
)

func Dispatch(EventSource io.Reader) {
	var events []event.Event
	bufEventSource := bufio.NewReader(EventSource)
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

		event, err := event.NewEvent(line)
		if err != nil {
			log.Printf("Skip line \"%s\" because of error: %s", line, err.Error())
			continue
		}

		events = append(events, event)
	}

	for _, event := range events {
		event_json, err := json.Marshal(event)
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't convert Event to JSON: %s", err.Error()))
		}
		fmt.Println(string(event_json))

	}
}
