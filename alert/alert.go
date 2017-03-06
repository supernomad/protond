// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package alert

import (
	"errors"

	"github.com/Supernomad/protond/common"
)

const (
	// NoopAlert defines a noop alert plugin used for testing.
	NoopAlert = "noop"
)

// Alert is the interface that plugins must adhere to for operation as an alert plugin.
type Alert interface {
	// Emit should send the event to the configured backend alert sink.
	Emit(event *common.Event)

	// Name returns the name of the alert plugin.
	Name() string
}

// New generates a alert plugin based on the passed in plugin and user defined configuration.
func New(alertPlugin string, config *common.Config, pluginConfig *common.PluginConfig) (Alert, error) {
	switch alertPlugin {
	case NoopAlert:
		return newNoop(config)
	}
	return nil, errors.New("specified alert plugin does not exist")
}
