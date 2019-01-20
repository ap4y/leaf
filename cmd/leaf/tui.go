package main

import (
	"fmt"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

type align int

const (
	alignLeft align = iota
	alignCenter
	alignRight
)

// TUI implements terminal UI.
type TUI struct {
	userInput []rune
}

// NewTUI construct a new TUI instance.
func NewTUI() *TUI {
	return &TUI{userInput: make([]rune, 0)}
}

// Render renders current ui state using termbox.
func (ui *TUI) Render(s *SessionState) error {
	ui.draw(s)

	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return nil
			}

			if s.Step == StepFinished {
				break
			}

			if s.Step == StepScore {
				s.advance()
				break
			}

			if ev.Key == termbox.KeyEnter {
				s.resolveAnswer(string(ui.userInput))
				ui.userInput = make([]rune, 0)
			} else if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
				if len(ui.userInput) > 0 {
					ui.userInput = ui.userInput[:len(ui.userInput)-1]
				}
			} else {
				var ch rune
				if ev.Key == termbox.KeySpace {
					ch = ' '
				} else {
					ch = ev.Ch
				}

				ui.userInput = append(ui.userInput, ch)
			}
		case termbox.EventError:
			return ev.Err
		}

		ui.draw(s)
	}
}

func (ui *TUI) draw(s *SessionState) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	w, h := termbox.Size()

	write(fmt.Sprintf("    Deck: %s", s.DeckName), 1, 1, 0, 0, 0)
	write(fmt.Sprintf("Progress: %d/%d", s.Total-s.Left, s.Total), 1, 2, 0, 0, 0)

	if s.Step == StepFinished {
		write("no more cards!", w/2, h/2-4, alignCenter, termbox.ColorGreen, 0)
		return
	}

	write(s.Question, w/2, h/2-4, alignCenter, termbox.ColorYellow|termbox.AttrBold, 0)
	write("(type answer below)", w/2, h/2-3, alignCenter, 0, 0)

	answer := s.Answer
	x := (w / 2) - (runewidth.StringWidth(answer) / 2)
	inputBox := []rune{}
	for range answer {
		inputBox = append(inputBox, '_')
	}
	write(string(inputBox)+string('⏎'), x, h/2, 0, termbox.ColorWhite, 0)

	switch s.Step {
	case StepAnswering:
		input := strings.Replace(string(ui.userInput), " ", "␣", -1)
		write(input, x, h/2, 0, termbox.ColorGreen, 0)
	case StepScore:
		if s.Result {
			write("✓", w/2, (h/2)+2, alignCenter, termbox.ColorGreen|termbox.AttrBold, 0)
		} else {
			write("✕", w/2, (h/2)+2, alignCenter, termbox.ColorRed|termbox.AttrBold, 0)
			write(answer, w/2, (h/2)+3, alignCenter, termbox.ColorWhite, 0)
		}
	}
}

func write(text string, x, y int, align align, fg, bg termbox.Attribute) {
	var xOffset int
	switch align {
	case alignLeft:
		xOffset = x
	case alignCenter:
		xOffset = x - runewidth.StringWidth(text)/2
	case alignRight:
		xOffset = x - runewidth.StringWidth(text)
	}

	for _, c := range text {
		termbox.SetCell(xOffset, y, c, fg, bg)
		xOffset += runewidth.RuneWidth(c)
	}
}
