//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"strings"
)

func DisplayWidth(data string) (int, error) {
	var width int

	emul := NewEmulator(func(e *Emulator, ch rune) error {
		width++
		return nil
	})
	err := emul.Input([]rune(data))
	if err != nil {
		return 0, err
	}

	return width, nil
}

func Trim(data string) (string, error) {
	var sb strings.Builder

	emul := NewEmulator(func(e *Emulator, ch rune) error {
		_, err := sb.WriteRune(ch)
		return err
	})
	err := emul.Input([]rune(data))
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}