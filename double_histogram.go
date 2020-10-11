package main

import "strings"

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

func (h *DoubleHistogram) FirstKey() string {
	res := ""
	for k, _ := range h.Histograms {
		if res == "" || strings.Compare(res, k) > 0 {
			res = k
		}

	}
	return res
}
