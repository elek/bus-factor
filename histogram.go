package main

import "sort"

type Histogram struct {
	Events map[string]float64
}

type HistogramEntry struct {
	Key        string
	Occurrence float64
}

func NewHistogram() Histogram {
	return Histogram{
		Events: make(map[string]float64),
	}
}

func (h *Histogram) Increment(key string) {
	if _, found := h.Events[key]; !found {
		h.Events[key] = 1
	} else {
		h.Events[key] = h.Events[key] + 1
	}
}

func (h *Histogram) Set(key string, value float64) {
	h.Events[key] = value
}

func (h *Histogram) SortedView() []HistogramEntry {
	result := make([]HistogramEntry, 0)
	for k, v := range h.Events {
		result = append(result, HistogramEntry{k, v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Occurrence > result[j].Occurrence
	})
	return result
}

func (h *Histogram) Sum() float64 {
	s := 0.0
	for _, e := range h.Events {
		s += e
	}
	return s
}

func (h *Histogram) Max() float64 {
	m := 0.0
	for _, e := range h.Events {
		if e > m {
			m = e
		}
	}
	return m
}

func (h *Histogram) Get(contributor string) (float64, bool) {
	k, f := h.Events[contributor]
	return k, f
}

func (h *Histogram) GetOrZero(contributor string) float64 {
	k, f := h.Events[contributor]
	if f {
		return k
	} else {
		return 0.0
	}
}