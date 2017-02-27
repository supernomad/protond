// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package main

import (
	"os"

	"github.com/Supernomad/protond/common"
	"github.com/Supernomad/protond/filter"
	"github.com/Supernomad/protond/input"
	"github.com/Supernomad/protond/output"
	"github.com/Supernomad/protond/worker"
)

func handleError(log *common.Logger, err error) {
	if err != nil {
		log.Error.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	logger := common.InfoLogger
	if os.Getenv("PROTOND_DEBUG") != "" {
		logger = common.DebugLogger
	}

	log := common.NewLogger(logger)

	cfg, err := common.NewConfig(log)
	handleError(cfg.Log, err)

	workers := make([]*worker.Worker, cfg.NumWorkers)

	filters := make([]filter.Filter, 0)
	for i := 0; i < len(cfg.Filters); i++ {
		temp, err := filter.New(cfg.Filters[i].Type, cfg, cfg.Filters[i])
		handleError(cfg.Log, err)

		filters = append(filters, temp)
	}

	if len(filters) == 0 {
		noop, _ := filter.New(filter.NoopFilter, cfg, nil)
		filters = append(filters, noop)
	}

	inputs := make([]input.Input, 0)
	for i := 0; i < len(cfg.Inputs); i++ {
		temp, err := input.New(cfg.Inputs[i].Type, cfg, cfg.Inputs[i])
		handleError(cfg.Log, err)

		err = temp.Open()
		handleError(cfg.Log, err)

		inputs = append(inputs, temp)
	}

	if len(inputs) == 0 {
		stdin, _ := input.New(input.StdinInput, cfg, nil)
		inputs = append(inputs, stdin)
	}

	outputs := make([]output.Output, 0)
	for i := 0; i < len(cfg.Outputs); i++ {
		temp, err := output.New(cfg.Outputs[i].Type, cfg, cfg.Outputs[i])
		handleError(cfg.Log, err)

		err = temp.Open()
		handleError(cfg.Log, err)

		outputs = append(outputs, temp)
	}

	if len(outputs) == 0 {
		stdout, _ := output.New(output.StdoutOutput, cfg, nil)
		outputs = append(outputs, stdout)
	}

	for i := 0; i < cfg.NumWorkers; i++ {
		workers[i] = worker.New(cfg, inputs, filters, outputs)
		workers[i].Start()
	}

	signaler := common.NewSignaler(log, cfg, nil, map[string]string{})

	log.Info.Println("[MAIN]", "protond start up complete.")

	err = signaler.Wait(true)
	handleError(cfg.Log, err)

	for i := 0; i < cfg.NumWorkers; i++ {
		workers[i].Stop()
	}
}
