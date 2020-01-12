package monitoring

import (
	"errors"
	"net/http"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/storage/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type SqlMetrics struct {
	logger       *zap.SugaredLogger
	exporterPort string // prometheus http metrics exporter port, if empty string exporter not be run
	storage      *sql.Storage
	registry     *prometheus.Registry
}

func NewSqlMetrics(storage *sql.Storage, exporterPort string, logger *zap.SugaredLogger) (*SqlMetrics, error) {
	if storage == nil {
		return nil, errors.New("storage is required")
	}

	registry := prometheus.NewRegistry()

	opts := prometheus.GaugeOpts{
		Subsystem: "pg",
		Name:      "live_rows",
		Help:      "Estimate number of live rows",
	}

	counterFn := prometheus.NewGaugeFunc(opts, func() float64 {
		value, err := storage.GetStatValueNLiveTup()
		if err != nil {
			if logger != nil {
				logger.Errorf("GetStatValueNLiveTup error %s", err)
			}
		}
		return float64(value)
	})

	err := registry.Register(counterFn)
	if err != nil {
		if logger != nil {
			logger.Errorf("can't register gauge function `%s` metric: %s", opts.Name, err)
		}
	}

	return &SqlMetrics{
		logger:       logger,
		exporterPort: exporterPort,
		storage:      storage,
		registry:     registry,
	}, nil
}

func (m *SqlMetrics) RegisterExporter() {
	go func() {

		if m.logger != nil {
			m.logger.Infof("Try run sql (pg) metrics prometheus exporter run on port %s", m.exporterPort)
		}

		// HTTP exporter for prometheus for sql (pg) metrics
		handler := promhttp.InstrumentMetricHandler(
			m.registry,
			promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
		)

		err := http.ListenAndServe(":"+m.exporterPort, handler)
		if err != nil {
			if m.logger != nil {
				m.logger.Errorf("SqlMetrics.RegisterExporter, http listen and serve failed, return error %s", err)
			}
		}
	}()
}
