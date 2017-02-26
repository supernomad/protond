// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"errors"

	"github.com/Supernomad/protond/common"
)

// Plugin defines the kind of input plugin to use.
type Plugin int

const (
	// StdinInput defines an input plugin that takes data from stdin.
	StdinInput Plugin = iota
)

// Input is the interface that plugins must adhere to for operation as an input plugin.
type Input interface {
	// Next should return the next event that is queued or received and a nil error object, if there is an error during the process the event should be nil and the error object should contain the error.
	Next() (*common.Event, error)

	// Name returns the name of the input plugin.
	Name() string

	// Close should completely terminate all functionality and destroy the input plugin.
	Close() error
}

// New generates an input plugin based on the passed in plugin and user defined configuration.
func New(inputPlugin Plugin, cfg *common.Config) (Input, error) {
	switch inputPlugin {
	case StdinInput:
		return newStdin(cfg)
	}
	return nil, errors.New("specified input plugin does not exist")
}
