package levee

import (
	"math"
	"sync"
	"time"
)

// maxDeviation defines the max ratio of the tick interval that is allowed for rounding
const maxDeviation = 0.05

// Bucket implements a token bucket that fills with new tokens at fixed rate.
// This implementation is based on an algorithm found here: http://en.wikipedia.org/wiki/Token_bucket.
// This algorithm allows us to compute the state of the bucket only when the consumer requests new tokens
type Bucket struct {

	// fillInterval stores the duration that should be passed for each token to be added
	fillInterval	time.Duration

	// firstTickTime stores the exact moment when the bucket created
	firstTickTime	time.Time

	// maxOffset stores the absolute maximum offset that we allow for rounding
	maxOffset		time.Duration

	// capacity stores the maximum amount of tokens that the bucket can hold
	capacity		int64

	// lastTick stores the current tick that has been computed upon the last request
	lastTick		int64

	// tokens store the currently available tokens
	tokens			int64

	// mu asserts that only one computation will happen at a time
	mu				sync.Mutex
}

// NewBucket returns a new bucket with the calculated fillInterval.
// The bucket accepts rate for new tokens to be added, a capacity that will
// limit the amount of tokens that can be hold by the bucket.
// It also requires the current time for testing purposes.
func NewBucket(rate float64, capacity int64, now time.Time) *Bucket {
	fillInterval := time.Duration(1e9 / rate)
	return &Bucket{
		fillInterval: 	fillInterval,
		firstTickTime: 	now,
		maxOffset:		time.Duration(maxDeviation * float64(fillInterval)),
		capacity:		capacity,
		mu: 			sync.Mutex{},
	}
}

// UpdateLimit updates the duration for new token to be added, defined by rate,
// and the capacity of the bucket
func (b *Bucket) UpdateLimit(rate float64, capacity int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fillInterval = time.Duration(1e9 / rate)
	b.capacity = capacity
}

// apply calculates and adjusts the current available tokens for the bucket
func (b *Bucket) apply(now time.Time) time.Duration {
	elapsed := now.Sub(b.firstTickTime)
	tick := int64(math.Floor(float64(elapsed) / float64(b.fillInterval)))
	terra := time.Duration(int64(elapsed) % int64(b.fillInterval))
	// if the duration terra is close to the next tick by less than the maximum offset
	// we accept it as a full tick
	if b.fillInterval - terra <= b.maxOffset {
		tick++
		terra = 0
	}
	b.adjustTick(tick)
	return terra
}

// apply calculates and adjusts the current available tokens for the bucket
func (b *Bucket) adjustTick(tick int64) {
	delta := tick - b.lastTick
	if delta > 0 {
		b.tokens += delta
		b.lastTick = tick
	}
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
}
// Take applies the current tick and calculates the time required
// for the provided amount of tokens to become available.
// it requires the current moment for testing purpose
func (b *Bucket) Take(now time.Time, amount int64) time.Duration {
	if amount <= 0 {
		return 0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	terra := b.apply(now)
	b.tokens -= amount
	if b.tokens >= 0 {
		return 0
	}
	targetTime := now.Add(time.Duration(-b.tokens) * b.fillInterval - terra)
	waitTime := targetTime.Sub(now)
	return waitTime
}
