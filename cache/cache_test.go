// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package cache

import (
	"testing"
	"time"

	"github.com/Supernomad/protond/common"
)

func TestNonExistentInputPlugin(t *testing.T) {
	nonExistent, err := New("doesn't exist", nil, nil)
	if err == nil {
		t.Fatal("Something is very very wrong.")
	}
	if nonExistent != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestNoop(t *testing.T) {
	noop, err := New(NoopCache, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	test := noop.Get("test")
	if test == nil || len(test) != 0 {
		t.Fatal("Something is very very wrong.")
	}

	noop.Store("test", nil)

	name := noop.Name()
	if name != "Noop" {
		t.Fatal("Something is very very wrong.")
	}
}

func TestMemory(t *testing.T) {
	memory, err := New(MemoryCache, nil, &common.PluginConfig{Name: "Memory Test"})
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	memory.Store("test", event)

	test := memory.Get("test")
	if test == nil || len(test) != 1 || test[0] != event {
		t.Fatal("Something is wrong memory cache did not store a single event correctly.")
	}

	memory.Store("test", event)

	test = memory.Get("test")
	if test == nil || len(test) != 2 || test[0] != event || test[1] != event {
		t.Fatal("Something is very very wrong.")
	}

	name := memory.Name()
	if name != "Memory Test" {
		t.Fatal("Something is very very wrong.")
	}
}
