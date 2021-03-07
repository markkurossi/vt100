//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

var (
	_ CharDisplay = &Display{}
)

// Display implements fixed size CharDisplay.
type Display struct {
	Blank Char
	size  Point
	Lines [][]Char
}

// NewDisplay creates a display with the given dimensions.
func NewDisplay(width, height int) *Display {
	d := &Display{
		Blank: Char{
			Code:       0xa0,
			Foreground: Black,
			Background: White,
		},
		size: Point{
			X: width,
			Y: height,
		},
	}
	d.Resize(width, height)
	return d
}

// Resize resizes the display to given dimensions.
func (d *Display) Resize(width, height int) {
	d.size.X = width
	d.size.Y = height

	for row := 0; row < height; row++ {
		var line []Char
		var start int
		if row < len(d.Lines) {
			line = d.Lines[row]
			start = len(line)
		} else {
			line = make([]Char, width)
			start = 0
			d.Lines = append(d.Lines, line)
		}
		for col := start; col < width; col++ {
			line[col] = d.Blank
		}
	}
}

// Size implements the CharDisplay.Size function.
func (d *Display) Size() Point {
	return d.size
}

// Clear implements the CharDisplay.Clear function.
func (d *Display) Clear(from, to Point) {
	for y := from.Y; y <= to.Y; y++ {
		for x := from.X; x <= to.X; x++ {
			d.Lines[y][x] = d.Blank
		}
	}
}

// DECALN implements the CharDisplay.DECALN function.
func (d *Display) DECALN(size Point) {
	ch := d.Blank
	ch.Code = 'E'

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			d.Lines[y][x] = ch
		}
	}
}

// Set implements the CharDisplay.Set function.
func (d *Display) Set(p Point, char Char) {
	d.Lines[p.Y][p.X] = char
}

// InsertChars implements the CharDisplay.InsertChars function.
func (d *Display) InsertChars(size, p Point, count int) {
	var line []Char
	var x int

	for ; x < p.X; x++ {
		line = append(line, d.Lines[p.Y][x])
	}
	for i := 0; i < count; i++ {
		line = append(line, d.Blank)
	}
	for ; x+count < size.X; x++ {
		line = append(line, d.Lines[p.Y][x])
	}
	line = append(line, d.Lines[p.Y][size.X:]...)
	d.Lines[p.Y] = line
}

// DeleteChars implements the CharDisplay.DeleteChars function.
func (d *Display) DeleteChars(size, p Point, count int) {
	var line []Char
	var x int

	for ; x < p.X; x++ {
		line = append(line, d.Lines[p.Y][x])
	}
	for ; x+count < size.X; x++ {
		line = append(line, d.Lines[p.Y][x+count])
	}
	for i := 0; i < count; i++ {
		line = append(line, d.Blank)
	}
	line = append(line, d.Lines[p.Y][size.X:]...)
	d.Lines[p.Y] = line
}

// ScrollUp implements the CharDisplay.ScrollUp function.
func (d *Display) ScrollUp(count int) {
	for i := 0; i < count; i++ {
		saved := d.Lines[0]
		d.Lines = append(d.Lines[1:], saved)
	}
}
