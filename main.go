// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package main

import (
	"os"

	"github.com/Supernomad/protond/cache"
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

	config, err := common.NewConfig(log)
	handleError(config.Log, err)

	internalCache, err := cache.New(cache.MemoryCache, config, &common.PluginConfig{Name: "memory"})
	handleError(config.Log, err)

	workers := make([]*worker.Worker, config.NumWorkers)

	filters := make([]filter.Filter, 0)
	for i := 0; i < len(config.Filters); i++ {
		temp, err := filter.New(config.Filters[i].Type, config, config.Filters[i], internalCache, nil)
		handleError(config.Log, err)

		filters = append(filters, temp)
	}

	if len(filters) == 0 {
		noop, _ := filter.New(filter.NoopFilter, config, nil, nil, nil)
		filters = append(filters, noop)
	}

	inputs := make([]input.Input, 0)
	for i := 0; i < len(config.Inputs); i++ {
		temp, err := input.New(config.Inputs[i].Type, config, config.Inputs[i])
		handleError(config.Log, err)

		err = temp.Open()
		handleError(config.Log, err)

		inputs = append(inputs, temp)
	}

	if len(inputs) == 0 {
		stdin, _ := input.New(input.StdinInput, config, nil)
		inputs = append(inputs, stdin)
	}

	outputs := make([]output.Output, 0)
	for i := 0; i < len(config.Outputs); i++ {
		temp, err := output.New(config.Outputs[i].Type, config, config.Outputs[i])
		handleError(config.Log, err)

		err = temp.Open()
		handleError(config.Log, err)

		outputs = append(outputs, temp)
	}

	if len(outputs) == 0 {
		stdout, _ := output.New(output.StdoutOutput, config, nil)
		outputs = append(outputs, stdout)
	}

	for i := 0; i < config.NumWorkers; i++ {
		workers[i] = worker.New(config, inputs, filters, outputs)
		workers[i].Start()
	}

	signaler := common.NewSignaler(log, config, nil, map[string]string{})

	log.Info.Println("[MAIN]", "protond start up complete.")

	err = signaler.Wait(true)
	handleError(config.Log, err)

	for i := 0; i < config.NumWorkers; i++ {
		workers[i].Stop()
	}
}
