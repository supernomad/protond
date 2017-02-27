// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"bufio"
	"net"
	"time"

	"github.com/Supernomad/protond/common"
)

// TCP is a struct representing the standard input plugin.
type TCP struct {
	cfg         *common.Config
	inOutConfig *common.InOutConfig
	messages    chan string
	listener    *net.TCPListener
}

func (tcp *TCP) accept() {
	for {
		conn, err := tcp.listener.AcceptTCP()
		if err != nil {
			tcp.cfg.Log.Error.Println("[TCP]", "Error accepting new connections with the tcp plugin.")
			break
		}

		go tcp.handleConn(conn)
	}
}

func (tcp *TCP) handleConn(conn *net.TCPConn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			tcp.cfg.Log.Error.Println("[TCP]", "Error reading from connection with the tcp plugin.")
			break
		}

		tcp.messages <- message
	}
}

// Next will return the next event from standard input.
func (tcp *TCP) Next() (*common.Event, error) {
	text := <-tcp.messages

	event := &common.Event{
		Timestamp: time.Now(),
		Input:     tcp.inOutConfig.Name,
		Data: map[string]interface{}{
			"message": text[:len(text)-1],
		},
	}

	return event, nil
}

// Name returns 'TCP'.
func (tcp *TCP) Name() string {
	return tcp.inOutConfig.Name
}

// Open will open the TCP plugin.
func (tcp *TCP) Open() error {
	addr, err := net.ResolveTCPAddr("tcp", tcp.inOutConfig.Config["host"]+":"+tcp.inOutConfig.Config["port"])
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	tcp.listener = l

	return nil
}

// Close will close the TCP plugin.
func (tcp *TCP) Close() error {
	err := tcp.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func newTCP(cfg *common.Config, inOutConfig *common.InOutConfig) (Input, error) {
	tcp := &TCP{
		cfg:         cfg,
		inOutConfig: inOutConfig,
	}

	return tcp, nil
}