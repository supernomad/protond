// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package common

import (
	"encoding/json"
	"time"
)

// Event represents an arbitrary event passing through protond.
type Event struct {
	Timestamp time.Time              `json:"timestamp"`
	Input     string                 `jsone:"input"`
	Data      map[string]interface{} `json:"data"`
}

// Bytes will return the byte slice representation of the event struct, optionally "pretty" printed, if there is an error during the marshalling process the returned byte slice will be nil.
func (e *Event) Bytes(pretty bool) []byte {
	var data []byte
	if pretty {
		data, _ = json.MarshalIndent(e, "", "    ")
	} else {
		data, _ = json.Marshal(e)
	}
	return data
}

// String will return the string representation of the event struct, optionally "pretty" printed, if there is an error during the marshalling process the returned string will be empty.
func (e *Event) String(pretty bool) string {
	return string(e.Bytes(pretty))
}

// ParseEventData will convert the supplied string to an Event struct pointer.
func ParseEventData(str string) (map[string]interface{}, error) {
	var eventData map[string]interface{}
	err := json.Unmarshal([]byte(str), &eventData)
	return eventData, err
}
