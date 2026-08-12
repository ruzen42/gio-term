package main

import (
	log "gio-term/logger"
	"os"
	"os/exec"

	"gioui.org/app"
	"gioui.org/unit"

	"github.com/creack/pty"
)

func main() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	log.Debug("starting shell: " + shell)

	cmd := exec.Command(shell)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal("error while creating pty: " + err.Error())
	}
	defer ptmx.Close()

	term := newTerminal(ptmx)

	go term.readPTY()

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Gio Terminal"), app.Size(unit.Dp(800), unit.Dp(600)))
		term.window = w

		if err := term.run(); err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	}()

	app.Main()
}
