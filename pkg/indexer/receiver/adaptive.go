package receiver

import "time"

const (
	ewmaAlpha     = 0.3
	thresholdHigh = 200.0
	thresholdLow  = 100.0
)

func (r *Module) adjustBulkSize(n int, elapsed time.Duration) {
	if n == 0 || r.cfg.RequestBulkSize == 1 {
		return
	}
	msPerBlock := float64(elapsed.Milliseconds()) / float64(n)

	r.ewmaMu.Lock()
	r.ewma = ewmaAlpha*msPerBlock + (1-ewmaAlpha)*r.ewma
	ema := r.ewma
	r.ewmaMu.Unlock()

	current := r.bulkSize.Load()

	r.Log.Debug().
		Float64("ema_ms", ema).
		Int64("bulk_size", current).
		Msg("current ema")
	switch {
	case ema > thresholdHigh && current > 1:
		if current-r.stepBulkSize > 1 {
			r.bulkSize.Store(current - r.stepBulkSize)
		} else {
			r.bulkSize.Store(1)
		}
		r.Log.Info().
			Float64("ema_ms", ema).
			Int64("bulk_size", r.bulkSize.Load()).
			Msg("decreased bulk_size")
	case ema < thresholdLow && current < r.maxBulkSize:
		if current+r.stepBulkSize < r.maxBulkSize {
			r.bulkSize.Store(current + r.stepBulkSize)
		} else {
			r.bulkSize.Store(r.maxBulkSize)
		}
		r.Log.Info().
			Float64("ema_ms", ema).
			Int64("bulk_size", r.bulkSize.Load()).
			Msg("increased bulk_size")
	}
}

func getStepBulkSize(value int64) int64 {
	if value > 20 {
		return value / 10
	}
	return 1
}
