// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"bufio"
	"errors"
	"net"
	"time"

	"github.com/Supernomad/protond/common"
)

// TCP is a struct representing the tcp input plugin.
type TCP struct {
	cfg          *common.Config
	pluginConfig *common.PluginConfig
	messages     chan string
	listener     *net.TCPListener
}

func (tcp *TCP) accept() {
	for {
		conn, err := tcp.listener.AcceptTCP()
		if err != nil {
			tcp.cfg.Log.Error.Println("[TCP]", "Error accepting new connections with the tcp plugin.")
			break
		}

		tcp.cfg.Log.Debug.Println("[TCP]", "New tcp connection received.")
		go tcp.handleConn(conn)
	}
}

func (tcp *TCP) handleConn(conn *net.TCPConn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			tcp.cfg.Log.Debug.Println("[TCP]", "Error reading from connection with the tcp plugin, considering connection dead and moving on.")
			break
		}

		tcp.cfg.Log.Debug.Println("[TCP]", "New tcp message received.")
		tcp.messages <- message
	}
}

// Next will return the next event from the internal event buffer.
func (tcp *TCP) Next() (*common.Event, error) {
	text := <-tcp.messages

	event := &common.Event{
		Timestamp: time.Now(),
		Input:     tcp.pluginConfig.Name,
		Data: map[string]interface{}{
			"message": text[:len(text)-1],
		},
	}

	return event, nil
}

// Name returns 'TCP'.
func (tcp *TCP) Name() string {
	return tcp.pluginConfig.Name
}

// Open will open the TCP plugin.
func (tcp *TCP) Open() error {
	addr, err := net.ResolveTCPAddr("tcp", tcp.pluginConfig.Config["host"]+":"+tcp.pluginConfig.Config["port"])
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	tcp.cfg.Log.Debug.Printf("[TCP] New tcp listener created on %s:%s.", tcp.pluginConfig.Config["host"], tcp.pluginConfig.Config["port"])
	tcp.listener = l

	go tcp.accept()

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

func newTCP(cfg *common.Config, pluginConfig *common.PluginConfig) (Input, error) {
	tcp := &TCP{
		cfg:          cfg,
		pluginConfig: pluginConfig,
		messages:     make(chan string, cfg.Backlog),
	}

	if tcp.pluginConfig.Config["port"] == "" {
		return nil, errors.New("configuration for the tcp input plugin is missing a port definition")
	}

	return tcp, nil
}
