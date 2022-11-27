package gui

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/utils"
	"io"
	"os"
	"path"
	"regexp"
)

func (gu *GUI) registerForm() {
	var register models.UserRegister
	gu.forms.registerForm.Clear(true)
	gu.forms.registerForm.AddInputField("Username", "", 20, nil, func(username string) {
		register.Username = username
	})
	gu.forms.registerForm.AddInputField("Email", "", 20, nil, func(email string) {
		register.Email = email
	})
	gu.forms.registerForm.AddPasswordField("Password", "", 20, rune(42), func(passwd string) {
		register.Password = passwd
	})
	gu.forms.registerForm.AddButton("Register", func() {
		user, userErr := gu.client.Register(&register)
		if userErr != nil {
			gu.texts.ChangeAuthText(userErr.Error(), false)
			gu.errorModalRender(userErr.Error(), "Register")
			return
		}
		gu.texts.ChangeAuthText("", true)
		gu.client.UseToken(user.Token)
		gu.collectionContent()
		gu.loggedContent()
		gu.panels.SetCurrentPanel("Collection")
	})
	gu.forms.registerForm.AddButton("Back", func() {
		//pages.SwitchToPage("Menu")
		gu.panels.SetCurrentPanel("Main")
	})
}

func (gu *GUI) loginForm() {
	var login models.UserLogin
	gu.forms.loginForm.Clear(true)
	gu.forms.loginForm.AddInputField("Email", "", 20, nil, func(email string) {
		login.Email = email
	})
	gu.forms.loginForm.AddPasswordField("Password", "", 20, rune(42), func(passwd string) {
		login.Password = passwd
	})
	gu.forms.loginForm.AddButton("Login", func() {
		user, userErr := gu.client.Login(&login)
		if userErr != nil {
			gu.texts.ChangeAuthText(userErr.Error(), false)
			gu.errorModalRender(userErr.Error(), "Main")
			return
		}
		gu.client.UseToken(user.Token)
		gu.texts.ChangeAuthText("", true)
		gu.collectionContent()
		gu.loggedContent()
		gu.panels.SetCurrentPanel("Collection")
	})
	gu.forms.loginForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Main")
	})
}

func (gu *GUI) newNoteForm() {
	var note models.NewNote
	gu.forms.newNoteForm.Clear(true)
	gu.forms.newNoteForm.AddInputField("Title", "", 25, nil, func(title string) {
		note.Title = title
	})

	gu.forms.newNoteForm.AddInputField("Text", "", 35, nil, func(text string) {
		note.Note = text
	})
	gu.forms.newNoteForm.AddButton("Save", func() {
		_, elemErr := gu.client.AddElement("notes", &note)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "NewNote")
			return
		}
		if contentErr := gu.elementsContent("notes"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Notes")
	})
	gu.forms.newNoteForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Notes")
	})
}

func (gu *GUI) editNoteForm(note *models.Note) {
	var editNote models.NewNote
	editNote.Title = note.Title
	editNote.Note = note.Note
	gu.forms.editNoteForm.Clear(true)
	gu.forms.editNoteForm.AddInputField("Title", note.Title, 25, nil, func(title string) {
		editNote.Title = title
	})
	gu.forms.editNoteForm.AddInputField("Text", note.Note, 35, nil, func(text string) {
		editNote.Note = text
	})
	gu.forms.editNoteForm.AddButton("Save", func() {
		_, elemErr := gu.client.EditElement("notes", &editNote, note.ID)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "EditNote")
			return
		}
		if contentErr := gu.elementsContent("notes"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Notes")
	})
	gu.forms.editNoteForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Notes")
	})
}

func (gu *GUI) newCardForm() {
	var card models.NewCard
	gu.forms.newCardForm.Clear(true)
	gu.forms.newCardForm.AddInputField("Title", "", 25, nil, func(title string) {
		card.Title = title
	})
	gu.forms.newCardForm.AddInputField("Card number", "", 35, func(textToCheck string, lastChar rune) bool {
		check := func(text string) int {
			return len(regexp.MustCompile(`[\-\s]`).FindAllString(text, -1))
		}

		if len(textToCheck) > 16+check(textToCheck) {
			return false
		}
		res := utils.CheckCardNumber(textToCheck)
		return res
	}, func(cardNum string) {
		card.CardNumber = cardNum
	})
	gu.forms.newCardForm.AddInputField("Card owner", "", 35, func(textToCheck string, lastChar rune) bool {
		if len(textToCheck) > 20 {
			return false
		}
		res := utils.CheckCardOwner(textToCheck)
		return res
	}, func(owner string) {
		card.CardOwner = owner
	})
	gu.forms.newCardForm.AddInputField("Card exp (mm/dd)", "", 35, func(textToCheck string, lastChar rune) bool {
		check := func(text string) int {
			return len(regexp.MustCompile(`[\/]`).FindAllString(text, -1))
		}

		if len(textToCheck) > 6+check(textToCheck) {
			return false
		}
		res := utils.CheckCardExp(textToCheck)
		return res
	}, func(cardExp string) {
		card.CardExp = cardExp
	})
	gu.forms.newCardForm.AddInputField("Note", "", 35, nil, func(note string) {
		card.Notes = note
	})

	gu.forms.newCardForm.AddButton("Save", func() {
		_, elemErr := gu.client.AddElement("cards", &card)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "NewCard")
			return
		}
		if contentErr := gu.elementsContent("cards"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Cards")
	})
	gu.forms.newCardForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Cards")
	})
}

func (gu *GUI) editCardForm(card *models.Card) {
	var editCard models.NewCard
	gu.forms.editCardForm.Clear(true)
	editCard.CardNumber = card.CardNumber
	editCard.CardExp = card.CardExp
	editCard.CardOwner = card.CardOwner
	editCard.Title = card.Title
	editCard.Notes = card.Notes

	gu.forms.editCardForm.AddInputField("Title", card.Title, 25, nil, func(title string) {
		editCard.Title = title
	})

	gu.forms.editCardForm.AddInputField("Card number", card.CardNumber, 35, func(textToCheck string, lastChar rune) bool {
		if len(textToCheck) > 19 {
			return false
		}
		res := utils.CheckCardNumber(textToCheck)
		return res
	}, func(cardNum string) {
		editCard.CardNumber = cardNum
	})
	gu.forms.editCardForm.AddInputField("Card owner", card.CardOwner, 35, func(textToCheck string, lastChar rune) bool {
		if len(textToCheck) > 20 {
			return false
		}
		res := utils.CheckCardOwner(textToCheck)
		return res
	}, func(cardName string) {
		editCard.CardOwner = cardName
	})
	gu.forms.editCardForm.AddInputField("Card exp (mm/dd)", card.CardExp, 35, func(textToCheck string, lastChar rune) bool {
		if len(textToCheck) > 7 {
			return false
		}
		res := utils.CheckCardExp(textToCheck)
		return res
	}, func(cardExp string) {
		editCard.CardExp = cardExp
	})
	gu.forms.editCardForm.AddInputField("Note", card.Notes, 35, nil, func(note string) {
		editCard.Notes = note
	})
	gu.forms.editCardForm.AddButton("Save", func() {
		_, elemErr := gu.client.EditElement("cards", &editCard, card.ID)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "EditCard")
			return
		}
		if contentErr := gu.elementsContent("cards"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Cards")
	})
	gu.forms.editCardForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Cards")
	})
}

func (gu *GUI) newCredForm() {
	var cred models.NewCred
	gu.forms.newCredForm.Clear(true)
	gu.forms.newCredForm.AddInputField("Title", "", 25, nil, func(title string) {
		cred.Title = title
	})
	gu.forms.newCredForm.AddInputField("Login", "", 35, nil, func(login string) {
		cred.Login = login
	})
	gu.forms.newCredForm.AddInputField("Password", "", 35, nil, func(passwd string) {
		cred.Passwd = passwd
	})
	gu.forms.newCredForm.AddInputField("Note", "", 35, nil, func(note string) {
		cred.Notes = note
	})
	gu.forms.newCredForm.AddButton("Save", func() {
		_, elemErr := gu.client.AddElement("creds", &cred)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "NewCred")
			return
		}
		if contentErr := gu.elementsContent("creds"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Credentials")
	})
	gu.forms.newCredForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Credentials")
	})
}

func (gu *GUI) editCredForm(cred *models.Cred) {
	var editCred models.NewCred
	editCred.Title = cred.Title
	editCred.Login = cred.Login
	editCred.Notes = cred.Notes
	gu.forms.editCredForm.Clear(true)
	gu.forms.editCredForm.AddInputField("Title", cred.Title, 25, nil, func(title string) {
		editCred.Title = title
	})
	gu.forms.editCredForm.AddInputField("Login", cred.Login, 35, nil, func(login string) {
		editCred.Login = login
	})
	gu.forms.editCredForm.AddInputField("Password", cred.Passwd, 35, nil, func(passwd string) {
		editCred.Passwd = passwd
	})
	gu.forms.editCredForm.AddButton("Save", func() {
		_, elemErr := gu.client.EditElement("creds", &editCred, cred.ID)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "NewCred")
			return
		}
		if contentErr := gu.elementsContent("creds"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Credentials")
	})
	gu.forms.editCredForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Credentials")
	})
}

func (gu *GUI) newFileForm() {
	var clientFile models.NewFile
	var filepath string
	gu.forms.newFileForm.Clear(true)
	gu.forms.newFileForm.AddInputField("Title", "", 25, nil, func(title string) {
		clientFile.Title = title
	})
	gu.forms.newFileForm.AddInputField("File name", "", 35, nil, func(name string) {
		clientFile.FileName = name
	})
	gu.forms.newFileForm.AddInputField("Filepath", "", 35, nil, func(path string) {
		filepath = path
	})
	gu.forms.newFileForm.AddInputField("Note", "", 35, nil, func(note string) {
		clientFile.Notes = note
	})

	gu.forms.newFileForm.AddButton("Save", func() {
		file, fileErr := os.Open(filepath)
		if fileErr != nil {
			gu.errorModalRender(fileErr.Error(), "NewFile")
			return
		}
		bodyBytes, readErr := io.ReadAll(file)
		if readErr != nil {
			gu.errorModalRender(readErr.Error(), "NewFile")
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				gu.errorModalRender(err.Error(), "NewFile")
				return
			}
		}(file)

		clientFile.File = bodyBytes

		_, elemErr := gu.client.AddElement("files", &clientFile)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "NewFile")
			return
		}
		if contentErr := gu.elementsContent("files"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Files")
	})
	gu.forms.newFileForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Files")
	})
}

func (gu *GUI) editFileForm(file *models.File) {
	var editClientFile models.NewFile
	var filepath string
	editClientFile.Title = file.Title
	editClientFile.FileName = file.FileName
	editClientFile.File = file.File
	editClientFile.Notes = file.Notes
	gu.forms.editFileForm.Clear(true)
	gu.forms.editFileForm.AddInputField("Title", file.Title, 25, nil, func(title string) {
		editClientFile.Title = title
	})
	gu.forms.editFileForm.AddInputField("File name", file.FileName, 35, nil, func(name string) {
		editClientFile.FileName = name
	})
	gu.forms.editFileForm.AddInputField("Filepath (blank if u dont want to change file)", "", 35, nil, func(path string) {
		filepath = path
	})
	gu.forms.editFileForm.AddInputField("Note", file.Notes, 35, nil, func(note string) {
		editClientFile.Notes = note
	})

	gu.forms.editFileForm.AddButton("Save", func() {
		if filepath != "" {
			tmpfile, fileErr := os.Open(filepath)
			if fileErr != nil {
				gu.errorModalRender(fileErr.Error(), "NewFile")
				return
			}
			bodyBytes, readErr := io.ReadAll(tmpfile)
			if readErr != nil {
				gu.errorModalRender(readErr.Error(), "NewFile")
				return
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					gu.errorModalRender(err.Error(), "NewFile")
					return
				}
			}(tmpfile)

			editClientFile.File = bodyBytes
		} else {
			editClientFile.File = file.File
		}
		_, elemErr := gu.client.EditElement("files", &editClientFile, file.ID)
		if elemErr != nil {
			gu.errorModalRender(elemErr.Error(), "EditFile")
			return
		}
		if contentErr := gu.elementsContent("files"); contentErr != nil {
			gu.errorModalRender(contentErr.Error(), "Collection")
			return
		}
		gu.panels.SetCurrentPanel("Files")
	})
	gu.forms.editFileForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("Files")
	})
}

func (gu *GUI) getFileForm(file *models.File) {
	var filepath string
	var filename string
	var fullFilepath string
	gu.forms.getFileForm.Clear(true)
	gu.forms.getFileForm.AddInputField("Folder", "", 40, nil, func(path string) {
		filepath = path
	})
	gu.forms.getFileForm.AddInputField("Filename", "", 30, nil, func(name string) {
		filename = name
	})

	gu.forms.getFileForm.AddButton("Get", func() {
		if filepath != "" {
			if filename == "" {
				fullFilepath = path.Join(filepath, file.FileName)
			} else {
				fullFilepath = path.Join(filepath, filename)
			}
		}
		if err := os.WriteFile(fullFilepath, file.File, 0666); err != nil {
			gu.errorModalRender(err.Error(), "GetFile")
			return
		}
		gu.panels.SetCurrentPanel("File")
	})
	gu.forms.getFileForm.AddButton("Back", func() {
		gu.panels.SetCurrentPanel("File")
	})
}
