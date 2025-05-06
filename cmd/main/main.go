package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ce7er2s/sunny_5_skiers/pkg/dispatcher"
)

var layoutTimeOnly string = "15:04:05"

func ParseConfigFromJSON(jsonData []byte) (*dispatcher.Config, error) {
	var cfg dispatcher.Config
	err := json.Unmarshal(jsonData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("JSON parse error: %s", err.Error())
	}

	if cfg.SrcStartTime != "" {
		timestamp, err := time.Parse(layoutTimeOnly, cfg.SrcStartTime)

		if err != nil {
			return nil, fmt.Errorf("start time parse error: %s", err.Error())
		}

		cfg.StartTime = timestamp
	}

	if cfg.SrcStartDelta != "" {
		timestamp, err := time.Parse(layoutTimeOnly, cfg.SrcStartDelta)

		if err != nil {
			return nil, fmt.Errorf("start time parse error: %s", err.Error())
		}

		// ewww stinky
		cfg.StartDelta = timestamp.Sub(time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC))
	}

	return &cfg, nil
}

func ParseConfigFromFile(filePath string) (*dispatcher.Config, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Can't read file \"%s\": %s", filePath, err.Error())
	}
	return ParseConfigFromJSON(jsonData)
}

func main() {
	cfg, err := ParseConfigFromFile("config.json")

	if err != nil {
		log.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	r, err := os.Open("events")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		os.Exit(2)
	}
	defer r.Close()

	dispatcher.Dispatch(r, *cfg)
}
