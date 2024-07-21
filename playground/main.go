package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	msgChan := make(chan []byte)
	app := tview.NewApplication()
	inputField :=
		tview.NewInputField().
			SetLabel("Enter text: ").
			SetFieldWidth(30)
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msgChan <- []byte(inputField.GetText())
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

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(inputField, 1, 0, true)

	go func() {
		for {
			select {
			case b := <-msgChan:
				fmt.Fprintf(textView, "%s\n", string(b))
				textView.ScrollToEnd()
			}
		}
	}()

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

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
