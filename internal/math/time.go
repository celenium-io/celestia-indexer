package math

import "time"

func TimeFromNano(ts uint64) time.Time {
	return time.Unix(0, int64(ts)).UTC()
}
