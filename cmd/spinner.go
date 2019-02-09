package main

import (
	"os"
	"time"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func updateSpinner(w *wow.Wow, text string, disabled bool) {
	if !disabled {
		w.Text(" " + text)
	}
}

func startSpinner(disabled bool) *wow.Wow {
	if !disabled {
		w := wow.New(os.Stdout, spin.Get(spin.Dots), " Loading Gorsair...")
		w.Start()
		return w
	}
	return nil
}

// HACK: Waiting for a fix to issue
// https://github.com/gernest/wow/issues/5
func clearOutput(w *wow.Wow, disabled bool) {
	if !disabled {
		w.Text("\b")
		time.Sleep(80 * time.Millisecond)
		w.Stop()
	}
}
