package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type State int

const (
	StateAnswering State = iota
	StateResult
	StateFinished
)

type UI struct {
	State     State
	UserInput []rune
	DeckName  string
	Total     int
	Left      int
	Question  string
	Answer    string
	Result    bool
}

func (ui *UI) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	w, h := termbox.Size()

	write(fmt.Sprintf("    Deck: %s", ui.DeckName), 1, 1, 0, 0, 0)
	write(fmt.Sprintf("Progress: %d/%d", ui.Total-ui.Left, ui.Total), 1, 2, 0, 0, 0)

	if ui.State == StateFinished {
		write("no more cards!", w/2, h/2-4, AlignCenter, termbox.ColorGreen, 0)
		return
	}

	write(fmt.Sprintf("%s", ui.Question), w/2, h/2-4, AlignCenter, termbox.ColorYellow|termbox.AttrBold, 0)
	write("(type answer below)", w/2, h/2-3, AlignCenter, 0, 0)

	answer := ui.Answer
	x := (w / 2) - (utf8.RuneCountInString(answer) / 2)
	inputBox := []rune{}
	for _, ch := range answer {
		if ch == ' ' {
			inputBox = append(inputBox, '␣')
		} else {
			inputBox = append(inputBox, '_')
		}
	}
	write(string(inputBox)+string('⏎'), x, h/2, 0, termbox.ColorWhite, 0)

	switch ui.State {
	case StateAnswering:
		input := strings.Replace(string(ui.UserInput), " ", "␣", -1)
		write(input, x, h/2, 0, termbox.ColorGreen, 0)
	case StateResult:
		if ui.Result {
			write("✓", w/2, (h/2)+2, AlignCenter, termbox.ColorGreen|termbox.AttrBold, 0)
		} else {
			write("✕", w/2, (h/2)+2, AlignCenter, termbox.ColorRed|termbox.AttrBold, 0)
			write(answer, w/2, (h/2)+3, AlignCenter, termbox.ColorWhite, 0)
		}
	}
}

func write(text string, x, y int, align Align, fg, bg termbox.Attribute) {
	switch align {
	case AlignLeft:
		x = x
	case AlignCenter:
		x = x - utf8.RuneCountInString(text)/2
	case AlignRight:
		x = x - utf8.RuneCountInString(text)
	}

	for _, c := range text {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
