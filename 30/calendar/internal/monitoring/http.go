package monitoring

import (
	"net/http"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/monitoring/counter"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/slok/go-http-metrics/middleware"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type HttpMetrics struct {
	logger         *zap.SugaredLogger
	requestCounter prometheus.Counter
	rpsCounter     *counter.RpsVecCounter
	exporterPort   string // prometheus http metrics exporter port, if empty string exporter not be run
}

func NewHttpMetrics(exporterPort string, logger *zap.SugaredLogger) *HttpMetrics {
	rpsGaugeVecOpts := prometheus.GaugeOpts{
		Subsystem: "http",
		Name:      "requests_per_second",
		Help:      "Max count of requests per second per scrape_interval",
	}
	rpsGaugeVec := prometheus.NewGaugeVec(rpsGaugeVecOpts, []string{"method"})

	var rpsCounter *counter.RpsVecCounter
	if err := prometheus.Register(rpsGaugeVec); err != nil {
		if logger != nil {
			logger.Errorf("can't register rps gauge vector `%s` metric: %s", rpsGaugeVecOpts.Name, err)
		}
	} else {
		rpsCounter = counter.NewRpsVecCounter(rpsGaugeVec)
	}

	requestCounterOpts := prometheus.CounterOpts{
		Subsystem: "http",
		Name:      "requests_count",
		Help:      "Total number of requests to http service",
	}

	var requestCounter prometheus.Counter

	requestCounter = prometheus.NewCounter(requestCounterOpts)
	if err := prometheus.Register(requestCounter); err != nil {
		if logger != nil {
			logger.Errorf("can't register counter `%s` metric: %s", requestCounterOpts.Name, err)
		}
	}

	return &HttpMetrics{
		logger:         logger,
		rpsCounter:     rpsCounter,
		requestCounter: requestCounter,
		exporterPort:   exporterPort,
	}
}

func (m *HttpMetrics) RegisterMiddleware(next http.Handler) http.Handler {

	// metrics middleware
	metricsMiddleware := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	handler := metricsMiddleware.Handler("", next)

	// requests and requests per seconds counter metrics
	handler = m.counterMiddleware(handler)

	m.runMetricsExporter()

	return handler
}

func (m *HttpMetrics) counterMiddleware(next http.Handler) http.Handler {
	// if there is not registered metrics will not wrap handler
	if m.rpsCounter == nil && m.requestCounter == nil {
		return next
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.requestCounter != nil {
			m.requestCounter.Inc()
		}
		if m.rpsCounter != nil {
			m.rpsCounter.Inc(r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})

	if m.rpsCounter != nil {
		m.rpsCounter.Run()
	}

	return handler
}

func (m *HttpMetrics) runMetricsExporter() {
	go func() {

		if m.logger != nil {
			m.logger.Infof("Try run http metrics prometheus exporter run on port %s", m.exporterPort)
		}

		// HTTP exporter for prometheus
		handler := func(w http.ResponseWriter, r *http.Request) {
			if m.rpsCounter != nil {
				m.rpsCounter.OutputMaxValues()
			}
			next := promhttp.Handler()
			next.ServeHTTP(w, r)
		}

		err := http.ListenAndServe(":"+m.exporterPort, http.HandlerFunc(handler))
		if err != nil {
			if m.logger != nil {
				m.logger.Errorf("HttpMetrics.runMetricsExporter, http listen and serve failed, return error %s", err)
			}
		}
	}()
}
