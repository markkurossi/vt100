//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"fmt"
	"strings"
	"testing"
)

var widthTests = []struct {
	i string
	o []string
	w int
	h int
}{
	{
		i: "Hello, world!",
		o: []string{"Hello, world!"},
		w: 13,
		h: 1,
	},
	{
		i: "\x1b[30;41mHello, world!\x1b[0m",
		o: []string{"Hello, world!"},
		w: 13,
		h: 1,
	},
	{
		i: "\x1b[?3l\x1b#8",
		o: []string{
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
		},
		w: 80,
		h: 24,
	},
}

func TestDisplayWidth(t *testing.T) {
	for idx, test := range widthTests {
		w, h, err := DisplayWidth(test.i)
		if err != nil {
			t.Errorf("test %d: DisplayWidth failed: %s", idx, err)
			continue
		}
		if w != test.w || h != test.h {
			t.Errorf("test %d failed: got %d,%d, expected %d,%d",
				idx, w, h, test.w, test.h)
		}

		trimmed, err := Trim(test.i)
		if err != nil {
			t.Errorf("test %d: Trim failed: %s", idx, err)
			continue
		}
		if len(trimmed) != len(test.o) {
			t.Errorf("test %d: Trim: got %d lines, expected %d",
				idx, len(trimmed), len(test.o))
			continue
		}
		for i, l := range trimmed {
			if l != test.o[i] {
				t.Errorf("test %d: Trim: line %d differs", idx, i)
			}
		}
	}
}

var emulTests = []struct {
	input  string
	output string
}{
	{
		input: `stdout:
00000000  1b 5b 3f 33 6c 1b 23 38  1b 5b 39 3b 31 30 48 1b  |.[?3l.#8.[9;10H.|
00000010  5b 31 4a 1b 5b 31 38 3b  36 30 48 1b 5b 30 4a 1b  |[1J.[18;60H.[0J.|
00000020  5b 31 4b 1b 5b 39 3b 37  31 48 1b 5b 30 4b 1b 5b  |[1K.[9;71H.[0K.[|
00000030  31 30 3b 31 30 48 1b 5b  31 4b 1b 5b 31 30 3b 37  |10;10H.[1K.[10;7|
00000040  31 48 1b 5b 30 4b 1b 5b  31 31 3b 31 30 48 1b 5b  |1H.[0K.[11;10H.[|
00000050  31 4b 1b 5b 31 31 3b 37  31 48 1b 5b 30 4b 1b 5b  |1K.[11;71H.[0K.[|
00000060  31 32 3b 31 30 48 1b 5b  31 4b 1b 5b 31 32 3b 37  |12;10H.[1K.[12;7|
00000070  31 48 1b 5b 30 4b 1b 5b  31 33 3b 31 30 48 1b 5b  |1H.[0K.[13;10H.[|
00000080  31 4b 1b 5b 31 33 3b 37  31 48 1b 5b 30 4b 1b 5b  |1K.[13;71H.[0K.[|
00000090  31 34 3b 31 30 48 1b 5b  31 4b 1b 5b 31 34 3b 37  |14;10H.[1K.[14;7|
000000a0  31 48 1b 5b 30 4b 1b 5b  31 35 3b 31 30 48 1b 5b  |1H.[0K.[15;10H.[|
000000b0  31 4b 1b 5b 31 35 3b 37  31 48 1b 5b 30 4b 1b 5b  |1K.[15;71H.[0K.[|
000000c0  31 36 3b 31 30 48 1b 5b  31 4b 1b 5b 31 36 3b 37  |16;10H.[1K.[16;7|
000000d0  31 48 1b 5b 30 4b 1b 5b  31 37 3b 33 30 48 1b 5b  |1H.[0K.[17;30H.[|
000000e0  32 4b 1b 5b 32 34 3b 31  66 2a 1b 5b 31 3b 31 66  |2K.[24;1f*.[1;1f|
000000f0  2a 1b 5b 32 34 3b 32 66  2a 1b 5b 31 3b 32 66 2a  |*.[24;2f*.[1;2f*|
00000100  1b 5b 32 34 3b 33 66 2a  1b 5b 31 3b 33 66 2a 1b  |.[24;3f*.[1;3f*.|
00000110  5b 32 34 3b 34 66 2a 1b  5b 31 3b 34 66 2a 1b 5b  |[24;4f*.[1;4f*.[|
00000120  32 34 3b 35 66 2a 1b 5b  31 3b 35 66 2a 1b 5b 32  |24;5f*.[1;5f*.[2|
00000130  34 3b 36 66 2a 1b 5b 31  3b 36 66 2a 1b 5b 32 34  |4;6f*.[1;6f*.[24|
00000140  3b 37 66 2a 1b 5b 31 3b  37 66 2a 1b 5b 32 34 3b  |;7f*.[1;7f*.[24;|
00000150  38 66 2a 1b 5b 31 3b 38  66 2a 1b 5b 32 34 3b 39  |8f*.[1;8f*.[24;9|
00000160  66 2a 1b 5b 31 3b 39 66  2a 1b 5b 32 34 3b 31 30  |f*.[1;9f*.[24;10|
00000170  66 2a 1b 5b 31 3b 31 30  66 2a 1b 5b 32 34 3b 31  |f*.[1;10f*.[24;1|
00000180  31 66 2a 1b 5b 31 3b 31  31 66 2a 1b 5b 32 34 3b  |1f*.[1;11f*.[24;|
00000190  31 32 66 2a 1b 5b 31 3b  31 32 66 2a 1b 5b 32 34  |12f*.[1;12f*.[24|
000001a0  3b 31 33 66 2a 1b 5b 31  3b 31 33 66 2a 1b 5b 32  |;13f*.[1;13f*.[2|
000001b0  34 3b 31 34 66 2a 1b 5b  31 3b 31 34 66 2a 1b 5b  |4;14f*.[1;14f*.[|
000001c0  32 34 3b 31 35 66 2a 1b  5b 31 3b 31 35 66 2a 1b  |24;15f*.[1;15f*.|
000001d0  5b 32 34 3b 31 36 66 2a  1b 5b 31 3b 31 36 66 2a  |[24;16f*.[1;16f*|
000001e0  1b 5b 32 34 3b 31 37 66  2a 1b 5b 31 3b 31 37 66  |.[24;17f*.[1;17f|
000001f0  2a 1b 5b 32 34 3b 31 38  66 2a 1b 5b 31 3b 31 38  |*.[24;18f*.[1;18|
00000200  66 2a 1b 5b 32 34 3b 31  39 66 2a 1b 5b 31 3b 31  |f*.[24;19f*.[1;1|
00000210  39 66 2a 1b 5b 32 34 3b  32 30 66 2a 1b 5b 31 3b  |9f*.[24;20f*.[1;|
00000220  32 30 66 2a 1b 5b 32 34  3b 32 31 66 2a 1b 5b 31  |20f*.[24;21f*.[1|
00000230  3b 32 31 66 2a 1b 5b 32  34 3b 32 32 66 2a 1b 5b  |;21f*.[24;22f*.[|
00000240  31 3b 32 32 66 2a 1b 5b  32 34 3b 32 33 66 2a 1b  |1;22f*.[24;23f*.|
00000250  5b 31 3b 32 33 66 2a 1b  5b 32 34 3b 32 34 66 2a  |[1;23f*.[24;24f*|
00000260  1b 5b 31 3b 32 34 66 2a  1b 5b 32 34 3b 32 35 66  |.[1;24f*.[24;25f|
00000270  2a 1b 5b 31 3b 32 35 66  2a 1b 5b 32 34 3b 32 36  |*.[1;25f*.[24;26|
00000280  66 2a 1b 5b 31 3b 32 36  66 2a 1b 5b 32 34 3b 32  |f*.[1;26f*.[24;2|
00000290  37 66 2a 1b 5b 31 3b 32  37 66 2a 1b 5b 32 34 3b  |7f*.[1;27f*.[24;|
000002a0  32 38 66 2a 1b 5b 31 3b  32 38 66 2a 1b 5b 32 34  |28f*.[1;28f*.[24|
000002b0  3b 32 39 66 2a 1b 5b 31  3b 32 39 66 2a 1b 5b 32  |;29f*.[1;29f*.[2|
000002c0  34 3b 33 30 66 2a 1b 5b  31 3b 33 30 66 2a 1b 5b  |4;30f*.[1;30f*.[|
000002d0  32 34 3b 33 31 66 2a 1b  5b 31 3b 33 31 66 2a 1b  |24;31f*.[1;31f*.|
000002e0  5b 32 34 3b 33 32 66 2a  1b 5b 31 3b 33 32 66 2a  |[24;32f*.[1;32f*|
000002f0  1b 5b 32 34 3b 33 33 66                           |.[24;33f|
stdout:
00000000  2a 1b 5b 31 3b 33 33 66  2a 1b 5b 32 34 3b 33 34  |*.[1;33f*.[24;34|
00000010  66 2a 1b 5b 31 3b 33 34  66 2a 1b 5b 32 34 3b 33  |f*.[1;34f*.[24;3|
00000020  35 66 2a 1b 5b 31 3b 33  35 66 2a 1b 5b 32 34 3b  |5f*.[1;35f*.[24;|
00000030  33 36 66 2a 1b 5b 31 3b  33 36 66 2a 1b 5b 32 34  |36f*.[1;36f*.[24|
00000040  3b 33 37 66 2a 1b 5b 31  3b 33 37 66 2a 1b 5b 32  |;37f*.[1;37f*.[2|
00000050  34 3b 33 38 66 2a 1b 5b  31 3b 33 38 66 2a 1b 5b  |4;38f*.[1;38f*.[|
00000060  32 34 3b 33 39 66 2a 1b  5b 31 3b 33 39 66 2a 1b  |24;39f*.[1;39f*.|
00000070  5b 32 34 3b 34 30 66 2a  1b 5b 31 3b 34 30 66 2a  |[24;40f*.[1;40f*|
00000080  1b 5b 32 34 3b 34 31 66  2a 1b 5b 31 3b 34 31 66  |.[24;41f*.[1;41f|
00000090  2a 1b 5b 32 34 3b 34 32  66 2a 1b 5b 31 3b 34 32  |*.[24;42f*.[1;42|
000000a0  66 2a 1b 5b 32 34 3b 34  33 66 2a 1b 5b 31 3b 34  |f*.[24;43f*.[1;4|
000000b0  33 66 2a 1b 5b 32 34 3b  34 34 66 2a 1b 5b 31 3b  |3f*.[24;44f*.[1;|
000000c0  34 34 66 2a 1b 5b 32 34  3b 34 35 66 2a 1b 5b 31  |44f*.[24;45f*.[1|
000000d0  3b 34 35 66 2a 1b 5b 32  34 3b 34 36 66 2a 1b 5b  |;45f*.[24;46f*.[|
000000e0  31 3b 34 36 66 2a 1b 5b  32 34 3b 34 37 66 2a 1b  |1;46f*.[24;47f*.|
000000f0  5b 31 3b 34 37 66 2a 1b  5b 32 34 3b 34 38 66 2a  |[1;47f*.[24;48f*|
00000100  1b 5b 31 3b 34 38 66 2a  1b 5b 32 34 3b 34 39 66  |.[1;48f*.[24;49f|
00000110  2a 1b 5b 31 3b 34 39 66  2a 1b 5b 32 34 3b 35 30  |*.[1;49f*.[24;50|
00000120  66 2a 1b 5b 31 3b 35 30  66 2a 1b 5b 32 34 3b 35  |f*.[1;50f*.[24;5|
00000130  31 66 2a 1b 5b 31 3b 35  31 66 2a 1b 5b 32 34 3b  |1f*.[1;51f*.[24;|
00000140  35 32 66 2a 1b 5b 31 3b  35 32 66 2a 1b 5b 32 34  |52f*.[1;52f*.[24|
00000150  3b 35 33 66 2a 1b 5b 31  3b 35 33 66 2a 1b 5b 32  |;53f*.[1;53f*.[2|
00000160  34 3b 35 34 66 2a 1b 5b  31 3b 35 34 66 2a 1b 5b  |4;54f*.[1;54f*.[|
00000170  32 34 3b 35 35 66 2a 1b  5b 31 3b 35 35 66 2a 1b  |24;55f*.[1;55f*.|
00000180  5b 32 34 3b 35 36 66 2a  1b 5b 31 3b 35 36 66 2a  |[24;56f*.[1;56f*|
00000190  1b 5b 32 34 3b 35 37 66  2a 1b 5b 31 3b 35 37 66  |.[24;57f*.[1;57f|
000001a0  2a 1b 5b 32 34 3b 35 38  66 2a 1b 5b 31 3b 35 38  |*.[24;58f*.[1;58|
000001b0  66 2a 1b 5b 32 34 3b 35  39 66 2a 1b 5b 31 3b 35  |f*.[24;59f*.[1;5|
000001c0  39 66 2a 1b 5b 32 34 3b  36 30 66 2a 1b 5b 31 3b  |9f*.[24;60f*.[1;|
000001d0  36 30 66 2a 1b 5b 32 34  3b 36 31 66 2a 1b 5b 31  |60f*.[24;61f*.[1|
000001e0  3b 36 31 66 2a 1b 5b 32  34 3b 36 32 66 2a 1b 5b  |;61f*.[24;62f*.[|
000001f0  31 3b 36 32 66 2a 1b 5b  32 34 3b 36 33 66 2a 1b  |1;62f*.[24;63f*.|
00000200  5b 31 3b 36 33 66 2a 1b  5b 32 34 3b 36 34 66 2a  |[1;63f*.[24;64f*|
00000210  1b 5b 31 3b 36 34 66 2a  1b 5b 32 34 3b 36 35 66  |.[1;64f*.[24;65f|
00000220  2a 1b 5b 31 3b 36 35 66  2a 1b 5b 32 34 3b 36 36  |*.[1;65f*.[24;66|
00000230  66 2a 1b 5b 31 3b 36 36  66 2a 1b 5b 32 34 3b 36  |f*.[1;66f*.[24;6|
00000240  37 66 2a 1b 5b 31 3b 36  37 66 2a 1b 5b 32 34 3b  |7f*.[1;67f*.[24;|
00000250  36 38 66 2a 1b 5b 31 3b  36 38 66 2a 1b 5b 32 34  |68f*.[1;68f*.[24|
00000260  3b 36 39 66 2a 1b 5b 31  3b 36 39 66 2a 1b 5b 32  |;69f*.[1;69f*.[2|
00000270  34 3b 37 30 66 2a 1b 5b  31 3b 37 30 66 2a 1b 5b  |4;70f*.[1;70f*.[|
00000280  32 34 3b 37 31 66 2a 1b  5b 31 3b 37 31 66 2a 1b  |24;71f*.[1;71f*.|
00000290  5b 32 34 3b 37 32 66 2a  1b 5b 31 3b 37 32 66 2a  |[24;72f*.[1;72f*|
000002a0  1b 5b 32 34 3b 37 33 66  2a 1b 5b 31 3b 37 33 66  |.[24;73f*.[1;73f|
000002b0  2a 1b 5b 32 34 3b 37 34  66 2a 1b 5b 31 3b 37 34  |*.[24;74f*.[1;74|
000002c0  66 2a 1b 5b 32 34 3b 37  35 66 2a 1b 5b 31 3b 37  |f*.[24;75f*.[1;7|
000002d0  35 66 2a 1b 5b 32 34 3b  37 36 66 2a 1b 5b 31 3b  |5f*.[24;76f*.[1;|
000002e0  37 36 66 2a 1b 5b 32 34  3b 37 37 66 2a 1b 5b 31  |76f*.[24;77f*.[1|
000002f0  3b 37 37 66 2a 1b 5b 32  34 3b 37 38 66 2a 1b 5b  |;77f*.[24;78f*.[|
00000300  31 3b 37 38 66 2a 1b 5b  32 34 3b 37 39 66 2a 1b  |1;78f*.[24;79f*.|
00000310  5b 31 3b 37 39 66 2a 1b  5b 32 34 3b 38 30 66 2a  |[1;79f*.[24;80f*|
00000320  1b 5b 31 3b 38 30 66 2a  1b 5b 32 3b 32 48 2b 1b  |.[1;80f*.[2;2H+.|
00000330  5b 31 44 1b 44 2b 1b 5b  31 44 1b 44 2b 1b 5b 31  |[1D.D+.[1D.D+.[1|
00000340  44 1b 44 2b 1b 5b 31 44  1b 44 2b 1b 5b 31 44 1b  |D.D+.[1D.D+.[1D.|
00000350  44 2b 1b 5b 31 44 1b 44  2b 1b 5b 31 44 1b 44 2b  |D+.[1D.D+.[1D.D+|
00000360  1b 5b 31 44 1b 44 2b 1b  5b 31 44 1b 44 2b 1b 5b  |.[1D.D+.[1D.D+.[|
00000370  31 44 1b 44 2b 1b 5b 31  44 1b 44 2b 1b 5b 31 44  |1D.D+.[1D.D+.[1D|
00000380  1b 44 2b 1b 5b 31 44 1b  44 2b 1b 5b 31 44 1b 44  |.D+.[1D.D+.[1D.D|
00000390  2b 1b 5b 31 44 1b 44 2b  1b 5b 31 44 1b 44 2b 1b  |+.[1D.D+.[1D.D+.|
000003a0  5b 31 44 1b 44 2b 1b 5b  31 44 1b 44 2b 1b 5b 31  |[1D.D+.[1D.D+.[1|
000003b0  44 1b 44 2b 1b 5b 31 44  1b 44 2b 1b 5b 31 44 1b  |D.D+.[1D.D+.[1D.|
000003c0  44 2b 1b 5b 31 44 1b 44  1b 5b 32 33 3b 37 39 48  |D+.[1D.D.[23;79H|
000003d0  2b 1b 5b 31 44 1b 4d 2b  1b 5b 31 44 1b 4d 2b 1b  |+.[1D.M+.[1D.M+.|
000003e0  5b 31 44 1b 4d 2b 1b 5b  31 44 1b 4d 2b 1b 5b 31  |[1D.M+.[1D.M+.[1|
000003f0  44 1b 4d 2b 1b 5b 31 44  1b 4d 2b 1b 5b 31 44 1b  |D.M+.[1D.M+.[1D.|
00000400  4d 2b 1b 5b 31 44 1b 4d  2b 1b 5b 31 44 1b 4d 2b  |M+.[1D.M+.[1D.M+|
00000410  1b 5b 31 44 1b 4d 2b 1b  5b 31 44 1b 4d 2b 1b 5b  |.[1D.M+.[1D.M+.[|
00000420  31 44 1b 4d 2b 1b 5b 31  44 1b 4d 2b 1b 5b 31 44  |1D.M+.[1D.M+.[1D|
00000430  1b 4d 2b 1b 5b 31 44 1b  4d 2b 1b 5b 31 44 1b 4d  |.M+.[1D.M+.[1D.M|
00000440  2b 1b 5b 31 44 1b 4d 2b  1b 5b 31 44 1b 4d 2b 1b  |+.[1D.M+.[1D.M+.|
00000450  5b 31 44 1b 4d 2b 1b 5b  31 44 1b 4d 2b 1b 5b 31  |[1D.M+.[1D.M+.[1|
00000460  44 1b 4d 2b 1b 5b 31 44  1b 4d 1b 5b 32 3b 31 48  |D.M+.[1D.M.[2;1H|
00000470  2a 1b 5b 32 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[2;80H*.[10D.E|
00000480  2a 1b 5b 33 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[3;80H*.[10D.E|
00000490  2a 1b 5b 34 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[4;80H*.[10D.E|
000004a0  2a 1b 5b 35 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[5;80H*.[10D.E|
000004b0  2a 1b 5b 36 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[6;80H*.[10D.E|
000004c0  2a 1b 5b 37 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[7;80H*.[10D.E|
000004d0  2a 1b 5b 38 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[8;80H*.[10D.E|
000004e0  2a 1b 5b 39 3b 38 30 48  2a 1b 5b 31 30 44 1b 45  |*.[9;80H*.[10D.E|
000004f0  2a 1b 5b 31 30 3b 38 30  48 2a 1b 5b 31 30 44 0d  |*.[10;80H*.[10D.|
00000500  0a 2a 1b 5b 31 31 3b 38  30 48 2a 1b 5b 31 30 44  |.*.[11;80H*.[10D|
00000510  0d 0a 2a 1b 5b 31 32 3b  38 30 48 2a 1b 5b 31 30  |..*.[12;80H*.[10|
00000520  44 0d 0a 2a 1b 5b 31 33  3b 38 30 48 2a 1b 5b 31  |D..*.[13;80H*.[1|
00000530  30 44 0d 0a 2a 1b 5b 31  34 3b 38 30 48 2a 1b 5b  |0D..*.[14;80H*.[|
00000540  31 30 44 0d 0a 2a 1b 5b  31 35 3b 38 30 48 2a 1b  |10D..*.[15;80H*.|
00000550  5b 31 30 44 0d 0a 2a 1b  5b 31 36 3b 38 30 48 2a  |[10D..*.[16;80H*|
00000560  1b 5b 31 30 44 0d 0a 2a  1b 5b 31 37 3b 38 30 48  |.[10D..*.[17;80H|
00000570  2a 1b 5b 31 30 44 0d 0a  2a 1b 5b 31 38 3b 38 30  |*.[10D..*.[18;80|
00000580  48 2a 1b 5b 31 30 44 0d  0a 2a 1b 5b 31 39 3b 38  |H*.[10D..*.[19;8|
00000590  30 48 2a 1b 5b 31 30 44  0d 0a 2a 1b 5b 32 30 3b  |0H*.[10D..*.[20;|
000005a0  38 30 48 2a 1b 5b 31 30  44 0d 0a 2a 1b 5b 32 31  |80H*.[10D..*.[21|
000005b0  3b 38 30 48 2a 1b 5b 31  30 44 0d 0a 2a 1b 5b 32  |;80H*.[10D..*.[2|
000005c0  32 3b 38 30 48 2a 1b 5b  31 30 44 0d 0a 2a 1b 5b  |2;80H*.[10D..*.[|
000005d0  32 33 3b 38 30 48 2a 1b  5b 31 30 44 0d 0a 1b 5b  |23;80H*.[10D...[|
000005e0  32 3b 31 30 48 1b 5b 34  32 44 1b 5b 32 43 2b 1b  |2;10H.[42D.[2C+.|
000005f0  5b 30 43 1b 5b 32 44 1b  5b 31 43 2b 1b 5b 30 43  |[0C.[2D.[1C+.[0C|
00000600  1b 5b 32 44 1b 5b 31 43  2b 1b 5b 30 43 1b 5b 32  |.[2D.[1C+.[0C.[2|
00000610  44 1b 5b 31 43 2b 1b 5b  30 43 1b 5b 32 44 1b 5b  |D.[1C+.[0C.[2D.[|
00000620  31 43 2b 1b 5b 30 43 1b  5b 32 44 1b 5b 31 43 2b  |1C+.[0C.[2D.[1C+|
00000630  1b 5b 30 43 1b 5b 32 44  1b 5b 31 43 2b 1b 5b 30  |.[0C.[2D.[1C+.[0|
00000640  43 1b 5b 32 44 1b 5b 31  43 2b 1b 5b 30 43 1b 5b  |C.[2D.[1C+.[0C.[|
00000650  32 44 1b 5b 31 43 2b 1b  5b 30 43 1b 5b 32 44 1b  |2D.[1C+.[0C.[2D.|
00000660  5b 31 43 2b 1b 5b 30 43  1b 5b 32 44 1b 5b 31 43  |[1C+.[0C.[2D.[1C|
00000670  2b 1b 5b 30 43 1b 5b 32  44 1b 5b 31 43 2b 1b 5b  |+.[0C.[2D.[1C+.[|
00000680  30 43 1b 5b 32 44 1b 5b  31 43 2b 1b 5b 30 43 1b  |0C.[2D.[1C+.[0C.|
00000690  5b 32 44 1b 5b 31 43 2b  1b 5b 30 43 1b 5b 32 44  |[2D.[1C+.[0C.[2D|
000006a0  1b 5b 31 43 2b 1b 5b 30  43 1b 5b 32 44 1b 5b 31  |.[1C+.[0C.[2D.[1|
000006b0  43 2b 1b 5b 30 43 1b 5b  32 44 1b 5b 31 43 2b 1b  |C+.[0C.[2D.[1C+.|
000006c0  5b 30 43 1b 5b 32 44 1b  5b 31 43 2b 1b 5b 30 43  |[0C.[2D.[1C+.[0C|
000006d0  1b 5b 32 44 1b 5b 31 43  2b 1b 5b 30 43 1b 5b 32  |.[2D.[1C+.[0C.[2|
000006e0  44 1b 5b 31 43 2b 1b 5b  30 43 1b 5b 32 44 1b 5b  |D.[1C+.[0C.[2D.[|
000006f0  31 43 2b 1b 5b 30 43 1b  5b 32 44 1b 5b 31 43 2b  |1C+.[0C.[2D.[1C+|
00000700  1b 5b 30 43 1b 5b 32 44  1b 5b 31 43 2b 1b 5b 30  |.[0C.[2D.[1C+.[0|
00000710  43 1b 5b 32 44 1b 5b 31  43 2b 1b 5b 30 43 1b 5b  |C.[2D.[1C+.[0C.[|
00000720  32 44 1b 5b 31 43 2b 1b  5b 30 43 1b 5b 32 44 1b  |2D.[1C+.[0C.[2D.|
00000730  5b 31 43 2b 1b 5b 30 43  1b 5b 32 44 1b 5b 31 43  |[1C+.[0C.[2D.[1C|
00000740  2b 1b 5b 30 43 1b 5b 32  44 1b 5b 31 43 2b 1b 5b  |+.[0C.[2D.[1C+.[|
00000750  30 43 1b 5b 32 44 1b 5b  31 43 2b 1b 5b 30 43 1b  |0C.[2D.[1C+.[0C.|
00000760  5b 32 44 1b 5b 31 43 2b  1b 5b 30 43 1b 5b 32 44  |[2D.[1C+.[0C.[2D|
00000770  1b 5b 31 43 2b 1b 5b 30  43 1b 5b 32 44 1b 5b 31  |.[1C+.[0C.[2D.[1|
00000780  43 2b 1b 5b 30 43 1b 5b  32 44 1b 5b 31 43 2b 1b  |C+.[0C.[2D.[1C+.|
00000790  5b 30 43 1b 5b 32 44 1b  5b 31 43 2b 1b 5b 30 43  |[0C.[2D.[1C+.[0C|
000007a0  1b 5b 32 44 1b 5b 31 43  2b 1b 5b 30 43 1b 5b 32  |.[2D.[1C+.[0C.[2|
000007b0  44 1b 5b 31 43 2b 1b 5b  30 43 1b 5b 32 44 1b 5b  |D.[1C+.[0C.[2D.[|
000007c0  31 43 2b 1b 5b 30 43 1b  5b 32 44 1b 5b 31 43 2b  |1C+.[0C.[2D.[1C+|
000007d0  1b 5b 30 43 1b 5b 32 44  1b 5b 31 43 2b 1b 5b 30  |.[0C.[2D.[1C+.[0|
000007e0  43 1b 5b 32 44 1b 5b 31  43 2b 1b 5b 30 43 1b 5b  |C.[2D.[1C+.[0C.[|
000007f0  32 44 1b 5b 31 43 2b 1b  5b 30 43 1b 5b 32 44 1b  |2D.[1C+.[0C.[2D.|
00000800  5b 31 43 2b 1b 5b 30 43  1b 5b 32 44 1b 5b 31 43  |[1C+.[0C.[2D.[1C|
00000810  2b 1b 5b 30 43 1b 5b 32  44 1b 5b 31 43 2b 1b 5b  |+.[0C.[2D.[1C+.[|
00000820  30 43 1b 5b 32 44 1b 5b  31 43 2b 1b 5b 30 43 1b  |0C.[2D.[1C+.[0C.|
00000830  5b 32 44 1b 5b 31 43 2b  1b 5b 30 43 1b 5b 32 44  |[2D.[1C+.[0C.[2D|
00000840  1b 5b 31 43 2b 1b 5b 30  43 1b 5b 32 44 1b 5b 31  |.[1C+.[0C.[2D.[1|
00000850  43 2b 1b 5b 30 43 1b 5b  32 44 1b 5b 31 43 2b 1b  |C+.[0C.[2D.[1C+.|
00000860  5b 30 43 1b 5b 32 44 1b  5b 31 43 2b 1b 5b 30 43  |[0C.[2D.[1C+.[0C|
00000870  1b 5b 32 44 1b 5b 31 43  2b 1b 5b 30 43 1b 5b 32  |.[2D.[1C+.[0C.[2|
00000880  44 1b 5b 31 43 2b 1b 5b  30 43 1b 5b 32 44 1b 5b  |D.[1C+.[0C.[2D.[|
00000890  31 43 2b 1b 5b 30 43 1b  5b 32 44 1b 5b 31 43 2b  |1C+.[0C.[2D.[1C+|
000008a0  1b 5b 30 43 1b 5b 32 44  1b 5b 31 43 2b 1b 5b 30  |.[0C.[2D.[1C+.[0|
000008b0  43 1b 5b 32 44 1b 5b 31  43 2b 1b 5b 30 43 1b 5b  |C.[2D.[1C+.[0C.[|
000008c0  32 44 1b 5b 31 43 2b 1b  5b 30 43 1b 5b 32 44 1b  |2D.[1C+.[0C.[2D.|
000008d0  5b 31 43 2b 1b 5b 30 43  1b 5b 32 44 1b 5b 31 43  |[1C+.[0C.[2D.[1C|
000008e0  2b 1b 5b 30 43 1b 5b 32  44 1b 5b 31 43 2b 1b 5b  |+.[0C.[2D.[1C+.[|
000008f0  30 43 1b 5b 32 44 1b 5b  31 43 2b 1b 5b 30 43 1b  |0C.[2D.[1C+.[0C.|
00000900  5b 32 44 1b 5b 31 43 2b  1b 5b 30 43 1b 5b 32 44  |[2D.[1C+.[0C.[2D|
00000910  1b 5b 31 43 2b 1b 5b 30  43 1b 5b 32 44 1b 5b 31  |.[1C+.[0C.[2D.[1|
00000920  43 2b 1b 5b 30 43 1b 5b  32 44 1b 5b 31 43 2b 1b  |C+.[0C.[2D.[1C+.|
00000930  5b 30 43 1b 5b 32 44 1b  5b 31 43 2b 1b 5b 30 43  |[0C.[2D.[1C+.[0C|
00000940  1b 5b 32 44 1b 5b 31 43  2b 1b 5b 30 43 1b 5b 32  |.[2D.[1C+.[0C.[2|
00000950  44 1b 5b 31 43 2b 1b 5b  30 43 1b 5b 32 44 1b 5b  |D.[1C+.[0C.[2D.[|
00000960  31 43 2b 1b 5b 30 43 1b  5b 32 44 1b 5b 31 43 2b  |1C+.[0C.[2D.[1C+|
00000970  1b 5b 30 43 1b 5b 32 44  1b 5b 31 43 2b 1b 5b 30  |.[0C.[2D.[1C+.[0|
00000980  43 1b 5b 32 44 1b 5b 31  43 2b 1b 5b 30 43 1b 5b  |C.[2D.[1C+.[0C.[|
00000990  32 44 1b 5b 31 43 2b 1b  5b 30 43 1b 5b 32 44 1b  |2D.[1C+.[0C.[2D.|
000009a0  5b 31 43 2b 1b 5b 30 43  1b 5b 32 44 1b 5b 31 43  |[1C+.[0C.[2D.[1C|
000009b0  2b 1b 5b 30 43 1b 5b 32  44 1b 5b 31 43 2b 1b 5b  |+.[0C.[2D.[1C+.[|
000009c0  30 43 1b 5b 32 44 1b 5b  31 43 1b 5b 32 33 3b 37  |0C.[2D.[1C.[23;7|
000009d0  30 48 1b 5b 34 32 43 1b  5b 32 44 2b 1b 5b 31 44  |0H.[42C.[2D+.[1D|
000009e0  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
000009f0  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000a00  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000a10  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000a20  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000a30  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000a40  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000a50  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000a60  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000a70  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000a80  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000a90  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000aa0  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000ab0  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000ac0  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000ad0  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000ae0  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000af0  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000b00  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000b10  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000b20  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000b30  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000b40  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000b50  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000b60  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000b70  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000b80  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000b90  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000ba0  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000bb0  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000bc0  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000bd0  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000be0  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000bf0  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000c00  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000c10  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000c20  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000c30  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000c40  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000c50  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000c60  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000c70  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000c80  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000c90  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000ca0  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000cb0  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000cc0  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000cd0  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000ce0  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000cf0  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000d00  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000d10  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000d20  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000d30  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000d40  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000d50  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000d60  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000d70  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000d80  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000d90  30 44 08 2b 1b 5b 31 44  1b 5b 31 43 1b 5b 30 44  |0D.+.[1D.[1C.[0D|
00000da0  08 2b 1b 5b 31 44 1b 5b  31 43 1b 5b 30 44 08 2b  |.+.[1D.[1C.[0D.+|
00000db0  1b 5b 31 44 1b 5b 31 43  1b 5b 30 44 08 2b 1b 5b  |.[1D.[1C.[0D.+.[|
00000dc0  31 44 1b 5b 31 43 1b 5b  30 44 08 2b 1b 5b 31 44  |1D.[1C.[0D.+.[1D|
00000dd0  1b 5b 31 43 1b 5b 30 44  08 2b 1b 5b 31 44 1b 5b  |.[1C.[0D.+.[1D.[|
00000de0  31 43 1b 5b 30 44 08 2b  1b 5b 31 44 1b 5b 31 43  |1C.[0D.+.[1D.[1C|
00000df0  1b 5b 30 44 08 2b 1b 5b  31 44 1b 5b 31 43 1b 5b  |.[0D.+.[1D.[1C.[|
00000e00  30 44 08 1b 5b 31 3b 31  48 1b 5b 31 30 41 1b 5b  |0D..[1;1H.[10A.[|
00000e10  31 41 1b 5b 30 41 1b 5b  32 34 3b 38 30 48 1b 5b  |1A.[0A.[24;80H.[|
00000e20  31 30 42 1b 5b 31 42 1b  5b 30 42 1b 5b 31 30 3b  |10B.[1B.[0B.[10;|
00000e30  31 32 48 20 20 20 20 20  20 20 20 20 20 20 20 20  |12H             |
00000e40  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000e50  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000e60  20 20 20 20 20 20 20 20  20 20 20 20 20 1b 5b 31  |             .[1|
00000e70  42 1b 5b 35 38 44 20 20  20 20 20 20 20 20 20 20  |B.[58D          |
00000e80  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000e90  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000ea0  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000eb0  1b 5b 31 42 1b 5b 35 38  44 20 20 20 20 20 20 20  |.[1B.[58D       |
00000ec0  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000ed0  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000ee0  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000ef0  20 20 20 1b 5b 31 42 1b  5b 35 38 44 20 20 20 20  |   .[1B.[58D    |
00000f00  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000f10  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000f20  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000f30  20 20 20 20 20 20 1b 5b  31 42 1b 5b 35 38 44 20  |      .[1B.[58D |
00000f40  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000f50  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000f60  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000f70  20 20 20 20 20 20 20 20  20 1b 5b 31 42 1b 5b 35  |         .[1B.[5|
00000f80  38 44 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |8D              |
00000f90  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000fa0  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
00000fb0  20 20 20 20 20 20 20 20  20 20 20 20 1b 5b 31 42  |            .[1B|
00000fc0  1b 5b 35 38 44 1b 5b 35  41 1b 5b 31 43 54 68 65  |.[58D.[5A.[1CThe|
00000fd0  20 73 63 72 65 65 6e 20  73 68 6f 75 6c 64 20 62  | screen should b|
00000fe0  65 20 63 6c 65 61 72 65  64 2c 20 20 61 6e 64 20  |e cleared,  and |
00000ff0  68 61 76 65 20 61 6e 20  75 6e 62 72 6f 6b 65 6e  |have an unbroken|
00001000  20 62 6f 72 2d 1b 5b 31  32 3b 31 33 48 64 65 72  | bor-.[12;13Hder|
00001010  20 6f 66 20 2a 27 73 20  61 6e 64 20 2b 27 73 20  | of *'s and +'s |
00001020  61 72 6f 75 6e 64 20 74  68 65 20 65 64 67 65 2c  |around the edge,|
00001030  20 20 20 61 6e 64 20 65  78 61 63 74 6c 79 20 69  |   and exactly i|
00001040  6e 20 74 68 65 1b 5b 31  33 3b 31 33 48 6d 69 64  |n the.[13;13Hmid|
00001050  64 6c 65 20 20 74 68 65  72 65 20 73 68 6f 75 6c  |dle  there shoul|
00001060  64 20 62 65 20 61 20 66  72 61 6d 65 20 6f 66 20  |d be a frame of |
00001070  45 27 73 20 61 72 6f 75  6e 64 20 74 68 69 73 20  |E's around this |
00001080  20 74 65 78 74 1b 5b 31  34 3b 31 33 48 77 69 74  | text.[14;13Hwit|
00001090  68 20 20 6f 6e 65 20 28  31 29 20 66 72 65 65 20  |h  one (1) free |
000010a0  70 6f 73 69 74 69 6f 6e  20 61 72 6f 75 6e 64 20  |position around |
000010b0  69 74 2e 20 20 20 20 50  75 73 68 20 3c 52 45 54  |it.    Push <RET|
000010c0  55 52 4e 3e                                       |URN>|
`,
		output: `
********************************************************************************
*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+        EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE        +*
*+        E                                                          E        +*
*+        E The screen should be cleared,  and have an unbroken bor- E        +*
*+        E der of *'s and +'s around the edge,   and exactly in the E        +*
*+        E middle  there should be a frame of E's around this  text E        +*
*+        E with  one (1) free position around it.    Push <RETURN>  E        +*
*+        E                                                          E        +*
*+        EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE        +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*+                                                                            +*
*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*
********************************************************************************
`,
	},
	{
		input: `stdout:
00000000  1b 5b 3f 33 6c 1b 5b 3f  33 6c 54 65 73 74 20 6f  |.[?3l.[?3lTest o|
00000010  66 20 61 75 74 6f 77 72  61 70 2c 20 6d 69 78 69  |f autowrap, mixi|
00000020  6e 67 20 63 6f 6e 74 72  6f 6c 20 61 6e 64 20 70  |ng control and p|
00000030  72 69 6e 74 20 63 68 61  72 61 63 74 65 72 73 2e  |rint characters.|
00000040  0d 0d 0a 54 68 65 20 6c  65 66 74 2f 72 69 67 68  |...The left/righ|
00000050  74 20 6d 61 72 67 69 6e  73 20 73 68 6f 75 6c 64  |t margins should|
00000060  20 68 61 76 65 20 6c 65  74 74 65 72 73 20 69 6e  | have letters in|
00000070  20 6f 72 64 65 72 3a 0d  0d 0a 1b 5b 33 3b 32 31  | order:....[3;21|
00000080  72 1b 5b 3f 36 68 1b 5b  31 39 3b 31 48 41 1b 5b  |r.[?6h.[19;1HA.[|
00000090  31 39 3b 38 30 48 61 0d  0a 1b 5b 31 38 3b 38 30  |19;80Ha...[18;80|
000000a0  48 61 42 1b 5b 31 39 3b  38 30 48 42 08 20 62 0d  |HaB.[19;80HB. b.|
000000b0  0a 1b 5b 31 39 3b 38 30  48 43 08 08 09 09 63 1b  |..[19;80HC....c.|
000000c0  5b 31 39 3b 32 48 08 43  0d 0a 1b 5b 31 39 3b 38  |[19;2H.C...[19;8|
000000d0  30 48 0d 0a 1b 5b 31 38  3b 31 48 44 1b 5b 31 38  |0H...[18;1HD.[18|
000000e0  3b 38 30 48 64 1b 5b 31  39 3b 31 48 45 1b 5b 31  |;80Hd.[19;1HE.[1|
000000f0  39 3b 38 30 48 65 0d 0a  1b 5b 31 38 3b 38 30 48  |9;80He...[18;80H|
00000100  65 46 1b 5b 31 39 3b 38  30 48 46 08 20 66 0d 0a  |eF.[19;80HF. f..|
00000110  1b 5b 31 39 3b 38 30 48  47 08 08 09 09 67 1b 5b  |.[19;80HG....g.[|
00000120  31 39 3b 32 48 08 47 0d  0a 1b 5b 31 39 3b 38 30  |19;2H.G...[19;80|
00000130  48 0d 0a 1b 5b 31 38 3b  31 48 48 1b 5b 31 38 3b  |H...[18;1HH.[18;|
00000140  38 30 48 68 1b 5b 31 39  3b 31 48 49 1b 5b 31 39  |80Hh.[19;1HI.[19|
00000150  3b 38 30 48 69 0d 0a 1b  5b 31 38 3b 38 30 48 69  |;80Hi...[18;80Hi|
00000160  4a 1b 5b 31 39 3b 38 30  48 4a 08 20 6a 0d 0a 1b  |J.[19;80HJ. j...|
00000170  5b 31 39 3b 38 30 48 4b  08 08 09 09 6b 1b 5b 31  |[19;80HK....k.[1|
00000180  39 3b 32 48 08 4b 0d 0a  1b 5b 31 39 3b 38 30 48  |9;2H.K...[19;80H|
00000190  0d 0a 1b 5b 31 38 3b 31  48 4c 1b 5b 31 38 3b 38  |...[18;1HL.[18;8|
000001a0  30 48 6c 1b 5b 31 39 3b  31 48 4d 1b 5b 31 39 3b  |0Hl.[19;1HM.[19;|
000001b0  38 30 48 6d 0d 0a 1b 5b  31 38 3b 38 30 48 6d 4e  |80Hm...[18;80HmN|
000001c0  1b 5b 31 39 3b 38 30 48  4e 08 20 6e 0d 0a 1b 5b  |.[19;80HN. n...[|
000001d0  31 39 3b 38 30 48 4f 08  08 09 09 6f 1b 5b 31 39  |19;80HO....o.[19|
000001e0  3b 32 48 08 4f 0d 0a 1b  5b 31 39 3b 38 30 48 0d  |;2H.O...[19;80H.|
000001f0  0a 1b 5b 31 38 3b 31 48  50 1b 5b 31 38 3b 38 30  |..[18;1HP.[18;80|
00000200  48 70 1b 5b 31 39 3b 31  48 51 1b 5b 31 39 3b 38  |Hp.[19;1HQ.[19;8|
00000210  30 48 71 0d 0a 1b 5b 31  38 3b 38 30 48 71 52 1b  |0Hq...[18;80HqR.|
00000220  5b 31 39 3b 38 30 48 52  08 20 72 0d 0a 1b 5b 31  |[19;80HR. r...[1|
00000230  39 3b 38 30 48 53 08 08  09 09 73 1b 5b 31 39 3b  |9;80HS....s.[19;|
00000240  32 48 08 53 0d 0a 1b 5b  31 39 3b 38 30 48 0d 0a  |2H.S...[19;80H..|
00000250  1b 5b 31 38 3b 31 48                              |.[18;1H|

init.js:86 syscall: {cmd: "write", fd: 1, data: Uint8Array(599), offset: 0, length: 599, …}
wasm_exec.js:399 actCSI: unsupported: ESC[3;21r (0x72)

wasm_exec.js:399 Unsupported ESC[?6h
process.js?_ts=1615105574977:119 process: {cmd: "result", id: 229, error: null, code: 599, buf: undefined, …}
init.js:86 syscall: {fd: 3, data: Uint8Array(36), offset: 0, length: 36, cmd: "write", …}
process.js?_ts=1615105574977:119 process: {cmd: "result", id: 230, error: null, code: 36, buf: undefined, …}
wasm_exec.js:399 stdout:
00000000  54 1b 5b 31 38 3b 38 30  48 74 1b 5b 31 39 3b 31  |T.[18;80Ht.[19;1|
00000010  48 55 1b 5b 31 39 3b 38  30 48 75 0d 0a 1b 5b 31  |HU.[19;80Hu...[1|
00000020  38 3b 38 30 48 75 56 1b  5b 31 39 3b 38 30 48 56  |8;80HuV.[19;80HV|
00000030  08 20 76 0d 0a 1b 5b 31  39 3b 38 30 48 57 08 08  |. v...[19;80HW..|
00000040  09 09 77 1b 5b 31 39 3b  32 48 08 57 0d 0a 1b 5b  |..w.[19;2H.W...[|
00000050  31 39 3b 38 30 48 0d 0a  1b 5b 31 38 3b 31 48 58  |19;80H...[18;1HX|
00000060  1b 5b 31 38 3b 38 30 48  78 1b 5b 31 39 3b 31 48  |.[18;80Hx.[19;1H|
00000070  59 1b 5b 31 39 3b 38 30  48 79 0d 0a 1b 5b 31 38  |Y.[19;80Hy...[18|
00000080  3b 38 30 48 79 5a 1b 5b  31 39 3b 38 30 48 5a 08  |;80HyZ.[19;80HZ.|
00000090  20 7a 0d 0a 1b 5b 3f 36  6c 1b 5b 72 1b 5b 32 32  | z...[?6l.[r.[22|
000000a0  3b 31 48 50 75 73 68 20  3c 52 45 54 55 52 4e 3e  |;1HPush <RETURN>|
`,
	},
	{
		input: `stdout:
00000000  1b 5b 3f 33 6c 1b 5b 32  4a 1b 5b 31 3b 31 48 54  |.[?3l.[2J.[1;1HT|
00000010  65 73 74 20 6f 66 20 63  75 72 73 6f 72 2d 63 6f  |est of cursor-co|
00000020  6e 74 72 6f 6c 20 63 68  61 72 61 63 74 65 72 73  |ntrol characters|
00000030  20 69 6e 73 69 64 65 20  45 53 43 20 73 65 71 75  | inside ESC sequ|
00000040  65 6e 63 65 73 2e 0d 0d  0a 42 65 6c 6f 77 20 73  |ences....Below s|
00000050  68 6f 75 6c 64 20 62 65  20 66 6f 75 72 20 69 64  |hould be four id|
00000060  65 6e 74 69 63 61 6c 20  6c 69 6e 65 73 3a 0d 0d  |entical lines:..|
00000070  0a 0d 0d 0a 41 20 42 20  43 20 44 20 45 20 46 20  |....A B C D E F |
00000080  47 20 48 20 49 0d 0d 0a  41 1b 5b 32 08 43 42 1b  |G H I...A.[2.CB.|
00000090  5b 32 08 43 43 1b 5b 32  08 43 44 1b 5b 32 08 43  |[2.CC.[2.CD.[2.C|
000000a0  45 1b 5b 32 08 43 46 1b  5b 32 08 43 47 1b 5b 32  |E.[2.CF.[2.CG.[2|
000000b0  08 43 48 1b 5b 32 08 43  49 1b 5b 32 08 43 0d 0d  |.CH.[2.CI.[2.C..|
000000c0  0a                                                |.|

init.js:86 syscall: {cmd: "write", fd: 1, data: Uint8Array(193), offset: 0, length: 193, …}
process.js?_ts=1615111086901:119 process: {cmd: "result", id: 181, error: null, code: 193, buf: undefined, …}
init.js:86 syscall: {length: 36, cmd: "write", fd: 3, data: Uint8Array(36), offset: 0, …}
process.js?_ts=1615111086901:119 process: {cmd: "result", id: 182, error: null, code: 36, buf: undefined, …}
wasm_exec.js:399 stdout:
00000000  41 20 1b 5b 0d 32 43 42  1b 5b 0d 34 43 43 1b 5b  |A .[.2CB.[.4CC.[|
00000010  0d 36 43 44 1b 5b 0d 38  43 45 1b 5b 0d 31 30 43  |.6CD.[.8CE.[.10C|
00000020  46 1b 5b 0d 31 32 43 47  1b 5b 0d 31 34 43 48 1b  |F.[.12CG.[.14CH.|
00000030  5b 0d 31 36 43 49 0d 0d  0a 1b 5b 32 30 6c 41 20  |[.16CI....[20lA |
00000040  1b 5b 31 0b 41 42 20 1b  5b 31 0b 41 43 20 1b 5b  |.[1.AB .[1.AC .[|
00000050  31 0b 41 44 20 1b 5b 31  0b 41 45 20 1b 5b 31 0b  |1.AD .[1.AE .[1.|
00000060  41 46 20 1b 5b 31 0b 41  47 20 1b 5b 31 0b 41 48  |AF .[1.AG .[1.AH|
00000070  20 1b 5b 31 0b 41 49 20  1b 5b 31 0b 41 0d 0d 0a  | .[1.AI .[1.A...|
00000080  0d 0d 0a 50 75 73 68 20  3c 52 45 54 55 52 4e 3e  |...Push <RETURN>|
`,
		output: `
Test of cursor-control characters inside ESC sequences.
Below should be four identical lines:

A B C D E F G H I
A B C D E F G H I
A B C D E F G H I
A B C D E F G H I

Push <RETURN>
`,
	},
}

func TestEmul(t *testing.T) {
	for idx, test := range emulTests {
		data, err := ParseHexDump([]byte(test.input))
		if err != nil {
			t.Errorf("TestEmul %d: failed to parse input: %s", idx, err)
			continue
		}
		lines, err := Trim(string(data))
		if err != nil {
			t.Errorf("TestEmul %d: Trim failed: %s", idx, err)
			continue
		}
		expected := strings.Split(test.output, "\n")
		if len(expected) > 0 && len(expected[0]) == 0 {
			expected = expected[1:]
		}
		if len(expected) > 0 && len(expected[len(expected)-1]) == 0 {
			expected = expected[:len(expected)-1]
		}
		if len(expected) == 0 {
			for _, l := range lines {
				fmt.Println(l)
			}
		} else {
			for i, l := range lines {
				if i >= len(expected) ||
					strings.TrimRight(l, " ") != expected[i] {
					t.Errorf("TestEmul %d: line %d differs:\n%s\n%s\n",
						idx, i, l, expected[i])
				}
			}
		}
	}
}
