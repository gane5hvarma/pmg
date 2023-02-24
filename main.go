package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gane5hvarma/pmg/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := config.ParseMetrics()
	generateMetrics(config.Metrics)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9000", nil)

}
func generateMetrics(metrics []config.Metric) {
	for _, metric := range metrics {
		switch metric.MetricType {
		case "counter":
			NewCounterMetric(metric)
		case "gauge":
			fmt.Println("came here")
			NewGaugeMetric(metric)
		case "histogram":
			NewHistogramMetric(metric)
		}
	}
}

func NewCounterMetric(metric config.Metric) {
	newCounterMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metric.Name,
		Help: metric.Help,
	}, metric.Labels)
	prometheus.Register(newCounterMetric)

	for _, metricGenerator := range metric.Generators {
		go func(metricGenerator config.Generator) {
			for {
				labels := metric.Labels
				var labelValues []string
				for _, label := range labels {
					labelValues = append(labelValues, metricGenerator.Labels[label])
				}
				newCounterMetric.WithLabelValues(labelValues...).Inc()
				time.Sleep(time.Duration(metricGenerator.Freq) * time.Second)
			}
		}(metricGenerator)
	}

}

func NewGaugeMetric(metric config.Metric) {
	newGaugeMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metric.Name,
		Help: metric.Help,
	}, metric.Labels)
	prometheus.Register(newGaugeMetric)

	for _, metricGenerator := range metric.Generators {
		go func(metricGenerator config.Generator) {
			for {
				labels := metric.Labels
				var labelValues []string
				for _, label := range labels {
					labelValues = append(labelValues, metricGenerator.Labels[label])
				}
				switch metricGenerator.Method {
				case "add":
					newGaugeMetric.WithLabelValues(labelValues...).Add(metricGenerator.Value)
				case "sub":
					newGaugeMetric.WithLabelValues(labelValues...).Sub(metricGenerator.Value)
				case "set":
					newGaugeMetric.WithLabelValues(labelValues...).Set(metricGenerator.Value)
				}
				time.Sleep(time.Duration(metricGenerator.Freq) * time.Second)
			}
		}(metricGenerator)
	}
}

func NewHistogramMetric(metric config.Metric) {
	newHistogramMetric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: metric.Name,
		Help: metric.Help,
	}, metric.Labels)
	prometheus.Register(newHistogramMetric)

	for _, metricGenerator := range metric.Generators {
		go func(metricGenerator config.Generator) {
			for {
				labels := metric.Labels
				var labelValues []string
				for _, label := range labels {
					labelValues = append(labelValues, metricGenerator.Labels[label])
				}
				newHistogramMetric.WithLabelValues(labelValues...).Observe(metricGenerator.Value)
				time.Sleep(time.Duration(metricGenerator.Freq) * time.Second)
			}
		}(metricGenerator)
	}
}
