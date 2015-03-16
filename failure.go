// Package failure implements the Phi Accrual Failure Detector
/*

Please see http://ddg.jaist.ac.jp/pub/HDY+04.pdf

To use the failure detection algorithm, you need a heartbeat loop that will
call Ping() at regular intervals.  At any point, you can call Phi() which will
report how suspicious it is that a heartbeat has not been heard since the last
time Ping() was called.

See https://issues.apache.org/jira/browse/CASSANDRA-2597 for an explanation
of the simplified math used in Phi().
*/

package failure

import (
	"math"
	"sync"
	"time"

	"github.com/dgryski/go-onlinestats"
)

// Detector is a failure detector
type Detector struct {
	w          *onlinestats.Windowed
	last       time.Time
	minSamples int
	mu         sync.Mutex
}

// New returns a new failure detector that considers the last windowSize
// samples, and ensures there are at least minSamples in the window before
// returning an answer
func New(windowSize, minSamples int) *Detector {
	d := &Detector{
		w:          onlinestats.NewWindowed(windowSize),
		minSamples: minSamples,
	}

	return d
}

// Ping registers a heart-beat at time now
func (d *Detector) Ping(now time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.last.IsZero() {
		d.w.Push(now.Sub(d.last).Seconds())
	}
	d.last = now
}

// See https://issues.apache.org/jira/browse/CASSANDRA-2597 for an explanation
// of the math.

var PHI_FACTOR = 1.0 / math.Log(10.0)

// Phi calculates the suspicion level at time 'now' that the remote end has failed
func (d *Detector) Phi(now time.Time) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.w.Len() < d.minSamples {
		return 0
	}

	t := now.Sub(d.last).Seconds()

	return PHI_FACTOR * t / d.w.Mean()
}
