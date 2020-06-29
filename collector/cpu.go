package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/load"
	"log"
)

type CpuLoad struct {}

// cpu info

func (CpuLoad) Name() string {
	return namespace + "_cpu_load"
}

func (CpuLoad) Scrape(ch chan <- prometheus.Metric) error  {
	cpuload,err  := load.Avg()
	if err != nil {
		log.Printf("cpu load is not ,%v\n",err)
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		NewDesc("cpu_load","one","cpu load",prometheus.Labels{"host":"12345"}),
		prometheus.GaugeValue,
		cpuload.Load1,
	)
	ch <- prometheus.MustNewConstMetric(
		NewDesc("cpu_load","five","cpu load",nil),
		prometheus.GaugeValue,
		cpuload.Load5,
	)
	ch <- prometheus.MustNewConstMetric(
		NewDesc("cpu_load","fifteen","cpu load",nil),
		prometheus.GaugeValue,
		cpuload.Load15,
	)
	return nil
}


