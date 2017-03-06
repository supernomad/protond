// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"bufio"
	"errors"
	"net"
	"time"

	"github.com/Supernomad/protond/common"
)

// TCP is a struct representing the standard input plugin.
type TCP struct {
	config       *common.Config
	pluginConfig *common.PluginConfig
	conn         *net.TCPConn
	writer       *bufio.Writer
}

func (tcp *TCP) handleConn(addr *net.TCPAddr) {
handle:
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		tcp.config.Log.Error.Printf("[TCP] New tcp connection to %s:%s could not be established: %s", tcp.pluginConfig.Config["host"], tcp.pluginConfig.Config["port"], err.Error())

		time.Sleep(10 * time.Second)
		goto handle
	}

	tcp.config.Log.Debug.Printf("[TCP] New tcp connection to %s:%s established.", tcp.pluginConfig.Config["host"], tcp.pluginConfig.Config["port"])

	tcp.conn = conn
	tcp.writer = bufio.NewWriter(conn)

	one := make([]byte, 1)
	for {
		if _, err := tcp.conn.Read(one); err != nil {
			tcp.conn.Close()

			tcp.conn = nil
			tcp.writer = nil

			tcp.config.Log.Debug.Printf("[TCP] tcp connection to %s:%s terminated: %s", tcp.pluginConfig.Config["host"], tcp.pluginConfig.Config["port"], err.Error())
			time.Sleep(10 * time.Second)
			goto handle
		}
	}
}

// Send will push the supplied event to the connected tcp server.
func (tcp *TCP) Send(event *common.Event) (err error) {
	defer func() {
		// Handle an interrupt to the javascript vm running the filter.
		if caught := recover(); caught != nil {
			err = errors.New("tcp connection to remote server is down")
		}
		return
	}()

	str := event.String(false)

	n, err := tcp.writer.WriteString(str + "\n")
	if err != nil {
		return err
	}
	if len(str)+1 != n {
		return errors.New("failed writing the entire event to the remote tcp server")
	}

	return tcp.writer.Flush()
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

	go tcp.handleConn(addr)

	return nil
}

// Close will close the TCP plugin.
func (tcp *TCP) Close() error {
	if tcp.conn == nil {
		return nil
	}

	err := tcp.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func newTCP(config *common.Config, pluginConfig *common.PluginConfig) (Output, error) {
	tcp := &TCP{
		config:       config,
		pluginConfig: pluginConfig,
	}

	if tcp.pluginConfig.Config["host"] == "" {
		return nil, errors.New("configuration for the tcp output plugin is missing a host definition")
	}

	if tcp.pluginConfig.Config["port"] == "" {
		return nil, errors.New("configuration for the tcp output plugin is missing a port definition")
	}

	return tcp, nil
}
