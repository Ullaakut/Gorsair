package wow

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gernest/wow/spin"
	"golang.org/x/crypto/ssh/terminal"
)

//go:generate go run ./scripts/generate_spin.go
const erase = "\033[2K\r"

// Wow writes beautiful spinners on the terminal.
type Wow struct {
	txt        string
	s          spin.Spinner
	out        io.Writer
	running    bool
	done       func()
	mu         sync.RWMutex
	IsTerminal bool
}

// New creates a new wow instance ready to start spinning.
func New(o io.Writer, s spin.Spinner, text string, options ...func(*Wow)) *Wow {
	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	wow := Wow{out: o, s: s, txt: text, IsTerminal: isTerminal}

	for _, option := range options {
		option(&wow)
	}
	return &wow
}

// Start starts the spinner. The frames are written based on the spinner
// interval.
func (w *Wow) Start() {
	if !w.running {
		ctx, done := context.WithCancel(context.Background())
		t := time.NewTicker(time.Duration(w.s.Interval) * time.Millisecond)
		w.done = done
		w.running = true
		go func() {
			at := 0
			for {
				select {
				case <-ctx.Done():
					t.Stop()
					break
				case <-t.C:
					txt := erase + w.s.Frames[at%len(w.s.Frames)] + w.txt
					if w.IsTerminal {
						fmt.Fprint(w.out, txt)
					}
					at++
				}
			}
		}()
	}
}

// Stop stops the spinner
func (w *Wow) Stop() {
	if w.done != nil {
		w.done()
	}
	w.running = false
}

// Spinner sets s to the current spinner
func (w *Wow) Spinner(s spin.Spinner) *Wow {
	w.Stop()
	w.s = s
	w.Start()
	return w
}

// Text adds text to the current spinner
func (w *Wow) Text(txt string) *Wow {
	w.mu.Lock()
	w.txt = txt
	w.mu.Unlock()
	return w
}

// Persist writes the last character of the currect spinner frames together with
// the text on stdout.
//
// A new line is added at the end to ensure the text stay that way.
func (w *Wow) Persist() {
	w.Stop()
	at := len(w.s.Frames) - 1
	txt := erase + w.s.Frames[at] + w.txt + "\n"
	if w.IsTerminal {
		fmt.Fprint(w.out, txt)
	}
}

// PersistWith writes the last frame of s together with text with a new line
// added to make it stick.
func (w *Wow) PersistWith(s spin.Spinner, text string) {
	w.Stop()
	var a string
	if len(s.Frames) > 0 {
		a = s.Frames[len(s.Frames)-1]
	}
	txt := erase + a + text + "\n"
	if w.IsTerminal {
		fmt.Fprint(w.out, txt)
	}
}

// ForceOutput forces all output even if not not outputting directly to a terminal
func ForceOutput(w *Wow) {
	w.IsTerminal = true
}

// LogSymbol is a log severity level
type LogSymbol uint

// common log symbos
const (
	INFO LogSymbol = iota
	SUCCESS
	WARNING
	ERROR
)

func (s LogSymbol) String() string {
	if runtime.GOOS != "windows" {
		switch s {
		case INFO:
			return "ℹ"
		case SUCCESS:
			return "✔"
		case WARNING:
			return "⚠"
		case ERROR:
			return "✖"
		}
	} else {
		switch s {
		case INFO:
			return "i"
		case SUCCESS:
			return "v"
		case WARNING:
			return "!!"
		case ERROR:
			return "x"
		}
	}
	return ""
}
