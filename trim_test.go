//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"embed"
	"fmt"
	"io"
	"path"
	"strings"
	"testing"
)

//go:embed testsuite/*.hex testsuite/*.txt
var testsuite embed.FS

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

func TestEmul(t *testing.T) {
	const dir = "testsuite"
	entries, err := testsuite.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read testsuit: %s", err)
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".hex") {
			continue
		}
		err = emulTest(t, dir, entry.Name()[0:len(entry.Name())-4])
		if err != nil {
			t.Errorf("Test '%s' failed: %s", entry.Name(), err)
		}
	}
}

func emulTest(t *testing.T, dir, name string) error {
	hexf, err := testsuite.Open(path.Join(dir, name+".hex"))
	if err != nil {
		return err
	}
	defer hexf.Close()
	input, err := io.ReadAll(hexf)
	if err != nil {
		return err
	}
	data, err := ParseHexDump(input)
	if err != nil {
		return err
	}
	lines, err := Trim(string(data))
	if err != nil {
		return err
	}

	txtf, err := testsuite.Open(path.Join(dir, name+".txt"))
	if err != nil {
		fmt.Printf("TestEmul %s: not output file: %s\n------ output ------\n",
			name, err)
		for _, l := range lines {
			fmt.Println(l)
		}
		return nil
	}
	defer txtf.Close()
	output, err := io.ReadAll(txtf)
	if err != nil {
		return err
	}
	expected := strings.Split(string(output), "\n")
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
				t.Errorf("TestEmul %s: line %d differs:\n%s\n%s\n",
					name, i, l, expected[i])
			}
		}
	}
	return nil
}
