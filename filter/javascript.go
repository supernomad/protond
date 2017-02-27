// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"errors"
	"time"

	"github.com/Supernomad/protond/common"
	"github.com/robertkrimen/otto"
)

var errHalt = errors.New("filter timed out")

// Javascript is a struct representing the javascript filter plugin.
type Javascript struct {
	cfg       *common.Config
	filterCfg *common.FilterConfig
}

// Run will return a parsed object based on the configured javascript filter.
func (js *Javascript) Run(event *common.Event) (ret *common.Event, err error) {
	defer func(js *Javascript, event *common.Event) {
		if caught := recover(); caught != nil {
			ret = event
			err = errors.New("filter '" + js.filterCfg.Name + "' paniced with '" + caught.(error).Error() + "' while parsing event: " + event.String(false))
		}
		return
	}(js, event)

	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)

	go func(js *Javascript, vm *otto.Otto) {
		time.Sleep(js.cfg.FilterTimeout)
		vm.Interrupt <- func() {
			panic(errHalt)
		}
	}(js, vm)

	vm.Set("event", event.Data)

	value, err := vm.Run(js.filterCfg.Code + "\nJSON.stringify(event)")
	if err != nil {
		return event, err
	}

	exported, _ := value.ToString()

	data, err := common.ParseEventData(exported)
	if err != nil {
		return event, errors.New("event data is no longer an object after running javascript filter, ensure that 'event' is always an object within the filter '" + js.filterCfg.Name + "', the returned value was: " + exported)
	}

	event.Data = data
	return event, nil
}

// Name returns configured name for the javascript filter.
func (js *Javascript) Name() string {
	return js.filterCfg.Name
}

func newJavascript(cfg *common.Config, filterCfg *common.FilterConfig) (Filter, error) {
	js := &Javascript{
		cfg:       cfg,
		filterCfg: filterCfg,
	}

	return js, nil
}
