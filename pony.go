package main

func (h *Histogram) PonyDevs() []HistogramEntry {
	result := make([]HistogramEntry, 0)
	limit := h.Sum() / 2
	i := 0.0
	for _, entry := range h.SortedView() {
		i += entry.Occurrence
		result = append(result, entry)
		if i > limit {
			return result
		}
	}
	return result
}
