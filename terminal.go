package main

import (
	log "gio-term/logger"
	"image/color"
	"io"
	"os"
	"sync"

	"gioui.org/app"
)

// cell is a single character position on screen, with its own colors so
// ANSI SGR sequences can color individual runs of text.
type cell struct {
	r  rune
	fg color.NRGBA
	bg color.NRGBA
}

// ansiState is the parser state machine for escape sequences split across
// multiple reads from the PTY.
type ansiState int

const (
	stateGround ansiState = iota
	stateEscape           // saw ESC
	stateCSI              // saw ESC [, collecting a CSI sequence
)

type Terminal struct {
	mu     sync.Mutex
	lines  [][]cell
	pty    *os.File
	window *app.Window

	// current SGR pen state, carried across processBytes calls
	curFG color.NRGBA
	curBG color.NRGBA

	// escape sequence parser state, carried across processBytes calls
	state  ansiState
	csiBuf []byte
}

func newTerminal(ptmx *os.File) *Terminal {
	return &Terminal{
		lines: [][]cell{{}},
		pty:   ptmx,
		curFG: defaultFG,
		curBG: defaultBG,
		state: stateGround,
	}
}

func (t *Terminal) readPTY() {
	buf := make([]byte, 4096)
	for {
		n, err := t.pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Error("Error while readPTY: " + err.Error())
			return
		}

		t.mu.Lock()
		t.processBytes(buf[:n])
		t.mu.Unlock()

		if t.window != nil {
			t.window.Invalidate()
		}
	}
}

func (t *Terminal) processBytes(data []byte) {
	for _, b := range data {
		switch t.state {
		case stateGround:
			t.processGroundByte(b)
		case stateEscape:
			t.processEscapeByte(b)
		case stateCSI:
			t.processCSIByte(b)
		}
	}
}

func (t *Terminal) processGroundByte(b byte) {
	switch b {
	case 0x1b: // ESC
		t.state = stateEscape
	case '\n':
		t.lines = append(t.lines, []cell{})
	case '\r':
	case '\b', 0x7f:
		last := len(t.lines) - 1
		if len(t.lines[last]) > 0 {
			t.lines[last] = t.lines[last][:len(t.lines[last])-1]
		}
	case '\t':
		last := len(t.lines) - 1
		for i := 0; i < 4; i++ {
			t.lines[last] = append(t.lines[last], cell{r: ' ', fg: t.curFG, bg: t.curBG})
		}
	default:
		if b >= 32 && b <= 126 {
			last := len(t.lines) - 1
			t.lines[last] = append(t.lines[last], cell{r: rune(b), fg: t.curFG, bg: t.curBG})
		}
		// Non-ASCII / control bytes outside CSI sequences are dropped.
		// A real UTF-8 decoder could be layered on top of this later.
	}
}

func (t *Terminal) processEscapeByte(b byte) {
	switch b {
	case '[':
		t.state = stateCSI
		t.csiBuf = t.csiBuf[:0]
	default:
		// Unsupported escape (e.g. OSC, charset select) - bail back to
		// ground state without acting on it.
		t.state = stateGround
	}
}

func (t *Terminal) processCSIByte(b byte) {
	// CSI parameter/intermediate bytes are 0x30-0x3f and 0x20-0x2f; the
	// sequence ends on a "final byte" in 0x40-0x7e.
	if b >= 0x40 && b <= 0x7e {
		t.applyCSI(b, t.csiBuf)
		t.state = stateGround
		t.csiBuf = t.csiBuf[:0]
		return
	}
	t.csiBuf = append(t.csiBuf, b)
}

// applyCSI handles the CSI sequences relevant to a text-only terminal:
// SGR (color/style, final byte 'm') and cursor movement is intentionally
// ignored since this terminal has no cursor-addressed grid yet.
func (t *Terminal) applyCSI(final byte, params []byte) {
	if final != 'm' {
		// Cursor movement, erase, etc. Not implemented - see project notes.
		return
	}
	t.applySGR(parseSGRParams(params))
}

// parseSGRParams splits a CSI params buffer like "1;38;2;255;0;0" into ints.
// A missing/empty params buffer means a single implicit 0 (reset).
func parseSGRParams(params []byte) []int {
	if len(params) == 0 {
		return []int{0}
	}
	var out []int
	cur := 0
	has := false
	for _, b := range params {
		if b == ';' {
			if has {
				out = append(out, cur)
			} else {
				out = append(out, 0)
			}
			cur, has = 0, false
			continue
		}
		if b >= '0' && b <= '9' {
			cur = cur*10 + int(b-'0')
			has = true
		}
	}
	if has {
		out = append(out, cur)
	} else {
		out = append(out, 0)
	}
	return out
}

// applySGR updates the current pen (fg/bg) from a parsed list of SGR codes,
// supporting 16-color, 256-color, and truecolor (24-bit) forms.
func (t *Terminal) applySGR(codes []int) {
	for i := 0; i < len(codes); i++ {
		c := codes[i]
		switch {
		case c == 0:
			t.curFG = defaultFG
			t.curBG = defaultBG
		case c == 39:
			t.curFG = defaultFG
		case c == 49:
			t.curBG = defaultBG
		case c >= 30 && c <= 37:
			t.curFG = ansi16[c-30]
		case c >= 90 && c <= 97:
			t.curFG = ansi16[8+c-90]
		case c >= 40 && c <= 47:
			t.curBG = ansi16[c-40]
		case c >= 100 && c <= 107:
			t.curBG = ansi16[8+c-100]
		case c == 38 || c == 48:
			// Extended color: 38;5;N (256-color) or 38;2;R;G;B (truecolor).
			// Same shape for background with 48.
			col, consumed, ok := parseExtendedColor(codes[i+1:])
			if !ok {
				continue
			}
			if c == 38 {
				t.curFG = col
			} else {
				t.curBG = col
			}
			i += consumed
		}
	}
}

// parseExtendedColor parses the tail of an SGR extended color sequence
// (everything after "38" or "48") and returns the color, how many extra
// codes were consumed, and whether parsing succeeded.
func parseExtendedColor(rest []int) (color.NRGBA, int, bool) {
	if len(rest) == 0 {
		return color.NRGBA{}, 0, false
	}
	switch rest[0] {
	case 5: // 256-color palette
		if len(rest) < 2 {
			return color.NRGBA{}, 1, false
		}
		return ansi256(rest[1]), 2, true
	case 2: // truecolor
		if len(rest) < 4 {
			return color.NRGBA{}, len(rest), false
		}
		r, g, b := rest[1], rest[2], rest[3]
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}, 4, true
	default:
		return color.NRGBA{}, 1, false
	}
}
