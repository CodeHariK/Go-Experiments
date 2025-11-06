package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Blur() tea.Msg
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type ShortAnswerField struct {
	textinput textinput.Model
}

type LongAnswerField struct {
	textarea textarea.Model
}

func (sa *ShortAnswerField) Value() string {
	return sa.textinput.Value()
}

func (la *LongAnswerField) Value() string {
	return la.textarea.Value()
}

func (sa *ShortAnswerField) Blur() tea.Msg {
	return sa.textinput.Blur
}

func (la *LongAnswerField) Blur() tea.Msg {
	return la.textarea.Blur
}

func (sa *ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	sa.textinput, cmd = sa.textinput.Update(msg)
	return sa, cmd
}

func (la *LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	la.textarea, cmd = la.textarea.Update(msg)
	return la, cmd
}

func (sa *ShortAnswerField) View() string {
	return sa.textinput.View()
}

func (la *LongAnswerField) View() string {
	return la.textarea.View()
}

// TextInput
func NewShortAnswerField() Input {
	ti := textinput.New()
	ti.Placeholder = "Your answer here"
	ti.Focus()
	return &ShortAnswerField{ti}
}

// TextArea
func NewLongAnswerField() Input {
	ta := textarea.New()
	ta.Placeholder = "Your answer here"
	ta.Focus()
	return &LongAnswerField{ta}
}
