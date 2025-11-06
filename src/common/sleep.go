package common

import (
	"math/rand"
	"time"
)

const MIN_SLEEP_BETWEEN_RETRIES = 5
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES = 30
const MIN_SLEEP_BETWEEN_RETRIES_SHORT = 2
const MIN_SLEEP_BETWEEN_RETRIES_VERY_SHORT = 20
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SHORT = 10
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_VERY_SHORT = 10
const MIN_SLEEP_ON_START = 2
const MAX_OFFSET_SLEEP_ON_START = 30
const MIN_SLEEP_SHORT = 5
const MAX_OFFSET_SLEEP_SHORT = 10
const MIN_SLEEP_BETWEEN_SH_CONTACTS = 30
const MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS = 20
const MIN_SLEEP_BETWEEN_TASKS = 0
const MAX_OFFSET_SLEEP_BETWEEN_TASKS = 3

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

func SleepBetweenTasks() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_TASKS + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_TASKS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

func SleepBetweenRetriesVeryShort() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_VERY_SHORT + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_VERY_SHORT)
	t := time.Duration(r) * time.Millisecond
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
	r := MIN_SLEEP_SHORT + randSource.Intn(MAX_OFFSET_SLEEP_SHORT)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

func SleepBetweenShareContactsShort() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_SH_CONTACTS + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

func SleepBetweenShareContactsLarge() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := 3*MIN_SLEEP_BETWEEN_SH_CONTACTS + randSource.Intn(4*MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}
