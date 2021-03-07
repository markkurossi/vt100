//
// Copyright (c) 2018-2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"fmt"
	"io"
	"strings"
)

// Point defines a 2D point.
type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

// Equal tests if the argument point is equal to this point.
func (p Point) Equal(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}

var (
	zeroPoint = Point{}
)

// RGBA defines display color value.s
type RGBA uint32

// Emulator color codes.
const (
	Black       = RGBA(0x000000ff)
	Red         = RGBA(0xcd0000ff)
	Green       = RGBA(0x00cd00ff)
	Yellow      = RGBA(0xcdcd00ff)
	Blue        = RGBA(0x0000eeff)
	Magenta     = RGBA(0xcd00cdff)
	Cyan        = RGBA(0x00cdcdff)
	White       = RGBA(0xe5e5e5ff)
	BrightWhite = RGBA(0xffffffff)
)

const (
	debug = false
)

// Char defines the column character and properties in emulator
// display.
type Char struct {
	Code       rune
	Foreground RGBA
	Background RGBA
	Bold       bool
	Italic     bool
	Underline  bool
}

// Clone creates a new character with the argument code. All other
// character attributes are copied.
func (ch Char) Clone(code rune) Char {
	result := ch
	result.Code = code
	return result
}

// CharDisplay implements terminal display.
type CharDisplay interface {
	// Size returns the display size.
	Size() Point
	// Clear clears the display region (inclusively).
	Clear(from, to Point)
	// DECALN fills the screen with 'E'.
	DECALN(size Point)
	// Set sets the character at the specified point.
	Set(p Point, char Char)
	// InsertChars insert count number of characters to the specified
	// point.
	InsertChars(size, p Point, count int)
	// DeleteChars delets count number of characters from the
	// specified point.
	DeleteChars(size, p Point, count int)
	// ScrollUp scrolls the screen up count lines.
	ScrollUp(count int)
}

// Emulator implements terminal emulator.
type Emulator struct {
	display      CharDisplay
	Size         Point
	scrollTop    int
	scrollBottom int
	Cursor       Point
	Default      Char
	ch           Char
	overflow     bool
	state        *state
	stdout       io.Writer
	stderr       io.Writer
}

// NewEmulator creates a new terminal emulator.
func NewEmulator(stdout, stderr io.Writer, display CharDisplay) *Emulator {
	e := &Emulator{
		display: display,
		Default: Char{
			Foreground: Black,
			Background: BrightWhite,
		},
		state:  stStart,
		stdout: stdout,
		stderr: stderr,
	}
	e.Reset()
	return e
}

// Reset resets the emulator to initial state.
func (e *Emulator) Reset() {
	e.Size = e.display.Size()
	e.ch = e.Default
	e.clear(true, true)
}

// Resize sets emulator display area.
func (e *Emulator) Resize(width, height int) {
	e.Size = e.display.Size()
	if e.Size.X > width {
		e.Size.X = width
	}
	if e.Size.Y > height {
		e.Size.Y = height
	}
}

func (e *Emulator) setState(state *state) {
	e.state = state
	e.state.reset()
}

func (e *Emulator) output(format string, a ...interface{}) {
	if e.stdout == nil {
		return
	}
	e.stdout.Write([]byte(fmt.Sprintf(format, a...)))
}

func (e *Emulator) debug(format string, a ...interface{}) {
	if e.stderr == nil {
		return
	}
	msg := fmt.Sprintf(format, a...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	e.stderr.Write([]byte(msg))
}

func (e *Emulator) setIconName(name string) {
	e.debug("Icon Name: %s", name)
}

func (e *Emulator) setWindowTitle(name string) {
	e.debug("Window Title: %s", name)
}

func (e *Emulator) clearLine(line, from, to int) {
	if line < 0 || line >= e.Size.Y {
		return
	}
	if to >= e.Size.X {
		to = e.Size.X - 1
	}
	e.display.Clear(Point{
		X: from,
		Y: line,
	}, Point{
		X: to,
		Y: line,
	})
}

func (e *Emulator) clear(start, end bool) {
	if start {
		if e.Cursor.Y > 0 {
			e.display.Clear(zeroPoint, Point{
				X: e.Size.X - 1,
				Y: e.Cursor.Y - 1,
			})
		}
		e.display.Clear(Point{
			X: 0,
			Y: e.Cursor.Y,
		}, Point{
			X: e.Cursor.X,
			Y: e.Cursor.Y,
		})
	}
	if end {
		e.display.Clear(Point{
			X: e.Cursor.X,
			Y: e.Cursor.Y,
		}, Point{
			X: e.Size.X - 1,
			Y: e.Cursor.Y,
		})
		e.display.Clear(Point{
			Y: e.Cursor.Y + 1,
		}, Point{
			X: e.Size.X - 1,
			Y: e.Size.Y - 1,
		})
	}
}

func (e *Emulator) moveTo(row, col int) {
	if col < 0 {
		col = 0
	}
	if col >= e.Size.X {
		col = e.Size.X - 1
	}
	e.Cursor.X = col

	if row < 0 {
		row = 0
	}
	if row >= e.Size.Y {
		e.scrollUp(e.Size.Y - row + 1)
		row = e.Size.Y - 1
	}
	e.Cursor.Y = row
	e.overflow = false
}

func (e *Emulator) scrollUp(count int) {
	if count >= e.Size.Y {
		e.clear(true, true)
		return
	}
	e.display.ScrollUp(count)

	for i := 0; i < count; i++ {
		e.clearLine(e.Size.Y-1-i, 0, e.Size.X)
	}
}

func (e *Emulator) insertChar(code int) {
	if e.overflow {
		if e.Cursor.Y+1 >= e.Size.Y {
			e.scrollUp(1)
			e.moveTo(e.Cursor.Y, 0)
		} else {
			e.moveTo(e.Cursor.Y+1, 0)
		}
	}
	e.display.Set(e.Cursor, e.ch.Clone(rune(code)))
	if e.Cursor.X+1 >= e.Size.X {
		e.overflow = true
	} else {
		e.moveTo(e.Cursor.Y, e.Cursor.X+1)
	}
}

func (e *Emulator) insertChars(row, col, count int) {
	if row < 0 {
		row = 0
	} else if row >= e.Size.Y {
		row = e.Size.Y - 1
	}
	if col < 0 {
		col = 0
	} else if col >= e.Size.X {
		return
	}
	if col+count >= e.Size.X {
		e.clearLine(row, col, e.Size.X)
		return
	}
	e.display.InsertChars(e.Size, Point{
		Y: row,
		X: col,
	}, count)
}

func (e *Emulator) deleteChars(row, col, count int) {
	if row < 0 {
		row = 0
	} else if row >= e.Size.Y {
		row = e.Size.Y - 1
	}
	if col < 0 {
		col = 0
	} else if col >= e.Size.X {
		return
	}
	if col+count >= e.Size.X {
		e.clearLine(row, col, e.Size.X)
		return
	}
	e.display.DeleteChars(e.Size, Point{
		Y: row,
		X: col,
	}, count)
}

// Input runs the terminal emulation with the next input code.
func (e *Emulator) Input(code int) {
	if debug {
		e.debug("Emulator.Input: %s<-0x%x (%d) '%c'", e.state, code, code, code)
	}
	next := e.state.input(e, code)
	if next != nil {
		if debug {
			e.debug("Emulator.Input: %s->%s", e.state, next)
		}
		e.setState(next)
	}
}
