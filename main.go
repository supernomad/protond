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
	log := common.NewLogger(common.InfoLogger)

	cfg, err := common.NewConfig(log)
	handleError(cfg.Log, err)

	workers := make([]*worker.Worker, cfg.NumWorkers)

	stdin, _ := input.New(input.StdinInput, cfg)
	noop, _ := filter.New(filter.NoopFilter, cfg)
	stdout, _ := output.New(output.StdoutOutput, cfg)

	for i := 0; i < cfg.NumWorkers; i++ {
		workers[i] = worker.New(cfg, []input.Input{stdin}, []filter.Filter{noop}, []output.Output{stdout})
		workers[i].Start()
	}

	signaler := common.NewSignaler(log, cfg, nil, map[string]string{})

	log.Info.Println("[MAIN] protond start up complete.")

	err = signaler.Wait(true)
	handleError(log, err)

	for i := 0; i < cfg.NumWorkers; i++ {
		workers[i].Stop()
	}
}
