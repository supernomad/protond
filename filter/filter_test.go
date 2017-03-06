// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"testing"
	"time"

	"github.com/Supernomad/protond/alert"
	"github.com/Supernomad/protond/cache"
	"github.com/Supernomad/protond/common"
)

var (
	config        *common.Config
	badCfg        *common.Config
	internalCache cache.Cache
	alerts        map[string]alert.Alert
)

func init() {
	filterTimeout, _ := time.ParseDuration("10s")
	badTimeout, _ := time.ParseDuration("1ns")
	log := common.NewLogger(common.NoopLogger)
	config = &common.Config{FilterTimeout: filterTimeout, Log: log}
	badCfg = &common.Config{FilterTimeout: badTimeout, Log: log}
	internalCache, _ = cache.New(cache.MemoryCache, config, &common.PluginConfig{Name: "memory"})
	noopAlert, _ := alert.New(alert.NoopAlert, config, nil)
	alerts = map[string]alert.Alert{"Noop": noopAlert}
}

func TestNonExistentFilterPlugin(t *testing.T) {
	nonExistent, err := New("doesn't exist", nil, nil, nil, nil)
	if err == nil {
		t.Fatal("Something is very very wrong.")
	}
	if nonExistent != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestNoop(t *testing.T) {
	noop, err := New(NoopFilter, nil, nil, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}
	name := noop.Name()
	if name != "Noop" {
		t.Fatal("Something is very very wrong.")
	}

	test, err := noop.Run(event)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	if test != event {
		t.Fatal("Something is very very wrong.")
	}
}

func TestJavascript(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event.message = "testing"
			event.added_field = "woot"
			event.new_array = ["this", "should", "be", "handled", 1, 2, 3]
			event.new_object = {"woot": 123, "hello": "world", "sub_array":[1,2,3,"woot"]}
		`,
	}
	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	name := javascript.Name()
	if name != "Test Filter" {
		t.Fatal("javascript filter name is not properly handled")
	}

	test, err := javascript.Run(event)
	if err != nil {
		t.Fatalf("Error occurred: %s", err.Error())
	}

	if test.Data["message"] != "testing" {
		t.Fatalf("javascript filter failed to overwrite existing 'message' field")
	}

	if test.Data["added_field"] == nil || test.Data["added_field"] != "woot" {
		t.Fatalf("javascript filter failed to add a new field 'added_field' and/or set its value correctly")
	}

	if test.Data["new_array"] == nil || len(test.Data["new_array"].([]interface{})) != 7 {
		t.Fatalf("javascript filter failed to add a new field 'added_field' and/or set its value correctly")
	}
}

func TestJavascriptAlert(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event.message = "testing"
			event.added_field = "woot"
			event.new_array = ["this", "should", "be", "handled", 1, 2, 3]
			event.new_object = {"woot": 123, "hello": "world", "sub_array":[1,2,3,"woot"]}
			alert.emit("Doesn't exist", event)
			alert.emit({this: "should fail"}, event)
			alert.emit("Noop", "this should fail")
			alert.emit("Noop", event)
		`,
	}
	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	_, err = javascript.Run(event)
	if err != nil {
		t.Fatalf("Something is very very wrong. %s", err.Error())
	}
}

func TestJavascriptInternalCache(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event.message = "testing"
			event.added_field = "woot"
			event.new_array = ["this", "should", "be", "handled", 1, 2, 3]
			event.new_object = {"woot": 123, "hello": "world", "sub_array":[1,2,3,"woot"]}
			cache.store("testing", event)
		`,
	}

	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "woot",
		},
	}

	_, err = javascript.Run(event)
	if err != nil {
		t.Fatalf("Something is very very wrong. %s", err.Error())
	}

	test := javascript.(*Javascript).internalCache.Get("testing")
	if test == nil || len(test) != 1 || test[0].Data["message"] != "testing" {
		t.Fatal("internal cache was not properly updated")
	}
}

func TestJavascriptInternalCacheGet(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			events = cache.get("testing")
			if(events.length == 1) {
				event.stored_events = 1
			}
		`,
	}

	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "woot",
		},
	}

	test, err := javascript.Run(event)
	if err != nil {
		t.Fatalf("Something is very very wrong. %s", err.Error())
	}

	if test == nil || test.Data["stored_events"].(float64) != 1 {
		t.Fatal("event was not properly updated based on cached events")
	}
}

func TestJavascriptInternalCacheObjectKeyGet(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			events = cache.get({this:"should fail"})
			if(events != undefined) {
				event.failed = true
			} else {
				event.failed = false
			}
		`,
	}

	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "woot",
		},
	}

	test, err := javascript.Run(event)
	if err != nil {
		t.Fatalf("Something is very very wrong. %s", err.Error())
	}

	if test == nil || test.Data["failed"].(bool) != false {
		t.Fatal("event was not properly updated based on cached events")
	}
}

func TestJavascriptInternalCacheObjectKeyStore(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			cache.store({this:"should fail"}, event)
		`,
	}

	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "woot",
		},
	}

	_, err = javascript.Run(event)
	if err != nil {
		t.Fatalf("Something is very very wrong. %s", err.Error())
	}

	test := javascript.(*Javascript).internalCache.Get("testing")
	if test == nil || len(test) != 1 {
		t.Fatal("internal cache was not properly updated")
	}
}

func TestJavascriptInternalCacheStringValueStore(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			cache.store("testing", "this should fail")
		`,
	}

	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "woot",
		},
	}

	_, err = javascript.Run(event)
	if err != nil {
		t.Fatalf("Something is very very wrong. %s", err.Error())
	}

	test := javascript.(*Javascript).internalCache.Get("testing")
	if test == nil || len(test) != 1 {
		t.Fatal("internal cache was not properly updated")
	}
}

func TestJavascriptImproperTypeReturn(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event = "testing"
		`,
	}
	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	test, err := javascript.Run(event)
	if err == nil {
		t.Fatal("javascript filter improperly set event type but passed.")
	}
	if test == nil || test.Timestamp != event.Timestamp || test.Data["message"] != event.Data["message"] {
		t.Fatal("javascript filter improperly set event value on failure, should be the unchanged supplied event object.")
	}
}

func TestJavascriptImproperScript(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event = "testing"
			setTimeout(function() {
				console.log("this will never work")
			}, 100)
		`,
	}
	javascript, err := New(JavascriptFilter, config, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	test, err := javascript.Run(event)
	if err == nil {
		t.Fatal("javascript filter improperly handled a return line in the filter.")
	}
	if test == nil || test.Timestamp != event.Timestamp || test.Data["message"] != event.Data["message"] {
		t.Fatal("javascript filter improperly set event value on failure, should be the unchanged supplied event object.")
	}
}

func TestJavascriptInterrupt(t *testing.T) {
	filterConfig := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			function square(a) {
				return a * a
			}

			event.value = 1
			while(true) {
				event.value = square(event.value)
			}
		`,
	}
	javascript, err := New(JavascriptFilter, badCfg, filterConfig, internalCache, alerts)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	test, err := javascript.Run(event)
	if err == nil {
		t.Fatal("javascript filter improperly handled an interrupt.")
	}
	if test == nil || test.Timestamp != event.Timestamp || test.Data["message"] != event.Data["message"] {
		t.Fatal("javascript filter improperly set event value on failure, should be the unchanged supplied event object.")
	}
}
