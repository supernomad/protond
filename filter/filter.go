// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"errors"

	"github.com/Supernomad/protond/alert"
	"github.com/Supernomad/protond/cache"
	"github.com/Supernomad/protond/common"
)

const (
	// NoopFilter defines a filter that does nothing.
	NoopFilter = "noop"

	// JavascriptFilter defines a javascript based filter.
	JavascriptFilter = "js"
)

// Filter is the interface that plugins must adhere to for operation as a filter plugin.
type Filter interface {
	// Run should take in the supplied event and preform the filtering, and then return the filtered event and a nil error object, if there is an error during the process the returned event should be the unchanged supplied event and the error object should contain the error.
	Run(*common.Event) (*common.Event, error)

	// Name returns the name of the filter plugin.
	Name() string
}

// New generates a filter plugin based on the passed in plugin and user defined configuration.
func New(filterPlugin string, config *common.Config, filterConfig *common.FilterConfig, internalCache cache.Cache, alerts map[string]alert.Alert) (Filter, error) {
	switch filterPlugin {
	case NoopFilter:
		return newNoop(config)
	case JavascriptFilter:
		return newJavascript(config, filterConfig, internalCache, alerts)
	}
	return nil, errors.New("specified filter plugin does not exist")
}
