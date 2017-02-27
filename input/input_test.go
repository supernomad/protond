// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"bufio"
	"io/ioutil"
	"net"
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
	tcp, _ := New(TCPInput, &common.InOutConfig{Name: "Testing TCP", Type: "tcp", Config: map[string]string{"host": "localhost", "port": "9090"}}, &common.Config{Backlog: 1024, Log: common.NewLogger(common.NoopLogger)})
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

	tcp.Close()
}
