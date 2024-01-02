package util

import "math"

func MinLengthOfStrings(strings []string) int {
	min := math.MaxInt
	for _, s := range strings {
		if len(s) < min {
			min = len(s)
		}
	}

	return min
}

func MaxLengthOfStrings(strings []string) int {
	max := math.MinInt
	for _, s := range strings {
		if len(s) > max {
			max = len(s)
		}
	}

	return max
}
