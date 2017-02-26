// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"errors"

	"github.com/Supernomad/protond/common"
)

// Type defines the kind of output plugin to use.
type Type int

const (
	// Stdout defines an output plugin that pushes data to stdout.
	Stdout Type = iota
)

// Output is the interface that plugins must adhere to for operation as an output plugin.
type Output interface {
	// Send should take the passed in event and send it to the arbitrary endpoint or data sink.
	Send(*common.Event) error

	// Close should completely terminate all functionality and destroy the output plugin.
	Close() error
}

// New generates an output plugin based on the passed in type and user defined configuration.
func New(outputType Type, cfg *common.Config) (Output, error) {
	switch outputType {
	case Stdout:

	}
	return nil, errors.New("specified output type does not exist")
}
