package receiver

import "time"

const (
	ewmaAlpha     = 0.3
	thresholdHigh = 200.0
	thresholdLow  = 100.0
)

func (r *Module) adjustBulkSize(n int, elapsed time.Duration) {
	if n == 0 {
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
		r.bulkSize.Store(current - 1)
		r.Log.Info().
			Float64("ema_ms", ema).
			Int64("bulk_size", current-1).
			Msg("decreased bulk_size")
	case ema < thresholdLow && current < r.maxBulkSize:
		r.bulkSize.Store(current + 1)
		r.Log.Info().
			Float64("ema_ms", ema).
			Int64("bulk_size", current+1).
			Msg("increased bulk_size")
	}
}
