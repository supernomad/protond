// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"errors"

	"github.com/Supernomad/protond/common"
)

// Plugin defines the kind of filter plugin to use.
type Plugin int

const (
	// NoopFilter defines a filter that does nothing.
	NoopFilter Plugin = iota
)

// Filter is the interface that plugins must adhere to for operation as a filter plugin.
type Filter interface {
	// Run should take in the supplied event and preform the filtering, and then return the filtered event and a nil error object, if there is an error during the process the returned event should be the unchanged supplied event and the error object should contain the error.
	Run(*common.Event) (*common.Event, error)

	// Name returns the name of the filter plugin.
	Name() string

	// Close should completely terminate all functionality and destroy the filter plugin.
	Close() error
}

// New generates a filter plugin based on the passed in plugin and user defined configuration.
func New(filterPlugin Plugin, cfg *common.Config) (Filter, error) {
	switch filterPlugin {
	case NoopFilter:
		return newNoop(cfg)
	}
	return nil, errors.New("specified filter plugin does not exist")
}
