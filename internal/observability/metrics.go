package observability

import (
	"sync"
)

// MetricsCollector collects and stores metrics.
// This is a simple in-memory implementation that can be replaced with Prometheus, etc.
type MetricsCollector struct {
	mu sync.RWMutex
	// Counters
	counters map[string]int64
	// Histograms (as simple slices for now)
	histograms map[string][]float64
	// Gauges
	gauges map[string]float64
}

// NewMetricsCollector creates a new metrics collector.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		counters:   make(map[string]int64),
		histograms: make(map[string][]float64),
		gauges:     make(map[string]float64),
	}
}

// IncCounter increments a counter by 1.
func (m *MetricsCollector) IncCounter(name string, labels map[string]string) {
	key := m.buildKey(name, labels)
	m.mu.Lock()
	m.counters[key]++
	m.mu.Unlock()
}

// AddCounter adds a value to a counter.
func (m *MetricsCollector) AddCounter(name string, value int64, labels map[string]string) {
	key := m.buildKey(name, labels)
	m.mu.Lock()
	m.counters[key] += value
	m.mu.Unlock()
}

// ObserveHistogram adds a value to a histogram.
func (m *MetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
	key := m.buildKey(name, labels)
	m.mu.Lock()
	m.histograms[key] = append(m.histograms[key], value)
	m.mu.Unlock()
}

// SetGauge sets a gauge value.
func (m *MetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
	key := m.buildKey(name, labels)
	m.mu.Lock()
	m.gauges[key] = value
	m.mu.Unlock()
}

// GetCounter returns the current value of a counter.
func (m *MetricsCollector) GetCounter(name string, labels map[string]string) int64 {
	key := m.buildKey(name, labels)
	m.mu.RLock()
	value := m.counters[key]
	m.mu.RUnlock()
	return value
}

// GetHistogram returns a copy of the histogram values.
func (m *MetricsCollector) GetHistogram(name string, labels map[string]string) []float64 {
	key := m.buildKey(name, labels)
	m.mu.RLock()
	values := make([]float64, len(m.histograms[key]))
	copy(values, m.histograms[key])
	m.mu.RUnlock()
	return values
}

// GetGauge returns the current value of a gauge.
func (m *MetricsCollector) GetGauge(name string, labels map[string]string) float64 {
	key := m.buildKey(name, labels)
	m.mu.RLock()
	value := m.gauges[key]
	m.mu.RUnlock()
	return value
}

// buildKey creates a key from metric name and labels.
func (m *MetricsCollector) buildKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	// Simple key building - in practice this would be more sophisticated
	key := name
	for k, v := range labels {
		key += "," + k + "=" + v
	}
	return key
}

// MetricsAdapter adapts our MetricsCollector to the Metrics interface.
type MetricsAdapter struct {
	collector *MetricsCollector
}

// NewMetricsAdapter creates a new metrics adapter.
func NewMetricsAdapter(collector *MetricsCollector) *MetricsAdapter {
	return &MetricsAdapter{collector: collector}
}

// IncCounter increments a counter by 1.
func (a *MetricsAdapter) IncCounter(name string, labels map[string]string) {
	a.collector.IncCounter(name, labels)
}

// AddCounter adds a value to a counter.
func (a *MetricsAdapter) AddCounter(name string, value int64, labels map[string]string) {
	a.collector.AddCounter(name, value, labels)
}

// ObserveHistogram adds a value to a histogram.
func (a *MetricsAdapter) ObserveHistogram(name string, value float64, labels map[string]string) {
	a.collector.ObserveHistogram(name, value, labels)
}

// SetGauge sets a gauge value.
func (a *MetricsAdapter) SetGauge(name string, value float64, labels map[string]string) {
	a.collector.SetGauge(name, value, labels)
}