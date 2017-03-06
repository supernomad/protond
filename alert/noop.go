// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package alert

import (
	"github.com/Supernomad/protond/common"
)

// Noop is a struct representing the standard input plugin.
type Noop struct {
	config *common.Config
	name   string
}

// Emit will return an empty list of events.
func (noop *Noop) Emit(event *common.Event) {
	return
}

// Name returns 'Noop'.
func (noop *Noop) Name() string {
	return noop.name
}

func newNoop(config *common.Config) (Alert, error) {
	noop := &Noop{
		config: config,
		name:   "Noop",
	}

	return noop, nil
}
