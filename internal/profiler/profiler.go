package profiler

import (
	"github.com/pyroscope-io/client/pyroscope"
)

type Config struct {
	Server          string `validate:"omitempty,http_url" yaml:"server"`
	Project         string `validate:"omitempty"          yaml:"project"`
	ApplicationName string `validate:"omitempty"          yaml:"app_name"`
}

func New(cfg *Config, service string) (*pyroscope.Profiler, error) {
	if cfg == nil || cfg.Server == "" {
		return nil, nil
	}

	return pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.ApplicationName,
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
