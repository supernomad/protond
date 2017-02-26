// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"bufio"
	"errors"
	"os"

	"github.com/Supernomad/protond/common"
)

// Stdout is a struct representing the standard output plugin.
type Stdout struct {
	cfg    *common.Config
	name   string
	writer *bufio.Writer
}

// Send writes the supplied event to standard output.
func (s *Stdout) Send(event *common.Event) error {
	str := event.String(true)

	n, err := s.writer.WriteString(str + "\n")
	if err != nil {
		return err
	}
	if len(str)+1 != n {
		return errors.New("failed writing the entire event to standard out")
	}

	return s.writer.Flush()
}

// Name returns 'Stdout'.
func (s *Stdout) Name() string {
	return s.name
}

// Close will close the Stdout plugin.
func (s *Stdout) Close() error {
	return nil
}

func newStdout(cfg *common.Config) (Output, error) {
	stdout := &Stdout{
		cfg:    cfg,
		name:   "Stdout",
		writer: bufio.NewWriter(os.Stdout),
	}

	return stdout, nil
}
