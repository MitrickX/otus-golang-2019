package counter

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RpsVecCounter struct {
	values   map[string][2]int // index 0 for current value, index 1 for maximum value
	mx       sync.Mutex
	ticker   *time.Ticker
	gaugeVec *prometheus.GaugeVec
}

func NewRpsVecCounter(gaugeVec *prometheus.GaugeVec) *RpsVecCounter {
	values := make(map[string][2]int)
	return &RpsVecCounter{
		values:   values,
		mx:       sync.Mutex{},
		ticker:   time.NewTicker(time.Second),
		gaugeVec: gaugeVec,
	}
}

func (r *RpsVecCounter) Run() {
	go func() {
		for range r.ticker.C {
			r.sync()
		}
	}()
}

func (r *RpsVecCounter) sync() {
	r.mx.Lock()
	defer r.mx.Unlock()
	for key, pair := range r.values {
		// if current value > maximum value
		if pair[0] > pair[1] {
			pair[1] = pair[0]
		}
		// clear current value
		pair[0] = 0
		r.values[key] = pair
	}
}

func (r *RpsVecCounter) Stop() {
	r.ticker.Stop()
}

func (r *RpsVecCounter) Inc(key string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	pair := r.values[key]
	pair[0]++
	r.values[key] = pair
}

func (r *RpsVecCounter) OutputMaxValues() {
	r.sync()
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.gaugeVec == nil {
		return
	}
	for key, pair := range r.values {
		maxValue := pair[1]
		r.gaugeVec.WithLabelValues(key).Set(float64(maxValue))
		pair[1] = 0
		r.values[key] = pair
	}
}
