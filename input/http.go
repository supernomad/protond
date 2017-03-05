// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Supernomad/protond/common"
)

// HTTP is a struct representing the http input plugin.
type HTTP struct {
	cfg          *common.Config
	pluginConfig *common.PluginConfig
	messages     chan map[string]interface{}
}

type response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
}

func (h *HTTP) setHeaders(w http.ResponseWriter) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Server", "protond")
}

func (h *HTTP) handleResponseError(err error) {
	if err != nil {
		h.cfg.Log.Error.Println("[INPUT]", "[HTTP]", "Error sending response to client:", err.Error())
	}
}

func (h *HTTP) handleRequestError(w http.ResponseWriter, requestError error) {
	body := response{
		Message: "Error handling request, POSTed data must be a json blob.",
		Error:   requestError.Error(),
	}
	resp, _ := json.Marshal(body)

	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write(resp)

	h.handleResponseError(err)
}

func (h *HTTP) handleSuccess(w http.ResponseWriter) {
	body := response{
		Message: "event received",
	}
	resp, _ := json.Marshal(body)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(resp)

	h.handleResponseError(err)
}

func (h *HTTP) handleEvents(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		h.handleRequestError(w, err)
		return
	}

	h.messages <- data
	h.handleSuccess(w)
}

func (h *HTTP) server() {
	http.HandleFunc(h.pluginConfig.Config["route"], h.handleEvents)
	for {
		err := http.ListenAndServe(h.pluginConfig.Config["host"]+":"+h.pluginConfig.Config["port"], nil)
		if err != nil {
			h.cfg.Log.Error.Println("[INPUT]", "[HTTP]", "Error initializing event api:", err.Error())
		}

		time.Sleep(10 * time.Second)
	}
}

// Next will return the next event on the internal event buffer.
func (h *HTTP) Next() (*common.Event, error) {
	msg := <-h.messages

	event := &common.Event{
		Timestamp: time.Now(),
		Input:     h.pluginConfig.Name,
		Data:      msg,
	}

	return event, nil
}

// Name returns the name of the current http plugin.
func (h *HTTP) Name() string {
	return h.pluginConfig.Name
}

// Open starts the internal http(s) server, which will start queuing events on its internal event buffer.
func (h *HTTP) Open() error {
	go h.server()
	return nil
}

// Close terminates the internal http(s) server and frees all resources associated with the plugin.
func (h *HTTP) Close() error {
	return nil
}

func newHTTP(cfg *common.Config, pluginConfig *common.PluginConfig) (Input, error) {
	h := &HTTP{
		cfg:          cfg,
		pluginConfig: pluginConfig,
		messages:     make(chan map[string]interface{}, cfg.Backlog),
	}

	if h.pluginConfig.Config["port"] == "" {
		return nil, errors.New("configuration for the http input plugin, '" + h.pluginConfig.Name + "', is missing a port definition")
	}

	if h.pluginConfig.Config["route"] == "" {
		h.cfg.Log.Warn.Println("[INPUT]", "[HTTP]", "No route definition for the http input plugin, '"+h.pluginConfig.Name+"', using default route '/'.")
		h.pluginConfig.Config["route"] = "/"
	}

	return h, nil
}
