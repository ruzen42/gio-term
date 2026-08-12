package main

import "image/color"

func rgb(c uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(c >> 16),
		G: uint8(c >> 8),
		B: uint8(c),
		A: 0xff,
	}
}

func rgba(c uint32, a uint8) color.NRGBA {
	return color.NRGBA{
		R: uint8(c >> 16),
		G: uint8(c >> 8),
		B: uint8(c),
		A: a,
	}
}

var (
	defaultFG = rgb(0xd4d4d4)
	defaultBG = rgb(0x1e1e1e)
)

var ansi16 = [16]color.NRGBA{
	rgb(0x000000), // 0 black
	rgb(0xcd3131), // 1 red
	rgb(0x0dbc79), // 2 green
	rgb(0xe5e510), // 3 yellow
	rgb(0x2472c8), // 4 blue
	rgb(0xbc3fbc), // 5 magenta
	rgb(0x11a8cd), // 6 cyan
	rgb(0xe5e5e5), // 7 white
	rgb(0x666666), // 8 bright black
	rgb(0xf14c4c), // 9 bright red
	rgb(0x23d18b), // 10 bright green
	rgb(0xf5f543), // 11 bright yellow
	rgb(0x3b8eea), // 12 bright blue
	rgb(0xd670d6), // 13 bright magenta
	rgb(0x29b8db), // 14 bright cyan
	rgb(0xe5e5e5), // 15 bright white
}

// ansi256 returns the color for an xterm 256-color palette index (SGR
// "38;5;N" / "48;5;N"). Indices 0-15 mirror the 16-color palette, 16-231
// are a 6x6x6 color cube, and 232-255 are a grayscale ramp.
func ansi256(n int) color.NRGBA {
	switch {
	case n < 16:
		return ansi16[n]
	case n < 232:
		n -= 16
		r := (n / 36) % 6
		g := (n / 6) % 6
		b := n % 6
		step := func(v int) uint8 {
			if v == 0 {
				return 0
			}
			return uint8(55 + v*40)
		}
		return color.NRGBA{R: step(r), G: step(g), B: step(b), A: 0xff}
	default:
		v := uint8(8 + (n-232)*10)
		return color.NRGBA{R: v, G: v, B: v, A: 0xff}
	}
}
