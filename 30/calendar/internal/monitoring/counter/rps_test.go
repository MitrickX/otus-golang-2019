package counter

import (
	"math/rand"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func newTestRpsCounter() *RpsCounter {

	rpsGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rps",
		Help: "Request per second",
	})

	counter := NewRpsCounter(rpsGauge)
	counter.ticker.Stop()
	counter.ticker = time.NewTicker(200 * time.Millisecond)
	return counter
}

func TestRpsCounting(t *testing.T) {
	counter := newTestRpsCounter()

	counter.Run()

	counter.Inc()
	counter.Inc()
	counter.Inc()

	time.Sleep(250 * time.Millisecond)

	counter.Inc()
	counter.Inc()
	counter.Inc()
	counter.Inc()

	counter.OutputMaxValues()

	expectedCount := 4
	count := int(testutil.ToFloat64(counter.gauge))
	if expectedCount != count {
		t.Errorf("unexpected count %d instreadof %d", count, expectedCount)
	}

	counter.OutputMaxValues()

	count = int(testutil.ToFloat64(counter.gauge))
	if count != 0 {
		t.Errorf("unexpected count %d instreadof %d", count, 0)
	}
}

func TestRpsConcurrentCounting(t *testing.T) {
	counter := newTestRpsCounter()

	counter.Run()

	incFunc := func() {
		d := rand.Intn(100)
		time.Sleep(time.Duration(d) * time.Millisecond)
		counter.Inc()
	}

	go incFunc()
	go incFunc()
	go incFunc()

	time.Sleep(250 * time.Millisecond)

	go incFunc()
	go incFunc()
	go incFunc()
	go incFunc()

	time.Sleep(250 * time.Millisecond)

	counter.OutputMaxValues()

	expectedCount := 4
	count := int(testutil.ToFloat64(counter.gauge))
	if expectedCount != count {
		t.Errorf("unexpected count %d instreadof %d", count, expectedCount)
	}

	counter.OutputMaxValues()

	count = int(testutil.ToFloat64(counter.gauge))
	if count != 0 {
		t.Errorf("unexpected count %d instreadof %d", count, 0)
	}
}
