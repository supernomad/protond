// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package alert

import (
	"testing"
	"time"

	"github.com/Supernomad/protond/common"
)

func TestNonExistentAlertPlugin(t *testing.T) {
	nonExistent, err := New("doesn't exist", nil, nil)
	if err == nil {
		t.Fatal("Something is very very wrong.")
	}
	if nonExistent != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestNoop(t *testing.T) {
	noop, err := New(NoopAlert, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	noop.Emit(event)

	name := noop.Name()
	if name != "Noop" {
		t.Fatal("Something is very very wrong.")
	}
}
