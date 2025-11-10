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
const MIN_SLEEP_BETWEEN_SH_CONTACTS = 40
const MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS = 30
const MIN_SLEEP_LARGE = 40
const MAX_OFFSET_SLEEP_LARGE = 20

// Sleep de alrededor de 40 segundo con un desvío de 20
func SleepLarge() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_LARGE + randSource.Intn(MAX_OFFSET_SLEEP_LARGE)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// @TODO ajustar esto
func SleepBetweenRetries() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// @TODO ajustar esto
func SleepBetweenRetriesShort() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_SHORT + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SHORT)
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
	r := 20*MIN_SLEEP_BETWEEN_SH_CONTACTS + randSource.Intn(4*MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}
