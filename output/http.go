// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/Supernomad/protond/common"
)

// HTTP is a struct representing the http output plugin.
type HTTP struct {
	config       *common.Config
	pluginConfig *common.PluginConfig
	uri          string
}

// Send takes the passed in event and sends it to the remote server.
func (h *HTTP) Send(event *common.Event) error {
	buf := event.Bytes(false)
	resp, err := http.Post(h.uri, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return err
	}

	if resp == nil || resp.StatusCode != 200 {
		return errors.New("error contacting remote server: " + err.Error())
	}

	return nil
}

// Name returns the name of the http output plugin.
func (h *HTTP) Name() string {
	return h.pluginConfig.Name
}

// Open starts the internal http(s) server, which will sending events to the remote server.
func (h *HTTP) Open() error {
	h.uri = fmt.Sprintf("%s://%s:%s%s", h.pluginConfig.Config["scheme"], h.pluginConfig.Config["host"], h.pluginConfig.Config["port"], h.pluginConfig.Config["route"])
	return nil
}

// Close terminates the internal http(s) server and frees all resources associated with the plugin.
func (h *HTTP) Close() error {
	return nil
}

func newHTTP(config *common.Config, pluginConfig *common.PluginConfig) (Output, error) {
	h := &HTTP{
		config:       config,
		pluginConfig: pluginConfig,
	}

	if h.pluginConfig.Config["scheme"] == "" {
		h.config.Log.Warn.Println("[OUTPUT]", "[HTTP]", "No scheme definition for the http output plugin, '"+h.pluginConfig.Name+"', using default scheme 'http'.")
		h.pluginConfig.Config["scheme"] = "http"
	}

	if h.pluginConfig.Config["host"] == "" {
		return nil, errors.New("configuration for the http output plugin, '" + h.pluginConfig.Name + "', is missing a host definition")
	}

	if h.pluginConfig.Config["port"] == "" {
		return nil, errors.New("configuration for the http output plugin, '" + h.pluginConfig.Name + "', is missing a port definition")
	}

	if h.pluginConfig.Config["route"] == "" {
		h.config.Log.Warn.Println("[OUTPUT]", "[HTTP]", "No route definition for the http output plugin, '"+h.pluginConfig.Name+"', using default route '/'.")
		h.pluginConfig.Config["route"] = "/"
	}

	return h, nil
}
