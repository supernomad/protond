// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package cache

import (
	"github.com/Supernomad/protond/common"
)

// Memory is a struct representing the standard input plugin.
type Memory struct {
	config       *common.Config
	pluginConfig *common.PluginConfig
	events       map[string][]*common.Event
}

// Get will return an empty list of events.
func (memory *Memory) Get(key string) []*common.Event {
	return memory.events[key]
}

// Store will memory the store process of a cache plugin.
func (memory *Memory) Store(key string, event *common.Event) {
	if _, ok := memory.events[key]; ok {
		memory.events[key] = append(memory.events[key], event)
	} else {
		memory.events[key] = []*common.Event{event}
	}
}

// Name returns the name of the memory cache.
func (memory *Memory) Name() string {
	return memory.pluginConfig.Name
}

func newMemory(config *common.Config, pluginConfig *common.PluginConfig) (Cache, error) {
	memory := &Memory{
		config:       config,
		pluginConfig: pluginConfig,
		events:       make(map[string][]*common.Event),
	}

	return memory, nil
}
