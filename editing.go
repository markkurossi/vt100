//
// vt100.go
//
// Copyright (c) 2018-2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"fmt"
	"io"
)

// CursorUp moves cursor one line up.
func CursorUp(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'A'})
	return err
}

// CursorDown moves cursor one line down.
func CursorDown(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'B'})
	return err
}

// CursorForward moves the cursor one column right. Stops at the right
// edge of screen.
func CursorForward(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'C'})
	return err
}

// CursorBackward moves the cursor one column left. Stops at the left
// edge of screen.
func CursorBackward(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'D'})
	return err
}

// ScrollUp scrolls the screen one line up.
func ScrollUp(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'S'})
	return err
}

// ScrollDown scrolls the screen one line down.
func ScrollDown(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'T'})
	return err
}

// Backspace moves cursor one column left. Backspace does nothing if
// the cursor is already at the leftmost column.
func Backspace(out io.Writer) error {
	_, err := out.Write([]byte{0x08})
	return err
}

// DeleteChar deletes character from the current cursor position.
func DeleteChar(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'P'})
	return err
}

// EraseLineHead clears the current line from the beginning of the
// line to cursor position (inclusively).
func EraseLineHead(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', '1', 'K'})
	return err
}

// EraseLineTail clears the current line from the cursor position to
// the end of line (inclusively).
func EraseLineTail(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'K'})
	return err
}

// EraseLine clears the current line.
func EraseLine(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', '2', 'K'})
	return err
}

// EraseScreenHead clears the screen from the beginning of the screen
// to cursor position (inclusively).
func EraseScreenHead(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', '1', 'J'})
	return err
}

// EraseScreenTail clears the screen from cursor position to the end
// of screen (inclusively).
func EraseScreenTail(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', 'J'})
	return err
}

// EraseScreen clears screen.
func EraseScreen(out io.Writer) error {
	_, err := out.Write([]byte{0x1b, '[', '2', 'J'})
	return err
}

// MoveTo moves cursor to the specified row and column.
func MoveTo(out io.Writer, row, col int) error {
	_, err := out.Write([]byte(fmt.Sprintf("\x1b[%d;%dH", row, col)))
	return err
}
