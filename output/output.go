// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"errors"

	"github.com/Supernomad/protond/common"
)

// Plugin defines the kind of output plugin to use.
type Plugin int

const (
	// StdoutOutput defines an output plugin that pushes data to stdout.
	StdoutOutput Plugin = iota
)

// Output is the interface that plugins must adhere to for operation as an output plugin.
type Output interface {
	// Send should take the passed in event and send it to the arbitrary endpoint or data sink.
	Send(*common.Event) error

	// Name returns the name of the output plugin.
	Name() string

	// Close should completely terminate all functionality and destroy the output plugin.
	Close() error
}

// New generates an output plugin based on the passed in plugin and user defined configuration.
func New(outputPlugin Plugin, cfg *common.Config) (Output, error) {
	switch outputPlugin {
	case StdoutOutput:
		return newStdout(cfg)
	}
	return nil, errors.New("specified output plugin does not exist")
}
