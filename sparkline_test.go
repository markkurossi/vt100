//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"fmt"
	"testing"
)

func TestSparkline(t *testing.T) {
	values := []int{83, 61, 33, 25, 12, 1, 0, 75}

	fmt.Printf("%s\n", Sparkline(nil))
	fmt.Printf("%s\n", Sparkline(values))
	fmt.Printf("%s\n", SparklineRange(0, 100, values))
	fmt.Printf("%s\n", SparklineRange(-100, 100, values))
	fmt.Printf("%s\n", SparklineRange(50, 50, values))
	fmt.Printf("%s\n", SparklineRange(25, 65, values))
}
