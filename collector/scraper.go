package collector

import "github.com/prometheus/client_golang/prometheus"

type Scraper interface {
	// Name of the Scraper. Should be unique.
	Name() string

	// Help describes the role of the Scraper.
	// Example: "Collect from SHOW ENGINE INNODB STATUS"
	//Help() string

	// Scrape collects data from client and sends it over channel as prometheus metric.
	Scrape(domains []Domain,ch chan<- prometheus.Metric) error
}
