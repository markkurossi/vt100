//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"strings"
)

// HBlock creates a horizontal block that is width characters
// long. The fract argument specifies the fraction ([0...1]) of the
// width that is rendered with the Unicode block drawing characters,
// starting from the left edge. The remaining empty fraction of the
// block is padded with the empty rune.
func HBlock(width int, fract float64, empty rune) string {
	if fract < 0 {
		fract = 0
	}
	if fract > 1 {
		fract = 1
	}
	w8 := float64(width * 8)
	w := int(w8 * fract)

	var sb strings.Builder

	var i int
	for i = 0; i < w/8; i++ {
		sb.WriteRune(0x2588)
	}
	if w%8 > 0 {
		sb.WriteRune(0x2590 - rune(w%8))
		i++
	}
	for ; i < width; i++ {
		sb.WriteRune(empty)
	}

	return sb.String()
}
