package gui

import (
	"AlexSarva/GophKeeper/models"
	"fmt"
)

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
