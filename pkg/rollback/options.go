package rollback

import "github.com/rs/zerolog"

// ModuleOption -
type ModuleOption func(module *Module)

// WithLogger -
func WithLogger(l zerolog.Logger) ModuleOption {
	return func(module *Module) {
		module.log = l.With().Str("module", module.Name()).Logger()
	}
}

// WithIndexerName -
func WithIndexerName(name string) ModuleOption {
	return func(module *Module) {
		module.indexName = name
	}
}
