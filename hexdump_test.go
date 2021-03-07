//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var hexDumpReaderTests = []string{
	"",
	"Hello, world!",
	`This is a longer input
with multiple lines
and still more lines
to test!`,
}

func TestHexDumpReader(t *testing.T) {
	for _, input := range hexDumpReaderTests {
		ibuf := []byte(input)

		obuf, err := ParseHexDump([]byte(hex.Dump(ibuf)))
		if err != nil {
			t.Errorf("ParseHexDump failed: %s", err)
			continue
		}
		if bytes.Compare(ibuf, obuf) != 0 {
			t.Errorf("Got invalid result: input:\n%soutput:\n%s",
				hex.Dump(ibuf), hex.Dump(obuf))
		}
	}
}
