package shared

import "time"

func TimeNow() time.Time {
	return time.Now().UTC()
}
