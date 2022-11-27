package gui

import (
	"AlexSarva/GophKeeper/models"
	"fmt"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

func (gu *GUI) welcomeContent() {
	gu.content.welcomeContent.Clear()
	registerItem := cview.NewListItem("Register")
	registerItem.SetSecondaryText("Press to register")
	registerItem.SetShortcut('1')
	registerItem.SetSelectedFunc(func() {
		gu.registerForm()
		gu.panels.SetCurrentPanel("Register")
	})

	loginItem := cview.NewListItem("Login")
	loginItem.SetSecondaryText("Press to login")
	loginItem.SetShortcut('2')
	loginItem.SetSelectedFunc(func() {
		gu.loginForm()
		gu.panels.SetCurrentPanel("Login")
	})

	quitItem := cview.NewListItem("Quit")
	quitItem.SetSecondaryText("Press to exit")
	quitItem.SetShortcut('q')
	quitItem.SetSelectedFunc(func() {
		gu.app.Stop()
	})

	gu.content.welcomeContent.AddItem(registerItem)
	gu.content.welcomeContent.AddItem(loginItem)
	gu.content.welcomeContent.AddItem(quitItem)

	gu.content.welcomeContent.SetPadding(0, 0, 2, 0)

}

func (gu *GUI) loggedContent() {
	gu.content.welcomeContent.Clear()
	collectItem := cview.NewListItem("Collection")
	collectItem.SetSecondaryText("Go to collection")
	collectItem.SetShortcut('1')
	collectItem.SetSelectedFunc(func() {
		gu.panels.SetCurrentPanel("Collection")
	})

	logoutItem := cview.NewListItem("Log Out")
	logoutItem.SetSecondaryText("Press to log out")
	logoutItem.SetShortcut('2')
	logoutItem.SetSelectedFunc(func() {
		gu.welcomeContent()
		gu.texts.ChangeAuthText("You are not logged in!", false)
		//gu.client.UseToken("")
		gu.panels.SetCurrentPanel("Main")
	})

	quitItem := cview.NewListItem("Quit")
	quitItem.SetSecondaryText("Press to exit")
	quitItem.SetShortcut('q')
	quitItem.SetSelectedFunc(func() {
		gu.app.Stop()
	})

	gu.content.welcomeContent.AddItem(collectItem)
	gu.content.welcomeContent.AddItem(logoutItem)
	gu.content.welcomeContent.AddItem(quitItem)

	//welcome.ContextMenuList().SetItemEnabled(2, false)

	gu.content.welcomeContent.SetPadding(0, 0, 2, 0)

}

func (gu *GUI) collectionContent() {
	gu.content.collectionContent.Clear()

	cards := cview.NewListItem("Cards")
	cards.SetSecondaryText("Go to cards")
	cards.SetShortcut('1')
	cards.SetSelectedFunc(func() {
		if elemErr := gu.elementsContent("cards"); elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Cards")
	})

	creds := cview.NewListItem("Credentials")
	creds.SetSecondaryText("Go to credentials")
	creds.SetShortcut('2')
	creds.SetSelectedFunc(func() {
		if elemErr := gu.elementsContent("creds"); elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Credentials")
	})

	notes := cview.NewListItem("Notes")
	notes.SetSecondaryText("Go to notes")
	notes.SetShortcut('3')
	notes.SetSelectedFunc(func() {
		if elemErr := gu.elementsContent("notes"); elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Notes")
	})

	files := cview.NewListItem("Files")
	files.SetSecondaryText("Go to files")
	files.SetShortcut('4')
	files.SetSelectedFunc(func() {
		if elemErr := gu.elementsContent("files"); elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Files")
	})

	quitItem := cview.NewListItem("Back")
	quitItem.SetSecondaryText("Go to Main page")
	quitItem.SetShortcut('b')
	quitItem.SetSelectedFunc(func() {
		gu.panels.SetCurrentPanel("Main")
	})

	emptyItem := cview.NewListItem("")
	gu.content.collectionContent.AddItem(cards)
	gu.content.collectionContent.AddItem(creds)
	gu.content.collectionContent.AddItem(notes)
	gu.content.collectionContent.AddItem(files)
	gu.content.collectionContent.AddItem(emptyItem)
	gu.content.collectionContent.AddItem(emptyItem)
	gu.content.collectionContent.AddItem(quitItem)

}

func (gu *GUI) elementsContent(infoType string) error {
	elems, elemsErr := gu.client.ElementList(infoType)
	if elemsErr != nil {
		return elemsErr
	}

	switch infoType {
	case "notes":
		el := elems.([]models.Note)
		gu.content.notesContent.Clear()

		if len(el) != 0 {
			for index, value := range el {
				item := cview.NewListItem(value.Title)
				date := value.Created.Format("02 Jan 2006 15:04:05")
				if value.Changed != nil {
					date = fmt.Sprintf("%s (created: %s)", value.Changed.Time.Format("02 Jan 2006 15:04:05"), date)
				}
				item.SetSecondaryText(date)
				item.SetShortcut(rune(49 + index))
				gu.content.notesContent.AddItem(item)
			}
		} else {
			noContentItem := cview.NewListItem("No content")
			noContentItem.SetSecondaryText("no content in database")
			noContentItem.SetShortcut('x')
			gu.content.notesContent.AddItem(noContentItem)
		}

		emptyItem := cview.NewListItem("")

		newItem := cview.NewListItem("New Note")
		newItem.SetSecondaryText("crete New Note")
		newItem.SetShortcut('n')
		newItem.SetSelectedFunc(func() {
			gu.newNoteForm()
			gu.panels.SetCurrentPanel("NewNote")
		})

		colItem := cview.NewListItem("To Collection")
		colItem.SetSecondaryText("Go to collection")
		colItem.SetShortcut('c')
		colItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Collection")
		})

		quitItem := cview.NewListItem("To Main")
		quitItem.SetSecondaryText("Go to main menu")
		quitItem.SetShortcut('m')
		quitItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Main")
		})

		gu.content.notesContent.AddItem(emptyItem)
		gu.content.notesContent.AddItem(emptyItem)
		gu.content.notesContent.AddItem(newItem)
		gu.content.notesContent.AddItem(colItem)
		gu.content.notesContent.AddItem(quitItem)

		gu.content.notesContent.SetSelectedFunc(func(index int, element *cview.ListItem) {
			if index < len(el) {
				gu.generateNote(&el[index])
				gu.panels.SetCurrentPanel("Note")
			}
		})

		return nil
	case "cards":
		el := elems.([]models.Card)
		gu.content.cardsContent.Clear()

		if len(el) != 0 {
			for index, value := range el {
				item := cview.NewListItem(value.Title)
				date := value.Created.Format("02 Jan 2006 15:04:05")
				if value.Changed != nil {
					date = fmt.Sprintf("%s (created: %s)", value.Changed.Time.Format("02 Jan 2006 15:04:05"), date)
				}
				item.SetSecondaryText(date)
				item.SetShortcut(rune(49 + index))
				gu.content.cardsContent.AddItem(item)
			}
		} else {
			noContentItem := cview.NewListItem("No content")
			noContentItem.SetSecondaryText("no content in database")
			noContentItem.SetShortcut('x')
			gu.content.cardsContent.AddItem(noContentItem)
		}

		emptyItem := cview.NewListItem("")

		newItem := cview.NewListItem("New Card")
		newItem.SetSecondaryText("crete New Card")
		newItem.SetShortcut('n')
		newItem.SetSelectedFunc(func() {
			gu.newCardForm()
			gu.panels.SetCurrentPanel("NewCard")
		})

		colItem := cview.NewListItem("To Collection")
		colItem.SetSecondaryText("Go to collection")
		colItem.SetShortcut('c')
		colItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Collection")
		})

		quitItem := cview.NewListItem("To Main")
		quitItem.SetSecondaryText("Go to main menu")
		quitItem.SetShortcut('m')
		quitItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Main")
		})

		gu.content.cardsContent.AddItem(emptyItem)
		gu.content.cardsContent.AddItem(emptyItem)
		gu.content.cardsContent.AddItem(newItem)
		gu.content.cardsContent.AddItem(colItem)
		gu.content.cardsContent.AddItem(quitItem)

		gu.content.cardsContent.SetSelectedFunc(func(index int, element *cview.ListItem) {
			if index < len(el) {
				gu.generateCard(&el[index])
				gu.panels.SetCurrentPanel("Card")
			}
		})

		return nil
	case "creds":
		el := elems.([]models.Cred)
		gu.content.credsContent.Clear()

		if len(el) != 0 {
			for index, value := range el {
				item := cview.NewListItem(value.Title)
				date := value.Created.Format("02 Jan 2006 15:04:05")
				if value.Changed != nil {
					date = fmt.Sprintf("%s (created: %s)", value.Changed.Time.Format("02 Jan 2006 15:04:05"), date)
				}
				item.SetSecondaryText(date)
				item.SetShortcut(rune(49 + index))
				gu.content.credsContent.AddItem(item)
			}
		} else {
			noContentItem := cview.NewListItem("No content")
			noContentItem.SetSecondaryText("no content in database")
			noContentItem.SetShortcut('x')
			gu.content.credsContent.AddItem(noContentItem)
		}

		emptyItem := cview.NewListItem("")

		newItem := cview.NewListItem("New Cred")
		newItem.SetSecondaryText("crete New Cred")
		newItem.SetShortcut('n')
		newItem.SetSelectedFunc(func() {
			gu.newCredForm()
			gu.panels.SetCurrentPanel("NewCred")
		})

		colItem := cview.NewListItem("To Collection")
		colItem.SetSecondaryText("Go to collection")
		colItem.SetShortcut('c')
		colItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Collection")
		})

		quitItem := cview.NewListItem("To Main")
		quitItem.SetSecondaryText("Go to main menu")
		quitItem.SetShortcut('m')
		quitItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Main")
		})

		gu.content.credsContent.AddItem(emptyItem)
		gu.content.credsContent.AddItem(emptyItem)
		gu.content.credsContent.AddItem(newItem)
		gu.content.credsContent.AddItem(colItem)
		gu.content.credsContent.AddItem(quitItem)

		gu.content.credsContent.SetSelectedFunc(func(index int, element *cview.ListItem) {
			if index < len(el) {
				gu.generateCred(&el[index])
				gu.panels.SetCurrentPanel("Cred")
			}
		})
		return nil
	case "files":
		el := elems.([]models.File)
		gu.content.filesContent.Clear()

		if len(el) != 0 {
			for index, value := range el {
				item := cview.NewListItem(value.Title)
				date := value.Created.Format("02 Jan 2006 15:04:05")
				if value.Changed != nil {
					date = fmt.Sprintf("%s (created: %s)", value.Changed.Time.Format("02 Jan 2006 15:04:05"), date)
				}
				item.SetSecondaryText(date)
				item.SetShortcut(rune(49 + index))
				gu.content.filesContent.AddItem(item)
			}
		} else {
			noContentItem := cview.NewListItem("No content")
			noContentItem.SetSecondaryText("no content in database")
			noContentItem.SetShortcut('x')
			gu.content.filesContent.AddItem(noContentItem)
		}

		emptyItem := cview.NewListItem("")

		newItem := cview.NewListItem("New File")
		newItem.SetSecondaryText("crete New File")
		newItem.SetShortcut('n')
		newItem.SetSelectedFunc(func() {
			gu.newFileForm()
			gu.panels.SetCurrentPanel("NewFile")
		})

		colItem := cview.NewListItem("To Collection")
		colItem.SetSecondaryText("Go to collection")
		colItem.SetShortcut('c')
		colItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Collection")
		})

		quitItem := cview.NewListItem("To Main")
		quitItem.SetSecondaryText("Go to main menu")
		quitItem.SetShortcut('m')
		quitItem.SetSelectedFunc(func() {
			gu.panels.SetCurrentPanel("Main")
		})

		gu.content.filesContent.AddItem(emptyItem)
		gu.content.filesContent.AddItem(emptyItem)
		gu.content.filesContent.AddItem(newItem)
		gu.content.filesContent.AddItem(colItem)
		gu.content.filesContent.AddItem(quitItem)

		gu.content.filesContent.SetSelectedFunc(func(index int, element *cview.ListItem) {
			if index < len(el) {
				gu.generateFile(&el[index])
				gu.panels.SetCurrentPanel("File")
			}
		})

		return nil
	}
	return nil
}

func (gu *GUI) generateNote(note *models.Note) {
	gu.layouts.elementPage.Clear()
	gu.content.elementMenuContent.Clear()

	date := fmt.Sprintf("Created: %s", note.Created.Format("02 Jan 2006 15:04:05"))
	if note.Changed != nil {
		date = fmt.Sprintf("%s, Changed: %s", date, note.Changed.Time.Format("02 Jan 2006 15:04:05"))
	}

	editItem := cview.NewListItem("Edit")
	editItem.SetSecondaryText("edit this Note")
	editItem.SetShortcut('e')
	editItem.SetSelectedFunc(func() {
		gu.editNoteForm(note)
		gu.panels.SetCurrentPanel("EditNote")
	})

	deleteItem := cview.NewListItem("Delete")
	deleteItem.SetSecondaryText("delete this Note")
	deleteItem.SetShortcut('d')
	deleteItem.SetSelectedFunc(func() {
		_, delErr := gu.client.Delete("notes", note.ID)
		if delErr != nil {
			gu.errorModalRender(delErr.Error(), "Notes")
			return
		}
		if contentErr := gu.elementsContent("notes"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Notes")
	})

	backItem := cview.NewListItem("To Notes")
	backItem.SetSecondaryText("Go to Notes")
	backItem.SetShortcut('b')
	backItem.SetSelectedFunc(func() {
		gu.panels.SetCurrentPanel("Notes")
	})

	gu.content.elementMenuContent.AddItem(editItem)
	gu.content.elementMenuContent.AddItem(deleteItem)
	gu.content.elementMenuContent.AddItem(backItem)
	gu.content.elementMenuContent.SetPadding(1, 0, 2, 0)

	gu.layouts.elementPage.AddItem(textPrimitive(note.Title, tcell.ColorKhaki, 1), 0, 0, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(gu.content.elementMenuContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive("ID: "+note.ID.String(), tcell.ColorDarkSalmon, 1), 0, 1, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(elementTextPrimitive("Text: "+note.Note), 1, 1, 1, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive(date, tcell.ColorDarkOrange, 1), 2, 1, 1, 1, 0, 0, false)
}

func (gu *GUI) generateCard(card *models.Card) {
	gu.layouts.elementPage.Clear()
	gu.content.elementMenuContent.Clear()

	date := fmt.Sprintf("Created: %s", card.Created.Format("02 Jan 2006 15:04:05"))
	if card.Changed != nil {
		date = fmt.Sprintf("%s, Changed: %s", date, card.Changed.Time.Format("02 Jan 2006 15:04:05"))
	}

	text := fmt.Sprintf("Card number: %s\nCard owner: %s\nCard exp: %s", card.CardNumber, card.CardOwner, card.CardExp)
	if card.Notes != "" {
		text = fmt.Sprintf("%s\n\n%s", text, card.Notes)
	}

	editItem := cview.NewListItem("Edit")
	editItem.SetSecondaryText("edit this Card")
	editItem.SetShortcut('e')
	editItem.SetSelectedFunc(func() {
		gu.editCardForm(card)
		gu.panels.SetCurrentPanel("EditCard")
	})

	deleteItem := cview.NewListItem("Delete")
	deleteItem.SetSecondaryText("delete this Card")
	deleteItem.SetShortcut('d')
	deleteItem.SetSelectedFunc(func() {
		_, delErr := gu.client.Delete("cards", card.ID)
		if delErr != nil {
			gu.errorModalRender(delErr.Error(), "Cards")
			return
		}
		if contentErr := gu.elementsContent("cards"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Cards")
	})

	backItem := cview.NewListItem("To Cards")
	backItem.SetSecondaryText("Go to Cards")
	backItem.SetShortcut('b')
	backItem.SetSelectedFunc(func() {
		gu.panels.SetCurrentPanel("Cards")
	})

	gu.content.elementMenuContent.AddItem(editItem)
	gu.content.elementMenuContent.AddItem(deleteItem)
	gu.content.elementMenuContent.AddItem(backItem)
	gu.content.elementMenuContent.SetPadding(1, 0, 2, 0)

	gu.layouts.elementPage.AddItem(textPrimitive(card.Title, tcell.ColorKhaki, 1), 0, 0, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(gu.content.elementMenuContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive("ID: "+card.ID.String(), tcell.ColorDarkSalmon, 1), 0, 1, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(elementTextPrimitive(text), 1, 1, 1, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive(date, tcell.ColorDarkOrange, 1), 2, 1, 1, 1, 0, 0, false)
}

func (gu *GUI) generateCred(cred *models.Cred) {
	gu.layouts.elementPage.Clear()
	gu.content.elementMenuContent.Clear()

	date := fmt.Sprintf("Created: %s", cred.Created.Format("02 Jan 2006 15:04:05"))
	if cred.Changed != nil {
		date = fmt.Sprintf("%s, Changed: %s", date, cred.Changed.Time.Format("02 Jan 2006 15:04:05"))
	}

	text := fmt.Sprintf("Login: %s\nPassword: %s", cred.Login, cred.Passwd)
	if cred.Notes != "" {
		text = fmt.Sprintf("%s\n\n%s", text, cred.Notes)
	}

	editItem := cview.NewListItem("Edit")
	editItem.SetSecondaryText("edit this Cred")
	editItem.SetShortcut('e')
	editItem.SetSelectedFunc(func() {
		gu.editCredForm(cred)
		gu.panels.SetCurrentPanel("EditCred")
	})

	deleteItem := cview.NewListItem("Delete")
	deleteItem.SetSecondaryText("delete this Cred")
	deleteItem.SetShortcut('d')
	deleteItem.SetSelectedFunc(func() {
		_, delErr := gu.client.Delete("creds", cred.ID)
		if delErr != nil {
			gu.errorModalRender(delErr.Error(), "Credentials")
			return
		}
		if contentErr := gu.elementsContent("creds"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Credentials")
	})

	backItem := cview.NewListItem("To Credentials")
	backItem.SetSecondaryText("Go to Credentials")
	backItem.SetShortcut('b')
	backItem.SetSelectedFunc(func() {
		gu.panels.SetCurrentPanel("Credentials")
	})

	gu.content.elementMenuContent.AddItem(editItem)
	gu.content.elementMenuContent.AddItem(deleteItem)
	gu.content.elementMenuContent.AddItem(backItem)
	gu.content.elementMenuContent.SetPadding(1, 0, 2, 0)

	gu.layouts.elementPage.AddItem(textPrimitive(cred.Title, tcell.ColorKhaki, 1), 0, 0, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(gu.content.elementMenuContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive("ID: "+cred.ID.String(), tcell.ColorDarkSalmon, 1), 0, 1, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(elementTextPrimitive(text), 1, 1, 1, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive(date, tcell.ColorDarkOrange, 1), 2, 1, 1, 1, 0, 0, false)
}

func (gu *GUI) generateFile(file *models.File) {
	gu.layouts.elementPage.Clear()
	gu.content.elementMenuContent.Clear()

	date := fmt.Sprintf("Created: %s", file.Created.Format("02 Jan 2006 15:04:05"))
	if file.Changed != nil {
		date = fmt.Sprintf("%s, Changed: %s", date, file.Changed.Time.Format("02 Jan 2006 15:04:05"))
	}

	text := fmt.Sprintf("File name: %s", file.FileName)
	if file.Notes != "" {
		text = fmt.Sprintf("%s\n\n%s", text, file.Notes)
	}

	getItem := cview.NewListItem("Get")
	getItem.SetSecondaryText("Get this File")
	getItem.SetShortcut('e')
	getItem.SetSelectedFunc(func() {
		gu.fileHandler(file)
		gu.panels.SetCurrentPanel("FileHandler")
	})

	editItem := cview.NewListItem("Edit")
	editItem.SetSecondaryText("edit this File")
	editItem.SetShortcut('e')
	editItem.SetSelectedFunc(func() {
		gu.editFileForm(file)
		gu.panels.SetCurrentPanel("EditFile")
	})

	deleteItem := cview.NewListItem("Delete")
	deleteItem.SetSecondaryText("delete this File")
	deleteItem.SetShortcut('d')
	deleteItem.SetSelectedFunc(func() {
		_, delErr := gu.client.Delete("files", file.ID)
		if delErr != nil {
			gu.errorModalRender(delErr.Error(), "Files")
			return
		}
		if contentErr := gu.elementsContent("files"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Files")
			return
		}
		gu.panels.SetCurrentPanel("Files")
	})

	backItem := cview.NewListItem("To Files")
	backItem.SetSecondaryText("Go to Files")
	backItem.SetShortcut('b')
	backItem.SetSelectedFunc(func() {
		gu.panels.SetCurrentPanel("Files")
	})

	gu.content.elementMenuContent.AddItem(getItem)
	gu.content.elementMenuContent.AddItem(editItem)
	gu.content.elementMenuContent.AddItem(deleteItem)
	gu.content.elementMenuContent.AddItem(backItem)
	gu.content.elementMenuContent.SetPadding(1, 0, 2, 0)

	gu.layouts.elementPage.AddItem(textPrimitive(file.Title, tcell.ColorKhaki, 1), 0, 0, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(gu.content.elementMenuContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive("ID: "+file.ID.String(), tcell.ColorDarkSalmon, 1), 0, 1, 1, 1, 0, 0, false)
	gu.layouts.elementPage.AddItem(elementTextPrimitive(text), 1, 1, 1, 1, 0, 0, true)
	gu.layouts.elementPage.AddItem(textPrimitive(date, tcell.ColorDarkOrange, 1), 2, 1, 1, 1, 0, 0, false)
}
