// Package event provides Tesla Sentry event metadata handling.
package event

import (
	"encoding/json"
	"os"
	"path"
)

// Event represents Tesla Sentry Mode event metadata from event.json.
type Event struct {
	Timestamp string `json:"timestamp"`
	City      string `json:"city"`
	EstLat    string `json:"est_lat"`
	EstLon    string `json:"est_lon"`
	Reason    string `json:"reason"`
	Camera    string `json:"camera"`
}

// ParseEvent reads and parses event.json from the given directory.
func ParseEvent(directory string) (*Event, error) {
	f, err := os.Open(path.Join(directory, "event.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var event Event
	err = dec.Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}
