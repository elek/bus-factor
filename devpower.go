package main

func (h *Histogram) DevPower() float64 {
	result := float64(0)
	base := float64(0)
	weight := float64(1)
	previous := float64(1)
	for _, entry := range h.SortedView() {
		if base == 0 {
			base = entry.Occurrence
		} else {
			if entry.Occurrence > 0 {
				weight = weight * entry.Occurrence / previous
			} else {
				weight = 0
			}
		}
		previous = entry.Occurrence

		result += (entry.Occurrence / base) * weight
	}
	return result
}

func (h *Histogram) LeadFactor() float64 {
	result := float64(0)
	base := float64(0)
	for _, entry := range h.SortedView() {
		if base == 0 {
			base = entry.Occurrence
		}
		result += entry.Occurrence / base
	}
	return result
}

func NewDoubleHistogram() DoubleHistogram {
	return DoubleHistogram{
		Histograms: make(map[string]*Histogram),
	}
}
