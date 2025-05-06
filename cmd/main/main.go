package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ce7er2s/sunny_5_skiers/pkg/event"
)

func main() {
	fmt.Println("Start.")
	evt, err := event.NewEvent("[09:55:00.000] 2 1 10:00:00.000")
	if err != nil {
		log.Fatal(err.Error())
	}

	evt_json, err := json.Marshal(*evt)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't convert Event to JSON: %s", err.Error()))
	}

	fmt.Println(string(evt_json))
}
