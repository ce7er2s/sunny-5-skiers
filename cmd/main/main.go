package main

import (
	"os"

	"github.com/ce7er2s/sunny_5_skiers/pkg/dispatcher"
)

func main() {
	r, err := os.Open("events")
	if err != nil {
		os.Exit(1)
	}
	defer r.Close()

	dispatcher.Dispatch(r)
}
