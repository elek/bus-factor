package main

import "math"

func (h *Histogram) HHI() (res float64) {
	sum := h.Sum()

	for _, entry := range h.SortedView() {
		res += math.Pow(entry.Occurrence/sum, 2)
	}
	return res
}
