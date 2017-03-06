package worker

import (
	"testing"
	"time"

	"github.com/Supernomad/protond/common"
	"github.com/Supernomad/protond/filter"
	"github.com/Supernomad/protond/input"
	"github.com/Supernomad/protond/output"
)

func TestWorker(t *testing.T) {
	config := &common.Config{Log: common.NewLogger(common.NoopLogger), Backlog: 1024}

	in, err := input.New(input.NoopInput, config, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	filt, err := filter.New(filter.NoopFilter, config, nil, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	out, err := output.New(output.NoopOutput, config, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	worker := New(config, []input.Input{in}, []filter.Filter{filt}, []output.Output{out})

	worker.Start()

	time.Sleep(5 * time.Second)

	err = worker.Stop()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	time.Sleep(1 * time.Second)
}
