// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"errors"
	"strings"
	"time"

	"github.com/Supernomad/protond/alert"
	"github.com/Supernomad/protond/cache"
	"github.com/Supernomad/protond/common"
	"github.com/robertkrimen/otto"
)

var errHalt = errors.New("filter timed out")

const (
	alertInternal = `var alert = {
		emit: function(pluginName, evt, extra_params) {
			strigifiedParams = JSON.stringify(extra_params);
			strigifiedEvt = JSON.stringify(evt);

			_alert(pluginName, strigifiedEvt, strigifiedParams);
		}
	};`

	cacheInternal = `var cache = {
		get: function(key){
			return _get(key);
		},
		store: function(key, evt) {
			strigifiedEvt = JSON.stringify(evt);

			_store(key, strigifiedEvt);
		}
	};`

	returnInternal = `JSON.stringify(event);`
)

// Javascript is a struct representing the javascript filter plugin.
type Javascript struct {
	config        *common.Config
	filterConfig  *common.FilterConfig
	alerts        map[string]alert.Alert
	internalCache cache.Cache
}

func renderScript(filterCode string) string {
	return strings.Join([]string{
		alertInternal,
		cacheInternal,
		filterCode,
		returnInternal,
	}, "\n")
}

func interrupt(js *Javascript, vm *otto.Otto) {
	time.Sleep(js.config.FilterTimeout)
	vm.Interrupt <- func() {
		panic(errHalt)
	}
}

// Run will return a parsed object based on the configured javascript filter.
func (js *Javascript) Run(event *common.Event) (ret *common.Event, err error) {
	defer func() {
		// Handle an interrupt to the javascript vm running the filter.
		if caught := recover(); caught != nil {
			ret = event
			err = errors.New("filter '" + js.filterConfig.Name + "' paniced with '" + caught.(error).Error() + "' while parsing event: " + event.String(false))
		}
		return
	}()

	vm := otto.New()

	vm.Interrupt = make(chan func(), 1)
	go interrupt(js, vm)

	vm.Set("_alert", func(call otto.FunctionCall) otto.Value {
		js.config.Log.Debug.Printf("[FILTER] [JS] Filter, '%s', alert plugin function 'emit' called with plugin name, '%s'.", js.filterConfig.Name, call.Argument(0))

		plugin, err := call.Argument(0).ToString()
		if err != nil || strings.Contains(plugin, "Object") {
			js.config.Log.Error.Printf("[FILTER] [JS] Filter, '%s', errored with call to 'alert.emit', first argument was not a string.", js.filterConfig.Name)
			return otto.Value{}
		}

		if alert, ok := js.alerts[plugin]; ok {
			strEvent, _ := call.Argument(1).ToString()
			data, err := common.ParseEventData(strEvent)
			if err != nil {
				js.config.Log.Error.Printf("[FILTER] [JS] Filter, '%s', errored with call to 'cache.store', second argument was not an event object.", js.filterConfig.Name)
				return otto.Value{}
			}

			event.Data = data
			alert.Emit(event)
		}
		return otto.Value{}
	})

	vm.Set("_get", func(call otto.FunctionCall) otto.Value {
		js.config.Log.Debug.Printf("[FILTER] [JS] Filter, '%s', cache plugin function 'get' called with key, '%s'.", js.filterConfig.Name, call.Argument(0))

		key, err := call.Argument(0).ToString()
		if err != nil || strings.Contains(key, "Object") {
			js.config.Log.Error.Printf("[FILTER] [JS] Filter, '%s', errored with call to 'cache.get', first argument was not a string.", js.filterConfig.Name)
			return otto.Value{}
		}

		events := js.internalCache.Get(key)
		val, _ := vm.ToValue(events)
		return val
	})

	vm.Set("_store", func(call otto.FunctionCall) otto.Value {
		js.config.Log.Debug.Printf("[FILTER] [JS] Filter, '%s', cache plugin function 'store' called with key, '%s', and value, '%s'.", js.filterConfig.Name, call.Argument(0), call.Argument(1))

		key, err := call.Argument(0).ToString()
		if err != nil || strings.Contains(key, "Object") {
			js.config.Log.Error.Printf("[FILTER] [JS] Filter, '%s', errored with call to 'cache.store', first argument was not a string.", js.filterConfig.Name)
			return otto.Value{}
		}

		strEvent, _ := call.Argument(1).ToString()
		data, err := common.ParseEventData(strEvent)
		if err != nil {
			js.config.Log.Error.Printf("[FILTER] [JS] Filter, '%s', errored with call to 'cache.store', second argument was not an event object.", js.filterConfig.Name)
			return otto.Value{}
		}

		event.Data = data
		js.internalCache.Store(key, event)
		return otto.Value{}
	})

	vm.Set("event", event.Data)
	value, err := vm.Run(renderScript(js.filterConfig.Code))
	if err != nil {
		return event, err
	}

	exported, _ := value.ToString()

	data, err := common.ParseEventData(exported)
	if err != nil {
		return event, errors.New("event data is no longer an object after running javascript filter, ensure that 'event' is always an object within the filter '" + js.filterConfig.Name + "', the returned value was: " + exported)
	}

	event.Data = data
	return event, nil
}

// Name returns configured name for the javascript filter.
func (js *Javascript) Name() string {
	return js.filterConfig.Name
}

func newJavascript(config *common.Config, filterConfig *common.FilterConfig, internalCache cache.Cache, alerts map[string]alert.Alert) (Filter, error) {
	if alerts == nil {
		alerts = make(map[string]alert.Alert)
	}

	js := &Javascript{
		config:        config,
		filterConfig:  filterConfig,
		internalCache: internalCache,
		alerts:        alerts,
	}

	return js, nil
}
