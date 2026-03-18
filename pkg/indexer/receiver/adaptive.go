package receiver

import "time"

const (
	ewmaAlpha        = 0.3
	thresholdHigh    = 200.0
	thresholdLow     = 100.0
	decreaseCooldown = 30 * time.Second
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
		r.decreaseBulkSize(current)
		r.Log.Info().
			Float64("ema_ms", ema).
			Int64("bulk_size", r.bulkSize.Load()).
			Msg("decreased bulk_size")
	case ema < thresholdLow && current < r.maxBulkSize:
		// Guard against concurrent increases immediately after a decrease.
		// Without this, concurrent goroutines fetching light blocks undo the
		// decrease before the heavy batch can benefit from it.
		if time.Since(time.Unix(0, r.lastDecreasedAt.Load())) < decreaseCooldown {
			return
		}
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

// decreaseBulkSize reduces bulk_size by one step and records the decrease time.
func (r *Module) decreaseBulkSize(current int64) {
	next := current - r.stepBulkSize
	if next < 1 {
		next = 1
	}
	r.bulkSize.Store(next)
	r.lastDecreasedAt.Store(time.Now().UnixNano())
}

// connectionErrorDecrease halves bulk_size immediately and resets the EWMA to
// a high value. Called when the node closes the connection mid-stream
// (unexpected EOF, connection reset, etc.) — a clear signal that the batch was
// too large for the node's response timeout.
func (r *Module) connectionErrorDecrease() {
	current := r.bulkSize.Load()
	if current <= 1 {
		return
	}
	next := max(int64(1), current/2)
	r.bulkSize.Store(next)
	r.lastDecreasedAt.Store(time.Now().UnixNano())

	r.ewmaMu.Lock()
	r.ewma = thresholdHigh * 2 // keep EWMA above threshold so normal decrease fires first
	r.ewmaMu.Unlock()

	r.Log.Info().
		Int64("old_bulk_size", current).
		Int64("bulk_size", next).
		Msg("halved bulk_size due to connection error")
}

func getStepBulkSize(value int64) int64 {
	if value > 20 {
		return value / 10
	}
	return 1
}
