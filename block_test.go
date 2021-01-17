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

func TestHorizontal(t *testing.T) {
	values := []float64{.83, .61, .33, .25, .12, .1, .05}

	for _, v := range values {
		fmt.Printf("%2d: %s\n", int(v*100), HBlock(40, v, ' '))
	}
}
