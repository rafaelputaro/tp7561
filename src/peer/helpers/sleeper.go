package helpers

import (
	"math/rand"
	"time"
)

const MIN_SLEEP_BETWEEN_RETRIES = 10
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES = 30

func SleepBetweenRetries() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}
