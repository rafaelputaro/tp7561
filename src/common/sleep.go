package common

import (
	"math/rand"
	"time"
)

const MIN_SLEEP_BETWEEN_RETRIES = 5
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES = 30
const MIN_SLEEP_BETWEEN_RETRIES_SHORT = 2
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SHORT = 10
const MIN_SLEEP_ON_START = 2
const MAX_OFFSET_SLEEP_ON_START = 30
const MIN_SLEEP_SHORT = 1
const MAX_OFFSET_SLEEP_SHORT = 10

func SleepBetweenRetries() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

func SleepBetweenRetriesShort() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_SHORT + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SHORT)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

func SleepOnStart(numberOfParticipants int) {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_ON_START*numberOfParticipants/4 + randSource.Intn(MAX_OFFSET_SLEEP_ON_START)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

func SleepShort(numberOfParticipants int) {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_SHORT*numberOfParticipants/2 + randSource.Intn(MAX_OFFSET_SLEEP_SHORT)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}
