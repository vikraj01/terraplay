package utils

import (
	"log"
	"time"
)

func NewProgressLogger(totalRetries int, authTimeout time.Duration) *ProgressLogger {
	return &ProgressLogger{
		startTime:    time.Now(),
		totalRetries: totalRetries,
		authTimeout:  authTimeout,
	}
}

type ProgressLogger struct {
	startTime    time.Time
	totalRetries int
	authTimeout  time.Duration
}

func (pl *ProgressLogger) LogProgress(retries int, message string) {
	elapsed := time.Since(pl.startTime)
	progress := (float64(retries) / float64(pl.totalRetries)) * 100
	remainingTime := pl.authTimeout - elapsed

	log.Printf(
		"[%0.1f%%] %s | Retries: %d/%d | Elapsed: %s | Remaining: ~%s",
		progress,
		message,
		retries,
		pl.totalRetries,
		elapsed.Truncate(time.Second),
		remainingTime.Truncate(time.Second),
	)
}

func (pl *ProgressLogger) LogCompletion(success bool, message string) {
	elapsed := time.Since(pl.startTime)
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}

	log.Printf(
		"[100%%] %s | Status: %s | Total Time: %s",
		message,
		status,
		elapsed.Truncate(time.Second),
	)
}
