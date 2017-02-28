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
	cfg := &common.Config{Log: common.NewLogger(common.NoopLogger), Backlog: 1024}

	in, err := input.New(input.NoopInput, cfg, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	filt, err := filter.New(filter.NoopFilter, cfg, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	out, err := output.New(output.NoopOutput, cfg, nil)
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}

	worker := New(cfg, []input.Input{in}, []filter.Filter{filt}, []output.Output{out})

	worker.Start()

	time.Sleep(5 * time.Second)

	err = worker.Stop()
	if err != nil {
		t.Fatal("Something is very very wrong.")
	}
}
