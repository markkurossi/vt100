//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"testing"
)

var widthTests = []struct {
	i string
	o int
}{
	{
		i: "Hello, world!",
		o: 13,
	},
	{
		i: "\x1b[30;41mHello, world!\x1b[0m",
		o: 13,
	},
}

func TestDisplayWidth(t *testing.T) {
	for idx, test := range widthTests {
		w, err := DisplayWidth(test.i)
		if err != nil {
			t.Errorf("test %d failed: %s", idx, err)
			continue
		}
		if w != test.o {
			t.Errorf("test %d failed: got %d, expected %d", idx, w, test.o)
		}
	}
}
