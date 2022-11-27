package gui

import (
	"AlexSarva/GophKeeper/models"
	"fmt"

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

func (gu *GUI) errorModalRender(errorText string, returnPage string) {
	gu.constrains.constrain.ClearButtons()
	gu.constrains.constrain.SetText(fmt.Sprintf("Mistake!\n%s", errorText))
	gu.constrains.constrain.AddButtons([]string{"Cancel"})
	gu.constrains.constrain.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			gu.panels.SetCurrentPanel(returnPage)
		}
	})
	gu.panels.SetCurrentPanel("Mistake")
}

func (gu *GUI) fileHandler(file *models.File) {
	gu.constrains.fileHandler.ClearButtons()
	gu.constrains.fileHandler.SetText(fmt.Sprintf("Enter folder path\nto save file %s", file.FileName))
	gu.constrains.fileHandler.AddButtons([]string{"Enter Folder Path", "Cancel"})
	gu.constrains.fileHandler.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Enter Folder Path" {
			gu.getFileForm(file)
			gu.panels.SetCurrentPanel("GetFile")
			return
		}
		if buttonLabel == "Cancel" {
			gu.panels.SetCurrentPanel("File")
			return
		}
	})
}
