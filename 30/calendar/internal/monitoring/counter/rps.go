package counter

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RpsCounter struct {
	value    int
	maxValue int
	mx       sync.Mutex
	ticker   *time.Ticker
	gauge    prometheus.Gauge
}

func NewRpsCounter(gauge prometheus.Gauge) *RpsCounter {
	return &RpsCounter{
		mx:     sync.Mutex{},
		ticker: time.NewTicker(time.Second),
		gauge:  gauge,
	}
}

func (r *RpsCounter) Run() {
	go func() {
		for range r.ticker.C {
			r.sync()
		}
	}()
}

func (r *RpsCounter) sync() {
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.value > r.maxValue {
		r.maxValue = r.value
	}
	r.value = 0
}

func (r *RpsCounter) Stop() {
	r.ticker.Stop()
}

func (r *RpsCounter) Inc() {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.value++
}

func (r *RpsCounter) OutputMaxValues() {
	r.sync()
	r.mx.Lock()
	defer r.mx.Unlock()
	if r.gauge == nil {
		return
	}
	r.gauge.Set(float64(r.maxValue))
	r.maxValue = 0
}
