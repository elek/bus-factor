package main

import (
	"sort"
	"strings"
)

type DoubleHistogram struct {
	Histograms map[string]*Histogram
}

func (h *DoubleHistogram) Get(name string) *Histogram {
	hist, found := h.Histograms[name]
	if !found {
		nh := NewHistogram()
		hist = &nh
		h.Histograms[name] = hist
	}
	return hist
}

func (h *DoubleHistogram) LastKey() string {
	res := ""
	for k, _ := range h.Histograms {
		if res == "" || strings.Compare(res, k) < 0 {
			res = k
		}

	}
	return res
}
func (h *DoubleHistogram) Keys() []string {
	keys := make([]string, 0)
	for k, _ := range h.Histograms {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (h *DoubleHistogram) FirstKey() string {
	res := ""
	for k, _ := range h.Histograms {
		if res == "" || strings.Compare(res, k) > 0 {
			res = k
		}

	}
	return res
}

func (h *DoubleHistogram) Timelined() (DoubleHistogram, error) {

	keys := h.Keys()

	contributions := NewHistogram()
	result := NewDoubleHistogram()
	weight := 0.5
	for _, category := range keys {

		contributed := make(map[string]bool)
		for contributor, contribution := range h.Get(category).Events {
			existing := contributions.GetOrZero(contributor)
			contributions.Set(contributor, existing*(1-weight)+weight*float64(contribution))
			contributed[contributor] = true
		}

		adjustements := make(map[string]float64)
		for contributor, _ := range contributions.Events {
			if _, found := contributed[contributor]; !found {
				discounted := (1 - weight) * contributions.GetOrZero(contributor)
				if discounted >= 1 {
					adjustements[contributor] = discounted
				} else {
					adjustements[contributor] = 0
				}
			}
		}
		for k, v := range adjustements {
			contributions.Set(k, v)
		}
		snapshot := CopyHistogram(contributions)
		result.Set(category, &snapshot)
	}
	return result, nil
}

func (h *DoubleHistogram) Set(category string, histogram *Histogram) {
	h.Histograms[category] = histogram
}
