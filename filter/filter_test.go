// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package filter

import (
	"testing"
	"time"

	"github.com/Supernomad/protond/common"
)

func TestJavascript(t *testing.T) {
	config := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event.message = "testing"
			event.added_field = "woot"
			event.new_array = ["this", "should", "be", "handled", 1, 2, 3]
			event.new_object = {"woot": 123, "hello": "world", "sub_array":[1,2,3,"woot"]}
		`,
	}
	javascript, err := newJavascript(nil, config)

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
		t.Fatalf("Error occured: %s", err.Error())
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

func TestJavascriptImproperTypeReturn(t *testing.T) {
	config := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event = "testing"
		`,
	}
	javascript, err := newJavascript(nil, config)

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
	config := &common.FilterConfig{
		Name: "Test Filter",
		Code: `
			event = "testing"
			setTimeout(function() {
				console.log("this will never work")
			}, 100)
		`,
	}
	javascript, err := newJavascript(nil, config)

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
