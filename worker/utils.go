package worker

import (
	"time"

	"github.com/fsn-dev/crossChain-Bridge/log"
)

var (
	maxRecallLifetime       = int64(10 * 24 * 3600)
	restIntervalInRecallJob = 60 * time.Second

	maxVerifyLifetime       = int64(7 * 24 * 3600)
	restIntervalInVerifyJob = 60 * time.Second

	maxDoSwapLifetime       = int64(7 * 24 * 3600)
	restIntervalInDoSwapJob = 60 * time.Second

	maxStableLifetime       = int64(7 * 24 * 3600)
	restIntervalInStableJob = 60 * time.Second

	maxScanLifetime        = int64(3 * 24 * 3600)
	restIntervalInScanJob  = 10 * time.Second
	retryIntervalInScanJob = 10 * time.Second
)

func now() int64 {
	return time.Now().Unix()
}

func logWorker(job, subject string, context ...interface{}) {
	log.Info("["+job+"] "+subject, context...)
}

func logWorkerError(job, subject string, err error, context ...interface{}) {
	fields := []interface{}{"err", err}
	fields = append(fields, context...)
	log.Error("["+job+"] "+subject, fields...)
}

func getSepTimeInFind(dist int64) int64 {
	return now() - dist
}

func restInJob(duration time.Duration) {
	time.Sleep(duration)
}
