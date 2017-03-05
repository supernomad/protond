package output

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Supernomad/protond/common"
	"github.com/Supernomad/protond/input"
)

func TestNonExistentOutputPlugin(t *testing.T) {
	nonExistent, err := New("doesn't exist", nil, nil)
	if err == nil {
		t.Fatal("Something is very very wrong.")
	}
	if nonExistent != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestNoop(t *testing.T) {
	noop, err := New(NoopOutput, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	name := noop.Name()
	if name != "Noop" {
		t.Fatal("Something is very very wrong.")
	}

	err = noop.Send(nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	err = noop.Open()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	err = noop.Close()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestStdin(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "stdout")
	defer os.Remove(file.Name())
	os.Setenv("_TESTING_PROTOND", file.Name())

	stdout, err := New(StdoutOutput, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	err = stdout.Send(event)
	if err != nil {
		t.Fatalf("Something is very very wrong: %s", err.Error())
	}

	name := stdout.Name()
	if name != "Stdout" {
		t.Fatal("Something is very very wrong.")
	}

	err = stdout.Open()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	err = stdout.Close()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestTCP(t *testing.T) {
	tcp, err := New(TCPOutput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"host": "localhost"}})
	if err == nil || tcp != nil {
		t.Fatal("tcp plugin did not throw an error when configured without a port definition.")
	}

	tcp, err = New(TCPOutput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"port": "8080"}})
	if err == nil || tcp != nil {
		t.Fatal("tcp plugin did not throw an error when configured without a port definition.")
	}

	inputTCP, err := input.New(input.TCPInput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"host": "127.0.0.1", "port": "9091"}})
	if err != nil {
		t.Fatalf("setting up input tcp plugin threw an error for no reason: %s", err.Error())
	}
	inputTCP.Open()

	time.Sleep(1 * time.Second)

	tcp, err = New(TCPOutput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"host": "127.0.0.1", "port": "9091"}})
	if err != nil {
		t.Fatalf("tcp plugin threw an error for no reason: %s", err.Error())
	}
	tcp.Open()

	time.Sleep(1 * time.Second)

	event := &common.Event{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": 101010101,
		},
	}

	err = tcp.Send(event)
	if err != nil {
		t.Fatalf("Something is very very wrong: %s", err.Error())
	}

	name := tcp.Name()
	if name != "Testing TCP" {
		t.Fatal("Something is wrong name wasn't handled properly.")
	}

	err = inputTCP.Close()
	if err != nil {
		t.Fatal("Something is wrong close wasn't handled properly.")
	}

	err = tcp.Close()
	if err != nil {
		t.Fatal("Something is wrong close wasn't handled properly.")
	}
	time.Sleep(1 * time.Second)
}
