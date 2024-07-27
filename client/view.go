package client

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	types2 "tchat/internal/types"
)

type app struct {
	sendMessageCh chan []byte
	renderTextCh  chan []string
	joinChannelCh chan types2.Channel
	exitChannelCh chan struct{}
	application   *tview.Application

	lobbyView   *tview.TextView
	currentView *tview.TextView
	inputView   *tview.InputField

	ctx *clientContext
}

func newView(ctx *clientContext, sendMessageChan chan []byte, renderTextCh chan []string, joinChannelCh chan types2.Channel, exitChannelCh chan struct{}) *app {
	return &app{
		sendMessageCh: sendMessageChan,
		renderTextCh:  renderTextCh,
		joinChannelCh: joinChannelCh,
		exitChannelCh: exitChannelCh,
		ctx:           ctx,
	}
}

func (v *app) setUp() {
	app := tview.NewApplication()
	inputField :=
		tview.NewInputField().
			SetLabel(">: ")
	inputField.SetBackgroundColor(tcell.ColorYellow)
	inputField.SetFieldBackgroundColor(tcell.ColorYellow)
	inputField.SetFieldTextColor(tcell.ColorBlack)

	inputField.SetDoneFunc(func(key tcell.Key) {
		txt := inputField.GetText()
		if key == tcell.KeyEnter {
			statePrefix := ""
			// kinda ugly but for now it works, if msg starts with / then it is also not a message command
			// but a channel command
			if v.ctx.currentChannel != nil && strings.Index(txt, "/") != 0 {
				statePrefix = "/message "
			}
			v.sendMessageCh <- []byte(statePrefix + txt)
			inputField.SetText("")
		}
	})
	textView := v.setUpLobbyView(app)
	v.application = app
	v.lobbyView = textView
	v.currentView = textView
	v.inputView = inputField

	v.runChannelConsumers()
}

func (v *app) setUpLobbyView(app *tview.Application) *tview.TextView {
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

	return textView
}

func (v *app) setUpChannelView(app *tview.Application, c types2.Channel) *tview.TextView {
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

	textView.SetText(c.WelcomeMessage)

	return textView
}

func (v *app) runChannelConsumers() {
	go func() {
		for {
			select {
			case text := <-v.renderTextCh:
				fmt.Fprintf(v.currentView, "%s\n", strings.Join(text, "\n"))
				v.currentView.ScrollToEnd()
			case channel := <-v.joinChannelCh:
				v.JoinChannel(channel.Name)
			case <-v.exitChannelCh:
				v.LeaveChannel()
			}
		}
	}()

}

func (v *app) Run() error {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.currentView, 0, 1, false).
		AddItem(v.inputView, 1, 0, true)
	if err := v.application.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
	return nil
}

func (v *app) JoinChannel(channelName string) {
	v.currentView = v.setUpChannelView(v.application, types2.Channel{Name: channelName})
	f := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.currentView, 0, 1, false).
		AddItem(v.inputView, 1, 0, true)
	v.application.SetRoot(f, true)
}

func (v *app) LeaveChannel() {
	v.renderTextCh <- []string{"heres view, leaving channel"}
	v.currentView = v.lobbyView
	f := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.currentView, 0, 1, false).
		AddItem(v.inputView, 1, 0, true)
	v.application.SetRoot(f, true)
}
