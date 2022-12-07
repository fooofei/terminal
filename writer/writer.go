// Package writer wraps the terminal, for simple new instance and write as a io.Writer.
package writer

import (
	"bytes"
	"fmt"

	"github.com/fooofei/terminal/runes"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

func init() {
	encoding.Register()
}

type Terminal struct {
	Screen    tcell.Screen // the tcell instance
	screenBuf *bytes.Buffer
	stopped   chan struct{}
	finiFunc  func()
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
		Screen:    s,
		screenBuf: bytes.NewBufferString(""),
		stopped:   make(chan struct{}),
	}
	for _, opt := range opts {
		opt(t)
	}
	go t.pollEvent()
	return t, nil
}

func (t *Terminal) markStopped() {
	select {
	case <-t.stopped:
	default:
		close(t.stopped)
	}
}

func (t *Terminal) isStopped() bool {
	select {
	case <-t.stopped:
		return true
	default:
		return false
	}
}

func (t *Terminal) Close() error {
	t.Screen.Fini()
	t.markStopped()
	if t.finiFunc != nil {
		t.finiFunc()
	}
	return nil
}

func (t *Terminal) Clear() {
	t.screenBuf.Reset()
}

func (t *Terminal) ForceClearScreen() {
	t.Screen.Clear()
}

// setContent wraps the screen.SetContent method for put first runes to the head
func (t *Terminal) setContentXY(x, y int, runeList []rune, style tcell.Style) {
	var width, _ = t.Screen.Size()
	for i := len(runeList); i < width; i++ {
		runeList = append(runeList, ' ')
	}
	t.Screen.SetContent(x, y, runeList[0], runeList[1:], style)
}

func (t *Terminal) setContent(content []byte) {
	var y int
	var tailRuneList, _ = runes.DecodeRuneOnNewLine(content, func(rs []rune) {
		t.setContentXY(0, y, rs, tcell.StyleDefault)
		y += 1
	})
	t.setContentXY(0, y, tailRuneList, tcell.StyleDefault)
}

func (t *Terminal) Write(p []byte) (int, error) {
	if t.isStopped() {
		return 0, fmt.Errorf("cannot write anymore, terminal screen alreay stopped")
	}
	var n, err = t.screenBuf.Write(p)
	if n > 0 {
		t.setContent(t.screenBuf.Bytes())
		t.Show()
	}
	return n, err
}

func (t *Terminal) Show() {
	t.Screen.Show()
}

func (t *Terminal) Sync() {
	t.Screen.Sync()
}

func (t *Terminal) GetContent() []byte {
	return t.screenBuf.Bytes()
}

// Flush will use buffered content to redraw terminal
func (t *Terminal) Flush() {
	t.ForceClearScreen()
	t.setContent(t.screenBuf.Bytes())
	t.Show()
}

func (t *Terminal) pollEvent() {
	s := t.Screen
loop:
	for {
		ev := s.PollEvent()
		if ev == nil {
			// got quit
			break loop
		}
		switch v := ev.(type) {
		case *tcell.EventResize:
			t.Sync()
		case *tcell.EventKey:
			if v.Key() == tcell.KeyEscape || v.Key() == tcell.KeyCtrlC {
				t.Close()
				break loop
			}
		}
	}
}
