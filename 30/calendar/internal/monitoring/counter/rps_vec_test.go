package counter

import (
	"math/rand"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/client_golang/prometheus"
)

func newTestCounter() *RpsVecCounter {

	rpsGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_rps",
		Help: "Http request per second",
	}, []string{"method"})

	counter := NewRpsVecCounter(rpsGaugeVec)
	counter.ticker.Stop()
	counter.ticker = time.NewTicker(200 * time.Millisecond)
	return counter
}

func TestCounting(t *testing.T) {
	counter := newTestCounter()

	counter.Run()

	counter.Inc("/create_event")
	counter.Inc("/create_event")
	counter.Inc("/create_event")

	counter.Inc("/update_event")
	counter.Inc("/update_event")

	time.Sleep(250 * time.Millisecond)

	counter.Inc("/create_event")
	counter.Inc("/create_event")
	counter.Inc("/create_event")
	counter.Inc("/create_event")

	counter.Inc("/update_event")

	counter.OutputMaxValues()

	expectedCount := 4
	count := int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/create_event")))
	if expectedCount != count {
		t.Errorf("unexpected count %d instreadof %d for label `/create_event`", count, expectedCount)
	}

	expectedCount = 2
	count = int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/update_event")))
	if expectedCount != count {
		t.Errorf("unexpected count %d instreadof %d for label `/update_event`", count, expectedCount)
	}

	counter.OutputMaxValues()

	count = int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/create_event")))
	if count != 0 {
		t.Errorf("unexpected count %d instreadof %d for label `/create_event`", count, 0)
	}

	count = int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/update_event")))
	if count != 0 {
		t.Errorf("unexpected count %d instreadof %d for label `/update_event`", count, 0)
	}
}

func TestConcurrentCounting(t *testing.T) {
	counter := newTestCounter()

	counter.Run()

	incFunc := func(method string) {
		d := rand.Intn(100)
		time.Sleep(time.Duration(d) * time.Millisecond)
		counter.Inc(method)
	}

	go incFunc("/create_event")
	go incFunc("/create_event")
	go incFunc("/create_event")

	go incFunc("/update_event")
	go incFunc("/update_event")

	time.Sleep(250 * time.Millisecond)

	go incFunc("/create_event")
	go incFunc("/create_event")
	go incFunc("/create_event")
	go incFunc("/create_event")

	go incFunc("/update_event")

	time.Sleep(250 * time.Millisecond)

	counter.OutputMaxValues()

	expectedCount := 4
	count := int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/create_event")))
	if expectedCount != count {
		t.Errorf("unexpected count %d instreadof %d for label `/create_event`", count, expectedCount)
	}

	expectedCount = 2
	count = int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/update_event")))
	if expectedCount != count {
		t.Errorf("unexpected count %d instreadof %d for label `/update_event`", count, expectedCount)
	}

	counter.OutputMaxValues()

	count = int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/create_event")))
	if count != 0 {
		t.Errorf("unexpected count %d instreadof %d for label `/create_event`", count, 0)
	}

	count = int(testutil.ToFloat64(counter.gaugeVec.WithLabelValues("/update_event")))
	if count != 0 {
		t.Errorf("unexpected count %d instreadof %d for label `/update_event`", count, 0)
	}

}
