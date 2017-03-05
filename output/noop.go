// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"github.com/Supernomad/protond/common"
)

// Noop is a struct representing the noop plugin.
type Noop struct {
	config *common.Config
	name   string
}

// Send will noop.
func (noop *Noop) Send(event *common.Event) error {
	return nil
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

func newNoop(config *common.Config) (Output, error) {
	noop := &Noop{
		config: config,
		name:   "Noop",
	}
	return noop, nil
}
