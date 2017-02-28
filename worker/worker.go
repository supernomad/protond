// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package worker

import (
	"github.com/Supernomad/protond/common"
	"github.com/Supernomad/protond/filter"
	"github.com/Supernomad/protond/input"
	"github.com/Supernomad/protond/output"
)

// Worker represents an individual protond worker.
type Worker struct {
	cfg *common.Config

	incoming chan *common.Event
	outgoing chan *common.Event

	stopReading   chan struct{}
	stopFiltering chan struct{}
	stopWriting   chan struct{}

	filters []filter.Filter
	inputs  []input.Input
	outputs []output.Output
}

func (w *Worker) input(input int) {
	for {
		select {
		case <-w.stopReading:
			close(w.incoming)
			return
		default:
			event, err := w.inputs[input].Next()
			if err != nil {
				w.cfg.Log.Error.Printf("errored getting next event from input '%s'\nerror: %s", w.inputs[input].Name(), err.Error())
			} else {
				w.incoming <- event
			}
		}
	}
}

func (w *Worker) filter() {
	for {
		select {
		case event := <-w.incoming:
			var err error

			for i := 0; i < len(w.filters); i++ {
				event, err = w.filters[i].Run(event)
				if err != nil {
					w.cfg.Log.Error.Printf("errored running filter '%s' on event: %s\nerror: %s", w.filters[i].Name(), event.String(false), err.Error())
					break
				}
			}

			if err == nil {
				w.outgoing <- event
			}
		case <-w.stopFiltering:
			close(w.stopFiltering)
			return
		}
	}
}

func (w *Worker) output() {
	for {
		select {
		case event := <-w.outgoing:
			for i := 0; i < len(w.outputs); i++ {
				err := w.outputs[i].Send(event)
				if err != nil {
					w.cfg.Log.Error.Printf("errored sending to output '%s' on event: %s\nerror: %s", w.outputs[i].Name(), event.String(false), err.Error())
				}
			}
		case <-w.stopWriting:
			close(w.stopWriting)
			close(w.outgoing)
			return
		}
	}
}

// Start the protond worker, so it will begin processing events.
func (w *Worker) Start() {
	for i := 0; i < len(w.inputs); i++ {
		go w.input(i)
	}
	go w.filter()
	go w.output()
}

// Stop the protond worker, terminating all processing of events.
func (w *Worker) Stop() error {
	for i := 0; i < len(w.inputs); i++ {
		w.stopReading <- struct{}{}
		err := w.inputs[i].Close()
		if err != nil {
			return err
		}
	}
	close(w.stopReading)

	w.stopFiltering <- struct{}{}
	w.stopWriting <- struct{}{}

	for i := 0; i < len(w.outputs); i++ {
		err := w.outputs[i].Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// New returns a worker object that is fully configured and ready to be started.
func New(cfg *common.Config, inputs []input.Input, filters []filter.Filter, outputs []output.Output) *Worker {
	return &Worker{
		cfg:           cfg,
		inputs:        inputs,
		filters:       filters,
		outputs:       outputs,
		incoming:      make(chan *common.Event, cfg.Backlog),
		outgoing:      make(chan *common.Event, cfg.Backlog),
		stopReading:   make(chan struct{}, len(inputs)),
		stopFiltering: make(chan struct{}, 1),
		stopWriting:   make(chan struct{}, 1),
	}
}
