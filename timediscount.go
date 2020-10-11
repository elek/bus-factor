package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeDiscountedGitLog(path string) (Histogram, error) {
	histogram, err := ReadMonthlyGitLog(path)
	if err != nil {
		return Histogram{}, err
	}
	until := StringMonthToNum(histogram.LastKey())
	lastMonth := time.Now().Year()*12 + int(time.Now().Month()-1)
	if lastMonth < until {
		until = lastMonth
	}
	contributions := NewHistogram()

	weight := 0.5
	for i := StringMonthToNum(histogram.FirstKey()); i <= until; i++ {
		year := i / 12
		month := i%12 + 1
		yearString := fmt.Sprintf("%d-%02d", year, month)

		contributed := make(map[string]bool)
		for contributor, contribution := range histogram.Get(yearString).Events {
			existing := contributions.GetOrZero(contributor)
			contributions.Set(contributor, existing*(1-weight)+weight*float64(contribution))
			contributed[contributor] = true
		}

		adjustements := make(map[string]float64)
		for contributor, _ := range contributions.Events {
			if _, found := contributed[contributor]; !found {
				discounted := (1 - weight) * contributions.GetOrZero(contributor)
				if discounted > 0.01 {
					adjustements[contributor] = discounted
				} else {
					adjustements[contributor] = 0
				}
			}
		}
		for k, v := range adjustements {
			contributions.Set(k, v)
		}
	}
	return contributions, nil
}

func StringMonthToNum(monthString string) int {
	parts := strings.Split(monthString, "-")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	return year*12 + month - 1
}
