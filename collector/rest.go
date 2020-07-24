package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

type Rest struct {
	Client http.Client
	Url    string
}

func (r *Rest) Name() string {
	r.Url = "rest"
	return r.Url
}

func (r *Rest) Scrape(domains []Domain, ch chan<- prometheus.Metric) error {
	r.Client = http.Client{}
	r.Client.Timeout = time.Duration(time.Millisecond * 10000)
	falsedomain := []Domain{}
	for _, v := range domains {
		suffix := "rest"
		r.Url = v.Ip + "/" + suffix
		req, err := http.NewRequest("GET", r.Url, nil)
		if err != nil {
			falsedomain = append(falsedomain,v)
			continue
		}
		start := time.Now()
		resp, err := r.Client.Do(req)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(
				//这里的label是固定标签 我们可以通过传递值传进去
				NewDesc(suffix, "status", suffix+"status ", []string{"domain"}, prometheus.Labels{"status": "超时"}),
				prometheus.GaugeValue,
				-1,
				//动态标签的值 可以有多个动态标签
				v.Domaininfo,
			)
			continue
		}
		end := time.Now().Sub(start).Milliseconds()
		ch <- prometheus.MustNewConstMetric(
			//这里的label是固定标签 我们可以通过传递值传进去
			NewDesc(suffix, "status", suffix+"status ", []string{"domain"}, prometheus.Labels{"status": resp.Status,"url":r.Url}),
			prometheus.GaugeValue,
			float64(end),
			//动态标签的值 可以有多个动态标签
			v.Domaininfo,
		)
	}
	if len(falsedomain) > 0 {
		log.Printf("此次收集失败的有%v\n",falsedomain)
	}
	return nil
}
