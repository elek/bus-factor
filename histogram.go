package main

import "sort"

type Histogram struct {
	Events map[string]int
}

type HistogramEntry struct {
	Key        string
	Occurrence int
}

func NewHistogram() Histogram {
	return Histogram{
		Events: make(map[string]int),
	}
}

func (h *Histogram) Increment(key string) {
	if _, found := h.Events[key]; !found {
		h.Events[key] = 1
	} else {
		h.Events[key] = h.Events[key] + 1
	}
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
func (h *Histogram) Sum() int {
	s := 0
	for _, e := range h.Events {
		s += e
	}
	return s
}
