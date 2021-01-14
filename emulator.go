//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

// Emulator implements VT100 emulator.
type Emulator struct {
	output Output
	state  state
}

type Output func(e *Emulator, ch rune) error

func NewEmulator(output Output) *Emulator {
	return &Emulator{
		output: output,
		state:  stStart,
	}
}

type state func(e *Emulator, ch rune) error

var states map[rune]state

func init() {
	states = map[rune]state{
		0x1b: stESC,
		0x9b: stCSI,
	}
}

func stStart(e *Emulator, ch rune) error {
	next, ok := states[ch]
	if !ok {
		return e.output(e, ch)
	}
	e.state = next
	return nil
}

func stESC(e *Emulator, ch rune) error {
	switch ch {
	case '[':
		e.state = stCSI

	default:
		e.state = stStart
	}
	return nil
}

func stCSI(e *Emulator, ch rune) error {
	if isAlphabetic(ch) {
		e.state = stStart
	}
	// XXX ignore CSI for now
	return nil
}

func (e *Emulator) Input(data []rune) error {
	for _, r := range data {
		err := e.state(e, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func isAlphabetic(ch rune) bool {
	return 0x40 <= ch && ch <= 0x74
}
