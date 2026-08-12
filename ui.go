package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (t *Terminal) run() error {
	var ops op.Ops
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	var keyTag struct{}

	for {
		e := t.window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			event.Op(gtx.Ops, &keyTag)
			gtx.Execute(key.FocusCmd{Tag: &keyTag})

			for {
				ev, ok := gtx.Event(
					key.FocusFilter{Target: &keyTag},
					key.Filter{Focus: &keyTag},
				)
				if !ok {
					break
				}
				switch ke := ev.(type) {
				case key.Event:
					if ke.State == key.Press {
						t.handleKeyEvent(ke)
					}
				case key.EditEvent:
					if ke.Text != "" {
						t.pty.Write([]byte(ke.Text))
					}
				}
			}

			paint.FillShape(gtx.Ops, rgba(0x1e1e1e, 217), clip.Rect{
				Max: gtx.Constraints.Max,
			}.Op())

			t.mu.Lock()
			layout.Inset{Top: unit.Dp(8), Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				list := layout.List{Axis: layout.Vertical}
				return list.Layout(gtx, len(t.lines), func(gtx layout.Context, i int) layout.Dimensions {
					return t.layoutLine(gtx, th, t.lines[i])
				})
			})
			t.mu.Unlock()

			e.Frame(gtx.Ops)
		}
	}
}

func (t *Terminal) layoutLine(gtx layout.Context, th *material.Theme, line []cell) layout.Dimensions {
	if len(line) == 0 {
		lbl := material.Label(th, unit.Sp(14), " ")
		lbl.Font.Typeface = "monospace"
		return lbl.Layout(gtx)
	}

	children := make([]layout.FlexChild, 0, len(line))
	start := 0
	flush := func(end int) {
		if end <= start {
			return
		}
		run := line[start:end]
		text := make([]rune, len(run))
		for i, c := range run {
			text[i] = c.r
		}
		fg, bg := run[0].fg, run[0].bg
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(14), string(text))
			lbl.Color = fg
			lbl.Font.Typeface = "monospace"
			if bg != defaultBG {
				macro := op.Record(gtx.Ops)
				dims := lbl.Layout(gtx)
				call := macro.Stop()
				paint.FillShape(gtx.Ops, bg, clip.Rect{Max: dims.Size}.Op())
				call.Add(gtx.Ops)
				return dims
			}
			return lbl.Layout(gtx)
		}))
	}

	for i := 1; i < len(line); i++ {
		if line[i].fg != line[i-1].fg || line[i].bg != line[i-1].bg {
			flush(i)
			start = i
		}
	}
	flush(len(line))

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

func (t *Terminal) handleKeyEvent(e key.Event) {
	var data []byte

	switch e.Name {
	case key.NameReturn, key.NameEnter:
		data = []byte("\r")
	case key.NameDeleteBackward:
		data = []byte{0x7f}
	case key.NameDeleteForward:
		data = []byte("\x1b[3~")
	case key.NameTab:
		data = []byte("\t")
	case key.NameEscape:
		data = []byte{0x1b}
	case key.NameUpArrow:
		data = []byte("\x1b[A")
	case key.NameDownArrow:
		data = []byte("\x1b[B")
	case key.NameRightArrow:
		data = []byte("\x1b[C")
	case key.NameLeftArrow:
		data = []byte("\x1b[D")
	default:
		// Ctrl+C, Ctrl+D, etc. Plain printable characters are delivered
		if e.Modifiers.Contain(key.ModCtrl) && len(e.Name) == 1 {
			r := e.Name[0]
			if r >= 'A' && r <= 'Z' {
				data = []byte{r - 'A' + 1}
			} else if r >= 'a' && r <= 'z' {
				data = []byte{r - 'a' + 1}
			}
		}
	}

	if len(data) > 0 {
		t.pty.Write(data)
	}
}
