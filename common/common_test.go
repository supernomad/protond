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

func TestFileExists(t *testing.T) {
	if FileExists("/this/path/should/never/exist") {
		t.Fatal("FileExists returned true for a non-existant file")
	}
	if !FileExists("common.go") {
		t.Fatal("FileExists returned false for a file that should always exist")
	}
}

func TestNewConfig(t *testing.T) {
	os.Setenv("PROTOND_CONF_FILE", confFile)
	os.Setenv("PROTOND_PID_FILE", "../protond.pid")

	os.Args = append(os.Args, "-w", "100")
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

	cfg.usage(false)
	cfg.version(false)
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
