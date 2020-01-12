package monitoring

import (
	"net/http"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/monitoring/counter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type SenderMetrics struct {
	logger       *zap.SugaredLogger
	exporterPort string // prometheus http metrics exporter port, if empty string exporter not be run
	registry     *prometheus.Registry
	rpsCounter   *counter.RpsCounter
}

func NewSenderMetrics(exporterPort string, logger *zap.SugaredLogger) *SenderMetrics {

	registry := prometheus.NewRegistry()

	opts := prometheus.GaugeOpts{
		Subsystem: "sender",
		Name:      "requests_per_second",
		Help:      "Max count of requests per second per scrape_interval",
	}

	rpsGauge := prometheus.NewGauge(opts)

	var rpsCounter *counter.RpsCounter
	if err := registry.Register(rpsGauge); err != nil {
		if logger != nil {
			logger.Errorf("can't register rps gauge `%s` metric: %s", opts.Name, err)
		}
	} else {
		rpsCounter = counter.NewRpsCounter(rpsGauge)
	}

	m := &SenderMetrics{
		logger:       logger,
		exporterPort: exporterPort,
		registry:     registry,
		rpsCounter:   rpsCounter,
	}

	return m
}

func (m *SenderMetrics) IncSenderCounter() {
	if m.rpsCounter != nil {
		m.rpsCounter.Inc()
	}
}

func (m *SenderMetrics) Stop() {
	if m.rpsCounter != nil {
		m.rpsCounter.Stop()
	}
}

func (m *SenderMetrics) RegisterExporter() {

	if m.rpsCounter != nil {
		m.rpsCounter.Run()
	}

	go func() {

		if m.logger != nil {
			m.logger.Infof("Try run sender metrics prometheus exporter run on port %s", m.exporterPort)
		}

		// HTTP exporter for prometheus for sender metrics
		handler := func(w http.ResponseWriter, r *http.Request) {

			if m.rpsCounter != nil {
				// out current max value of counter to gauge metric
				m.rpsCounter.OutputMaxValues()
			}

			prom := promhttp.InstrumentMetricHandler(
				m.registry,
				promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
			)
			prom.ServeHTTP(w, r)
		}

		err := http.ListenAndServe(":"+m.exporterPort, http.HandlerFunc(handler))
		if err != nil {
			if m.logger != nil {
				m.logger.Errorf("SenderMetrics.RegisterExporter, http listen and serve failed, return error %s", err)
			}
		}
	}()
}
