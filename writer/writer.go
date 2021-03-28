// package writer wraps the terminal, for simple new instance and write as a io.Writer.
package writer

import (
	"fmt"

	"github.com/fooofei/terminal/runes"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"go.uber.org/atomic"
)

func init() {
	encoding.Register()
}

type Terminal struct {
	Screen   tcell.Screen // the tcell instance
	stopped  *atomic.Bool // marked screen be finished
	x        int
	y        int
	finiFunc func()
}

type Opts func(*Terminal)

// WithOnFinish will set called f when terminal finished
// f can be called in anther goroutine, must can be thread-safe
func WithOnFinish(f func()) Opts {
	return func(terminal *Terminal) {
		terminal.finiFunc = f
	}
}

func New(opts ...Opts) (*Terminal, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err = s.Init(); err != nil {
		s.Fini()
		return nil, err
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	t := &Terminal{
		Screen:  s,
		stopped: atomic.NewBool(false),
		x:       0,
		y:       0,
	}
	for _, opt := range opts {
		opt(t)
	}
	go t.pollEvent()
	return t, nil
}

func (term *Terminal) Close() error {
	term.Screen.Fini()
	term.stopped.Store(true)
	term.x = 0
	term.y = 0
	return nil
}

func (term *Terminal) Clear() {
	term.x = 0
	term.y = 0
	term.Screen.Clear()
}

// setContent wraps the screen.SetContent method for put first runes to the head
func (term *Terminal) setContent(x, y int, runes []rune, style tcell.Style) {
	if len(runes) <= 0 {
		return
	}
	term.Screen.SetContent(x, y, runes[0], runes[1:], style)
}

func (term *Terminal) Write(p []byte) (int, error) {
	if term.stopped.Load() {
		return 0, fmt.Errorf("terminal screen alreay stopped")
	}
	tailRunes, size := runes.DecodeRuneOnNewLine(p, func(rs []rune) {
		term.setContent(term.x, term.y, rs, tcell.StyleDefault)
		term.y += 1
		term.x = 0
	})
	term.setContent(term.x, term.y, tailRunes, tcell.StyleDefault)
	term.x += len(tailRunes)
	term.Screen.Show()
	return size, nil
}

func (term *Terminal) Show() {
	term.Screen.Show()
}

func (term *Terminal) Sync() {
	term.Screen.Sync()
}

func (term *Terminal) pollEvent() {
	s := term.Screen
loop:
	for {
		ev := s.PollEvent()
		if ev == nil {
			// got quit
			break loop
		}
		switch v := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if v.Key() == tcell.KeyEscape || v.Key() == tcell.KeyCtrlC {
				term.stopped.Store(true)
				s.Fini()
				if term.finiFunc != nil {
					term.finiFunc()
				}
				break loop
			}
		}
	}
}
