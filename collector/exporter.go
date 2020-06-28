package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	name string= "test_exporter"
	namespace string = "test"
	exporter string = "exporter"
)
func Name() string {
	return  name
}

type Exporter struct {
	metrics Metrics
	scrapers []Scraper
}

type Metrics struct {
	TotalScrapes prometheus.Counter
	ScrapeErrors *prometheus.CounterVec
	Error        prometheus.Gauge
	//HarborUp     prometheus.Gauge
}

func NewDesc(subsystem, name, help string,label prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		help, nil, label,
	)
}

var  _ prometheus.Collector = (*Exporter)(nil)

func New(metrics Metrics, scrapers []Scraper) *Exporter{
	return &Exporter{
		metrics: metrics,
		scrapers: scrapers,
	}
}

func NewMetrics() Metrics {
	subsystem := exporter
	return Metrics{
		TotalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "scrapes_total",
			Help:      "Total number of times harbor was scraped for metrics.",
		}),
		ScrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occurred scraping a harbor.",
		}, []string{"collector"}),
		Error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from harbor resulted in an error (1 for error, 0 for success).",
		}),
	}
}



func (e *Exporter) Describe(ch chan <-  *prometheus.Desc) {
	ch <- e.metrics.TotalScrapes.Desc()
	ch <- e.metrics.Error.Desc()
	e.metrics.ScrapeErrors.Describe(ch)
}

func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
	e.scrape(ch)
	ch <- e.metrics.TotalScrapes
	ch <- e.metrics.Error
	e.metrics.ScrapeErrors.Collect(ch)
}

func (e *Exporter) scrape(ch chan <- prometheus.Metric) {
	var (
		wg sync.WaitGroup
		err error
	)

	defer wg.Wait()
	for _,scraper := range  e.scrapers {
		wg.Add(1)
		//使用匿名函数 并且并发的收集指标
		go func(scraper Scraper) {
			defer wg.Done()
			label := scraper.Name()
			err = scraper.Scrape(ch)
			if err != nil {
				log.WithField("scraper", scraper.Name()).Error(err)
				e.metrics.ScrapeErrors.WithLabelValues(label).Inc()
				e.metrics.Error.Set(1)
			}
		}(scraper)
	}
}