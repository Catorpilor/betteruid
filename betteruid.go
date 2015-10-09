package betteruid

import (
	"math/rand"
	"sync"
	"time"
)

const (
	PUSH_CHARS = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
)

var (
	//Timestamp of last push, used to prevent local collision if you push
	//twice in one ms
	lastPushTime int64
	//We generate 72-bits of randomness which get turned into 12 characters
	//and appended to the timestamp to prevent collisions with other clients.
	//We store the last characters we generated because in the event of
	//collision, we'll use those same characters except increment by one.
	lastRandChars [12]int
	mu            sync.Mutex
	rnd           *rand.Rand
)

func init() {
	//generate new seeds
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 12; i++ {
		lastRandChars[i] = rnd.Intn(64)
	}
}

func New() string {
	var id [20]byte
	mu.Lock()
	timeMS := time.Now().UnixNano() / 1e6
	if timeMS == lastPushTime {
		//increment lastRandChars
		for i := 0; i < 12; i++ {
			lastRandChars[i]++
			if lastRandChars[i] < 64 {
				break
			}
			lastRandChars[i] = 0
		}
	}
	lastPushTime = timeMS
	//put random as second part
	for i := 0; i < 12; i++ {
		id[19-i] = PUSH_CHARS[lastRandChars[i]]
	}
	mu.Unlock()

	//put current at beginning
	for i := 7; i >= 0; i-- {
		n := int(timeMS % 64)
		id[i] = PUSH_CHARS[n]
		timeMS = timeMS / 64
	}
	return string(id[:])
}
