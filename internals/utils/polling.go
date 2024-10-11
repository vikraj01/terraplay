package utils

import "time"



func CalculateBackoff(retry int) time.Duration {
	baseRetryDelay := 2 * time.Second
	return time.Duration(retry) * baseRetryDelay
}