// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
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
	noop, err := New(NoopInput, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	name := noop.Name()
	if name != "Noop" {
		t.Fatal("Something is very very wrong.")
	}

	test, err := noop.Next()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	if test == nil {
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
	file, _ := ioutil.TempFile(os.TempDir(), "stdin")
	defer os.Remove(file.Name())
	os.Setenv("_TESTING_PROTOND", file.Name())

	stdin, err := New(StdinInput, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	writer := bufio.NewWriter(file)
	writer.WriteString("test\n")
	writer.Flush()

	test, err := stdin.Next()
	if err != nil {
		t.Fatalf("Something is very very wrong: %s", err.Error())
	}

	if test == nil {
		t.Fatal("Something is very very wrong: event is nil.")
	}

	if test.Data["message"] != "test" {
		t.Fatal("Something is very very wrong: event was improperlly parsed.")
	}

	name := stdin.Name()
	if name != "Stdin" {
		t.Fatal("Something is very very wrong.")
	}

	err = stdin.Open()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	err = stdin.Close()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}
}

func TestTCP(t *testing.T) {
	tcp, err := New(TCPInput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"host": "localhost"}})
	if err == nil || tcp != nil {
		t.Fatal("tcp plugin did not throw an error when configured without a port definition.")
	}

	tcp, err = New(TCPInput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"host": "localhost", "port": "9090"}})
	if err != nil {
		t.Fatalf("tcp plugin threw an error for no reason: %s", err.Error())
	}
	tcp.Open()

	time.Sleep(1 * time.Second)

	conn, _ := net.Dial("tcp", "127.0.0.1:9090")
	writer := bufio.NewWriter(conn)
	writer.WriteString("test\n")
	writer.Flush()
	conn.Close()

	time.Sleep(1 * time.Second)

	test, err := tcp.Next()
	if err != nil {
		t.Fatalf("Calling next on tcp input plugin errored: %s", err.Error())
	}

	if test == nil {
		t.Fatal("Something is wrong tcp input didn't error but returned event is nil.")
	}

	if test.Data["message"] != "test" {
		t.Fatal("tcp plugin improperlly parsed event.")
	}

	name := tcp.Name()
	if name != "Testing TCP" {
		t.Fatal("Something is wrong name wasn't handled properly.")
	}

	err = tcp.Close()
	if err != nil {
		t.Fatal("Something is wrong close wasn't handled properly.")
	}
}

func TestHttp(t *testing.T) {
	h, err := New(HTTPInput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing Http", Type: "http", Config: map[string]string{"host": "localhost"}})
	if err == nil || h != nil {
		t.Fatal("http plugin did not throw an error when configured without a port definition.")
	}

	h, err = New(HTTPInput, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)}, &common.PluginConfig{Name: "Testing Http", Type: "http", Config: map[string]string{"host": "localhost", "port": "9093"}})
	if err != nil {
		t.Fatal("http plugin did not throw an error when configured without a port definition.")
	}
	h.Open()

	time.Sleep(1 * time.Second)
	postData := map[string]string{
		"message": "test",
	}

	buf, _ := json.Marshal(postData)
	resp, err := http.Post("http://localhost:9093", "application/json", bytes.NewBuffer(buf))
	if err != nil || resp == nil || resp.StatusCode != 200 {
		t.Fatal("Something is wrong sent data wasn't handled properly.")
	}

	test, err := h.Next()
	if err != nil || test == nil {
		t.Fatal("Something is wrong couldn't retrieve sent data.")
	}

	if test.Data["message"] != "test" {
		t.Fatal("Something is wrong http plugin improperlly parsed event.")
	}

	resp, err = http.Post("http://localhost:9093", "text/plain", bytes.NewBuffer([]byte("testing string handling")))
	if err != nil || resp == nil || resp.StatusCode == 200 {
		t.Fatal("Something is wrong sent data wasn't handled properly.")
	}

	name := h.Name()
	if name != "Testing Http" {
		t.Fatal("Something is wrong name wasn't handled properly.")
	}

	err = h.Close()
	if err != nil {
		t.Fatal("Something is wrong close wasn't handled properly.")
	}
}
