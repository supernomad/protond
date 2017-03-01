// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package cache

import (
	"testing"
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
}
