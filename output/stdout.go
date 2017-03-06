// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package output

import (
	"bufio"
	"os"

	"github.com/Supernomad/protond/common"
)

// Stdout is a struct representing the standard output plugin.
type Stdout struct {
	config *common.Config
	name   string
	writer *bufio.Writer
}

// Send writes the supplied event to standard output.
func (stdout *Stdout) Send(event *common.Event) error {
	str := event.String(true)

	n, err := stdout.writer.WriteString(str + "\n")
	if err != nil || len(str)+1 != n {
		return err
	}

	return stdout.writer.Flush()
}

// Name returns 'Stdout'.
func (stdout *Stdout) Name() string {
	return stdout.name
}

// Open will open the Stdout plugin.
func (stdout *Stdout) Open() error {
	return nil
}

// Close will close the Stdout plugin.
func (stdout *Stdout) Close() error {
	return nil
}

func newStdout(config *common.Config) (Output, error) {
	stdout := &Stdout{
		config: config,
		name:   "Stdout",
	}

	if tmpFile := os.Getenv("_TESTING_PROTOND"); tmpFile != "" {
		file, _ := os.OpenFile(tmpFile, os.O_APPEND|os.O_RDWR, os.ModeAppend)

		stdout.writer = bufio.NewWriter(file)
	} else {
		stdout.writer = bufio.NewWriter(os.Stdout)
	}

	return stdout, nil
}
