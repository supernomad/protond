// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package common

import (
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"
)

const (
	confFile = "../dist/test/protond.yml"
)

func TestPathExists(t *testing.T) {
	if PathExists("/this/path/should/never/exist") {
		t.Fatal("PathExists returned true for a non-existent file")
	}
	if !PathExists("common.go") {
		t.Fatal("PathExists returned false for a file that should always exist")
	}
}

func TestNewConfig(t *testing.T) {
	os.Setenv("PROTOND_CONF_FILE", confFile)
	os.Setenv("PROTOND_PID_FILE", "../protond.pid")

	os.Args = append(os.Args, "-w", "100", "-f", "../dist/test/filters.d")
	cfg, err := NewConfig(NewLogger(NoopLogger))

	if err != nil {
		t.Fatalf("NewConfig returned an error, %s", err)
	}
	if cfg == nil {
		t.Fatal("NewConfig returned a blank config")
	}
	if cfg.NumWorkers != runtime.NumCPU() {
		t.Fatal("NewConfig didn't pick up the cli replacement for NumWorkers")
	}
	if cfg.DataDir != "/var/lib/testing-protond" {
		t.Fatal("NewConfig didn't pick up the config file replacement for DataDir")
	}
	if cfg.PidFile != "../protond.pid" {
		t.Fatal("NewConfig didn't pick up the environment variable replacement for PidFile")
	}

	cfg.parseSpecial([]string{"-h", "-v"}, false)
}

func TestNewLogger(t *testing.T) {
	log := NewLogger(NoopLogger)
	if log.Error == nil {
		t.Fatal("NewLogger returned a nil Error log.")
	}
	if log.Warn == nil {
		t.Fatal("NewLogger returned a nil Warn log.")
	}
	if log.Info == nil {
		t.Fatal("NewLogger returned a nil Info log.")
	}
	if log.Debug == nil {
		t.Fatal("NewLogger returned a nil Debug log.")
	}
}

func TestSignaler(t *testing.T) {
	log := NewLogger(NoopLogger)
	cfg, err := NewConfig(log)
	signaler := NewSignaler(log, cfg, []int{1}, map[string]string{"PROTOND_TESTING": "woot"})

	go func() {
		time.Sleep(1 * time.Second)
		signaler.signals <- syscall.SIGHUP
		time.Sleep(1 * time.Second)
		signaler.signals <- syscall.SIGINT
	}()

	err = signaler.Wait(false)
	if err != nil {
		t.Fatal("Wait returned an error: " + err.Error())
	}
	err = signaler.Wait(false)
	if err != nil {
		t.Fatal("Wait returned an error: " + err.Error())
	}
}

func TestEvent(t *testing.T) {
	event := &Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"number": 101010101,
			"array": []interface{}{
				"woot",
				"awesome",
				10101,
				false,
				[]interface{}{"hello", "sub", "array"},
			},
			"map": map[string]interface{}{
				"first":  "value",
				"second": 42,
				"third":  []interface{}{"this", "is", "another", 2, "sub", "array"},
				"fourth": map[string]interface{}{"sub": "map", "working": 30303, "yes": true},
			},
		},
	}

	bytes := event.Bytes(false)
	if bytes == nil {
		t.Fatal("Event.Bytes(false) returned nil for a valid event.")
	}

	bytes = event.Bytes(true)
	if bytes == nil {
		t.Fatal("Event.Bytes(true) returned nil for a valid event.")
	}

	str := event.String(false)
	if str == "" {
		t.Fatal("Event.String(false) returned an empty string for a valid event.")
	}

	str = event.String(true)
	if str == "" {
		t.Fatal("Event.String(true) returned an empty string for a valid event.")
	}
}

func TestParseEventData(t *testing.T) {
	testStr := `{"woot": 234, "sub_obj":{"hello": "world", "array":[1,2,3,true]}, "sub_array":["woot", {"sub":"object"}]}`
	data, err := ParseEventData(testStr)
	if err != nil {
		t.Fatal("ParseEventData returned an error parsing arbitrary event data.")
	}
	if data["woot"].(float64) != 234 {
		t.Fatal("ParseEventData returned an incorrect value for arbitrary data.")
	}
}
