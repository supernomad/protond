// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package cache

import (
	"github.com/Supernomad/protond/common"
)

// Noop is a struct representing the standard input plugin.
type Noop struct {
	config *common.Config
	name   string
	events []*common.Event
}

// Get will return an empty list of events.
func (noop *Noop) Get(key string) []*common.Event {
	return noop.events
}

// Store will noop the store process of a cache plugin.
func (noop *Noop) Store(key string, event *common.Event) {
	return
}

// Name returns 'Noop'.
func (noop *Noop) Name() string {
	return noop.name
}

func newNoop(config *common.Config) (Cache, error) {
	noop := &Noop{
		config: config,
		name:   "Noop",
		events: make([]*common.Event, 0),
	}

	return noop, nil
}
