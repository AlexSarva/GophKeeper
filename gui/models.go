package gui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type Constrains struct {
	constrain   *cview.Modal
	fileHandler *cview.Modal
}

func InitConstrains() *Constrains {
	constrain := cview.NewModal()
	fileHandler := cview.NewModal()
	return &Constrains{
		constrain:   constrain,
		fileHandler: fileHandler,
	}
}

type Layouts struct {
	mainPage       *cview.Grid
	collectionPage *cview.Grid
	notesPage      *cview.Grid
	elementPage    *cview.Grid
	cardsPage      *cview.Grid
	filesPage      *cview.Grid
	credsPage      *cview.Grid
}

func InitLayouts() *Layouts {
	mainGrid := cview.NewGrid()
	mainGrid.SetColumns(45, 0)
	mainGrid.SetRows(1, 1, 0)
	mainGrid.SetBorders(true)
	mainGrid.SetGap(1, 0)
	mainGrid.AddItem(textPrimitive("Welcome", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)

	collectionGrid := cview.NewGrid()
	collectionGrid.SetColumns(55, 0)
	collectionGrid.SetRows(1, 1, 0)
	collectionGrid.SetBorders(true)
	collectionGrid.SetGap(1, 0)
	collectionGrid.AddItem(textPrimitive("Collection: ", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)

	cardsGrid := cview.NewGrid()
	cardsGrid.SetColumns(60, 0)
	cardsGrid.SetRows(1, 1, 0)
	cardsGrid.SetBorders(true)
	cardsGrid.SetGap(1, 0)
	cardsGrid.AddItem(textPrimitive("Cards: ", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)

	credsGrid := cview.NewGrid()
	credsGrid.SetColumns(60, 0)
	credsGrid.SetRows(1, 1, 0)
	credsGrid.SetBorders(true)
	credsGrid.SetGap(1, 0)
	credsGrid.AddItem(textPrimitive("Credentials: ", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)

	notesGrid := cview.NewGrid()
	notesGrid.SetColumns(60, 0)
	notesGrid.SetRows(1, 1, 0)
	notesGrid.SetBorders(true)
	notesGrid.SetGap(1, 0)
	notesGrid.AddItem(textPrimitive("Notes: ", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)

	elementGrid := cview.NewGrid()
	elementGrid.SetColumns(40, 0)
	elementGrid.SetRows(1, 0, 1)
	elementGrid.SetBorders(true)
	elementGrid.SetGap(1, 0)

	filesGrid := cview.NewGrid()
	filesGrid.SetColumns(60, 0)
	filesGrid.SetRows(1, 1, 0)
	filesGrid.SetBorders(true)
	filesGrid.SetGap(1, 0)
	filesGrid.AddItem(textPrimitive("Files: ", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)

	return &Layouts{
		mainPage:       mainGrid,
		collectionPage: collectionGrid,
		notesPage:      notesGrid,
		elementPage:    elementGrid,
		cardsPage:      cardsGrid,
		filesPage:      filesGrid,
		credsPage:      credsGrid,
	}
}

type Content struct {
	welcomeContent     *cview.List
	collectionContent  *cview.List
	notesContent       *cview.List
	elementMenuContent *cview.List
	cardsContent       *cview.List
	credsContent       *cview.List
	filesContent       *cview.List
}

func InitContent() *Content {
	welcomeContent := cview.NewList()
	collectionContent := cview.NewList()
	notesContent := cview.NewList()
	elementMenuContent := cview.NewList()
	cardsContent := cview.NewList()
	credsContent := cview.NewList()
	filesContent := cview.NewList()
	return &Content{
		welcomeContent:     welcomeContent,
		collectionContent:  collectionContent,
		notesContent:       notesContent,
		elementMenuContent: elementMenuContent,
		cardsContent:       cardsContent,
		credsContent:       credsContent,
		filesContent:       filesContent,
	}
}

type Forms struct {
	registerForm *cview.Form
	loginForm    *cview.Form
	newNoteForm  *cview.Form
	editNoteForm *cview.Form
	newCardForm  *cview.Form
	editCardForm *cview.Form
	newCredForm  *cview.Form
	editCredForm *cview.Form
	newFileForm  *cview.Form
	editFileForm *cview.Form
	getFileForm  *cview.Form
}

func InitForms() *Forms {
	registerForm := cview.NewForm()
	loginForm := cview.NewForm()
	newNoteForm := cview.NewForm()
	editNoteForm := cview.NewForm()
	newCardForm := cview.NewForm()
	editCardForm := cview.NewForm()
	newCredForm := cview.NewForm()
	editCredForm := cview.NewForm()
	newFileForm := cview.NewForm()
	editFileForm := cview.NewForm()
	getFileForm := cview.NewForm()
	return &Forms{
		registerForm: registerForm,
		loginForm:    loginForm,
		newNoteForm:  newNoteForm,
		editNoteForm: editNoteForm,
		newCardForm:  newCardForm,
		editCardForm: editCardForm,
		newCredForm:  newCredForm,
		editCredForm: editCredForm,
		newFileForm:  newFileForm,
		editFileForm: editFileForm,
		getFileForm:  getFileForm,
	}
}

type Texts struct {
	authText *cview.TextView
}

func InitTexts() *Texts {
	authText := cview.NewTextView()
	return &Texts{
		authText: authText,
	}
}

func (t *Texts) ChangeAuthText(text string, auth bool) {
	if !auth {
		t.authText.SetTextColor(tcell.ColorRed)
		t.authText.SetText(text)
		return
	}
	t.authText.SetTextColor(tcell.ColorGreen)
	t.authText.SetText("You are successfully logged in!")

}
