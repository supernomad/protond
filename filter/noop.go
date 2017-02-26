// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"github.com/Supernomad/protond/common"
)

// Noop is a struct representing the standard input plugin.
type Noop struct {
	cfg  *common.Config
	name string
}

// Run will return the event unchanged.
func (noop *Noop) Run(event *common.Event) (*common.Event, error) {
	return event, nil
}

// Name returns 'Noop'.
func (noop *Noop) Name() string {
	return noop.name
}

func newNoop(cfg *common.Config) (Filter, error) {
	noop := &Noop{
		cfg:  cfg,
		name: "Noop",
	}

	return noop, nil
}
