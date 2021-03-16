//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"math"
	"strings"
)

// Sparkline creates a histogram chart of values. The chart is scaled
// to [min...max] values in the values array.
func Sparkline(values []int) string {
	min := math.MaxInt32
	max := math.MinInt32
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return SparklineRange(min, max, values)
}

// SparklineRange creates a histogram chart of values. The chart is
// scaled to [min...max]. Values smaller than min are rendered with
// space (u0020) and values larger than max are rendered with 'Light
// Shade' (u2591).
func SparklineRange(min, max int, values []int) string {
	if len(values) == 0 {
		return ""
	}
	if max < min {
		min = max
	}

	delta := max - min

	var sb strings.Builder
	for _, v := range values {
		var tick rune
		if v < min {
			tick = ' '
		} else if v > max {
			tick = 0x2591
		} else if delta == 0 {
			tick = 0x2581 + 4
		} else {
			tick = rune(0x2581 + (v-min)*7/delta)
		}
		sb.WriteRune(tick)
	}
	return sb.String()
}
