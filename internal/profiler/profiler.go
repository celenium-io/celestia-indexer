// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package profiler

import (
	"fmt"
	"runtime"

	"github.com/grafana/pyroscope-go"
)

type Config struct {
	Server  string `validate:"omitempty,http_url" yaml:"server"`
	Project string `validate:"omitempty"          yaml:"project"`
}

func New(cfg *Config, service string) (*pyroscope.Profiler, error) {
	if cfg == nil || cfg.Server == "" {
		return nil, nil
	}

	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	return pyroscope.Start(pyroscope.Config{
		ApplicationName: fmt.Sprintf("%s-%s", cfg.Project, service),
		ServerAddress:   cfg.Server,
		Tags: map[string]string{
			"project": cfg.Project,
			"service": service,
		},

		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}
