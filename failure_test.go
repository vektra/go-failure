package failure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vektra/neko"
)

func TestDetector(t *testing.T) {
	n := neko.Start(t)

	n.It("returns a good value when getting good heartbeats", func() {
		d := New(20, 10)

		for i := time.Duration(0); i <= 20; i++ {
			d.Ping(time.Now().Add((i * time.Second)))
		}

		phi := d.Phi(time.Now().Add(21 * time.Second))

		assert.True(t, phi < 1.0)
	})

	n.It("gets worried when heartbeats have slowed", func() {
		d := New(20, 10)

		for i := time.Duration(0); i <= 20; i++ {
			d.Ping(time.Now().Add(i * time.Second))
		}

		phi := d.Phi(time.Now().Add(40 * time.Second))

		assert.True(t, phi > 8.0)
	})

	n.It("approaches infinity given no new pings", func() {
		d := New(20, 10)

		for i := time.Duration(0); i <= 20; i++ {
			d.Ping(time.Now().Add(i * time.Second))
		}

		last := 0.0

		for i := time.Duration(1); i < 100; i++ {
			phi := d.Phi(time.Now().Add((i * 100) * time.Second))
			assert.True(t, phi > last)
			last = phi
		}
	})

	n.Meow()
}
