// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package main

import (
	"os"

	"github.com/Supernomad/protond/common"
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
	handleError(log, err)

	signaler := common.NewSignaler(log, cfg, nil, map[string]string{})

	log.Info.Println("[MAIN] protond start up complete.")

	err = signaler.Wait(true)
	handleError(log, err)
}
