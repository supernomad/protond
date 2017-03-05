// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"errors"

	"github.com/Supernomad/protond/common"
)

const (
	// NoopInput defines a no operation input plugin for testing.
	NoopInput = "noop"

	// StdinInput defines an input plugin that takes data from stdin.
	StdinInput = "stdin"

	// TCPInput defines an input plugin that takes data from a tcp socket.
	TCPInput = "tcp"

	// HTTPInput defins an input plugin that taks json data being posted from http clients, which can run with or without TLS.
	HTTPInput = "http"
)

// Input is the interface that plugins must adhere to for operation as an input plugin.
type Input interface {
	// Next should return the next event that is queued or received and a nil error object, if there is an error during the process the event should be nil and the error object should contain the error.
	Next() (*common.Event, error)

	// Name returns the name of the input plugin.
	Name() string

	// Open should fully initialize the input plugin.
	Open() error

	// Close should completely terminate all functionality and destroy the input plugin.
	Close() error
}

// New generates an input plugin based on the passed in plugin and user defined configuration.
func New(inputPlugin string, cfg *common.Config, pluginConfig *common.PluginConfig) (Input, error) {
	switch inputPlugin {
	case NoopInput:
		return newNoop(cfg)
	case StdinInput:
		return newStdin(cfg)
	case TCPInput:
		return newTCP(cfg, pluginConfig)
	case HTTPInput:
		return newHTTP(cfg, pluginConfig)
	}
	return nil, errors.New("specified input plugin does not exist")
}
