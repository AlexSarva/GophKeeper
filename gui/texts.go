package gui

import "github.com/gdamore/tcell/v2"

func (gu *GUI) authMessage() {
	gu.texts.authText.SetTextColor(tcell.ColorRed)
	gu.texts.authText.SetTextAlign(1)
	gu.texts.authText.SetText("You are not logged in!")
}
