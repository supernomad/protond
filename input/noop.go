// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"time"

	"github.com/Supernomad/protond/common"
)

// Noop is a struct representing the noop plugin.
type Noop struct {
	cfg  *common.Config
	name string
}

// Next will return the noop event.
func (noop *Noop) Next() (*common.Event, error) {
	event := &common.Event{
		Timestamp: time.Now(),
		Input:     noop.name,
		Data: map[string]interface{}{
			"message": "noop message",
		},
	}
	return event, nil
}

// Name returns 'Noop'.
func (noop *Noop) Name() string {
	return noop.name
}

// Open will open the Noop plugin.
func (noop *Noop) Open() error {
	return nil
}

// Close will close the Noop plugin.
func (noop *Noop) Close() error {
	return nil
}

func newNoop(cfg *common.Config) (Input, error) {
	noop := &Noop{
		cfg:  cfg,
		name: "Noop",
	}
	return noop, nil
}
