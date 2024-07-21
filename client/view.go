package client

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
)

type app struct {
	sendMessageChan chan []byte
	ui              io.Writer
	application     *tview.Application

	textView   *tview.TextView
	inputField *tview.InputField
}

func newView(sendMessageChan chan []byte) *app {
	return &app{
		sendMessageChan: sendMessageChan,
	}
}

func (v *app) setUp() {
	app := tview.NewApplication()
	inputField :=
		tview.NewInputField().
			SetLabel(">: ").
			SetFieldWidth(30)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			v.sendMessageChan <- []byte(inputField.GetText())
			inputField.SetText("")
		}
	})

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetScrollable(true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			row, n := textView.GetScrollOffset()
			textView.ScrollTo(row-1, n)
		case tcell.KeyDown:
			row, n := textView.GetScrollOffset()
			textView.ScrollTo(row+1, n)
		}
		return event
	})

	v.application = app
	v.ui = textView
	v.textView = textView
	v.inputField = inputField
}

func (v *app) UI() io.Writer {
	return v.ui
}

func (v *app) Run() error {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.textView, 0, 1, false).
		AddItem(v.inputField, 1, 0, true)
	if err := v.application.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
	return nil
}
