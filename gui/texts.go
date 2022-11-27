package gui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

func textPrimitive(text string, color tcell.Color, align int) cview.Primitive {
	tv := cview.NewTextView()
	tv.SetTextColor(color)
	tv.SetText(text)
	tv.SetTextAlign(align)
	tv.Blur()
	return tv
}

func aboutMessage() cview.Primitive {
	tv := cview.NewTextView()
	tv.SetTextColor(tcell.ColorOrangeRed)
	tv.SetText("Добро пожаловать в GophKeeper! \n\n Сервис представляет собой клиент-серверную систему,\nпозволяющую пользователю надёжно и безопасно хранить логины,\nпароли, бинарные данные и прочую приватную информацию.")
	tv.SetTextAlign(1)
	tv.SetPadding(1, 0, 3, 3)
	tv.Blur()
	return tv
}

func elementTextPrimitive(text string) cview.Primitive {
	tv := cview.NewTextView()
	tv.SetTextColor(tcell.ColorWhiteSmoke)
	tv.SetWordWrap(true)
	tv.SetScrollable(true)
	tv.SetText(text)
	tv.SetTextAlign(1)
	tv.SetPadding(1, 0, 2, 2)
	return tv
}

func (gu *GUI) authMessage() {
	gu.texts.authText.SetTextColor(tcell.ColorRed)
	gu.texts.authText.SetTextAlign(1)
	gu.texts.authText.SetText("You are not logged in!")
}
