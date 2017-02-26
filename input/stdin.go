// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package input

import (
	"bufio"
	"os"
	"time"

	"github.com/Supernomad/protond/common"
)

// Stdin is a struct representing the standard input plugin.
type Stdin struct {
	cfg    *common.Config
	name   string
	reader *bufio.Reader
}

// Next will return the next event from standard input.
func (stdin *Stdin) Next() (*common.Event, error) {
	text, _ := stdin.reader.ReadString('\n')
	event := &common.Event{
		Timestamp: time.Now(),
		Input:     stdin.name,
		Data: map[string]interface{}{
			"message": text[:len(text)-1],
		},
	}

	return event, nil
}

// Name returns 'Stdin'.
func (stdin *Stdin) Name() string {
	return stdin.name
}

// Close will close the Stdin plugin.
func (stdin *Stdin) Close() error {
	return nil
}

func newStdin(cfg *common.Config) (Input, error) {
	stdin := &Stdin{
		cfg:    cfg,
		name:   "Stdin",
		reader: bufio.NewReader(os.Stdin),
	}

	return stdin, nil
}
