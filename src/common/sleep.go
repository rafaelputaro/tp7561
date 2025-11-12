package common

import (
	"math/rand"
	"time"
)

const MIN_SLEEP_BETWEEN_RETRIES = 15
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES = 5

const MIN_SLEEP_BETWEEN_RETRIES_ADD_FILE = 5
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_ADD_FILE = 1

const MIN_SLEEP_BETWEEN_RETRIES_GET_FILE = 5
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_GET_FILE = 1

const MIN_SLEEP_BETWEEN_RETRIES_SH_CTS = 2
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SH_CTS = 1

const MIN_SLEEP_BETWEEN_RETRIES_SHORT = 2
const MIN_SLEEP_BETWEEN_RETRIES_VERY_SHORT = 20

const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SHORT = 10
const MAX_OFFSET_SLEEP_BETWEEN_RETRIES_VERY_SHORT = 10

const MIN_SLEEP_ON_START = 2
const MAX_OFFSET_SLEEP_ON_START = 30

const MIN_SLEEP_SHORT = 5
const MAX_OFFSET_SLEEP_SHORT = 10

const MIN_SLEEP_BETWEEN_SH_CONTACTS = 50 //40
const MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS = 10

const MIN_SLEEP_LARGE = 40
const MAX_OFFSET_SLEEP_LARGE = 20

// Sleep de alrededor de 40 segundo con un desvío de 20
func SleepLarge() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_LARGE + randSource.Intn(MAX_OFFSET_SLEEP_LARGE)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 5 segundos con desvío de 1
func SleepBetweenRetriesAddFile() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_ADD_FILE + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_ADD_FILE)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 5 segundos con desvío de 1
func SleepBetweenRetriesGetFile() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_GET_FILE + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_GET_FILE)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 2 segundos con desvío de 1
func SleepBetweenRetriesShareContactsRecip() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_SH_CTS + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SH_CTS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 15 segundos con un desvío de 5 segundos
func SleepBetweenRetries() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 2 segundos con un desvìo de 10 segundos
func SleepBetweenRetriesShort() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_RETRIES_SHORT + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_RETRIES_SHORT)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 20ms con un desvío de 10ms
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

// Sleep de 50 segundos con un desvío de 10 ms
func SleepBetweenShareContactsShort() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := MIN_SLEEP_BETWEEN_SH_CONTACTS + randSource.Intn(MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}

// Sleep de 13 minutos con un desvío de 20 segundos
func SleepBetweenShareContactsLarge() {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := 16*MIN_SLEEP_BETWEEN_SH_CONTACTS + randSource.Intn(2*MAX_OFFSET_SLEEP_BETWEEN_SH_CONTACTS)
	t := time.Duration(r) * time.Second
	time.Sleep(t)
}
