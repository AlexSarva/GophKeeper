package gui

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/workclient"
	"log"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

// GUI graphical user interface, it allows users to interact with
// GophKeeper service through graphical icons
type GUI struct {
	app        *cview.Application
	client     *workclient.Client
	panels     *cview.Panels
	layouts    *layouts
	content    *content
	forms      *forms
	texts      *texts
	constrains *constrains
}

// InitGUI initialize GUI, cfg should provide information about service address,
// crypto keys location, size, and secret for ecnrypt files
func InitGUI(cfg *models.GUIConfig) *GUI {
	app := cview.NewApplication()
	defer app.HandlePanic()
	app.EnableMouse(true)
	cli, cliErr := workclient.InitClient(cfg)
	if cliErr != nil {
		log.Fatalln(cliErr)
	}
	workLayouts := initLayouts()
	workPanels := cview.NewPanels()
	workContent := initContent()
	workForms := initForms()
	workTexts := initTexts()
	workConstrains := initConstrains()
	return &GUI{
		app:        app,
		client:     cli,
		layouts:    workLayouts,
		panels:     workPanels,
		content:    workContent,
		forms:      workForms,
		texts:      workTexts,
		constrains: workConstrains,
	}
}

// Render provides basic pages, layouts, modals, texts and forms
func (gu *GUI) Render() {
	gu.welcomeContent()
	gu.authMessage()
	// main page
	gu.layouts.mainPage.AddItem(textPrimitive("Welcome", tcell.ColorBlue, 1), 0, 0, 1, 1, 0, 0, false)
	gu.layouts.mainPage.AddItem(gu.content.welcomeContent, 2, 0, 1, 1, 0, 0, true)
	gu.layouts.mainPage.AddItem(gu.texts.authText, 1, 0, 1, 1, 0, 0, true)
	gu.layouts.mainPage.AddItem(aboutMessage(), 0, 1, 3, 1, 0, 0, false)

	// collection page
	gu.layouts.collectionPage.AddItem(gu.content.collectionContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.collectionPage.AddItem(textPrimitive("", tcell.ColorBlue, 1), 0, 1, 3, 1, 0, 0, false)

	// notes page
	gu.layouts.notesPage.AddItem(gu.content.notesContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.notesPage.AddItem(textPrimitive("", tcell.ColorBlue, 1), 0, 1, 3, 1, 0, 0, false)

	// creds page
	gu.layouts.credsPage.AddItem(gu.content.credsContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.credsPage.AddItem(textPrimitive("", tcell.ColorBlue, 1), 0, 1, 3, 1, 0, 0, false)

	// cards page
	gu.layouts.cardsPage.AddItem(gu.content.cardsContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.cardsPage.AddItem(textPrimitive("", tcell.ColorBlue, 1), 0, 1, 3, 1, 0, 0, false)

	// files page
	gu.layouts.filesPage.AddItem(gu.content.filesContent, 1, 0, 2, 1, 0, 0, true)
	gu.layouts.filesPage.AddItem(textPrimitive("", tcell.ColorBlue, 1), 0, 1, 3, 1, 0, 0, false)

	gu.panels.AddPanel("Main", gu.layouts.mainPage, true, true)
	gu.panels.AddPanel("Register", gu.forms.registerForm, true, false)
	gu.panels.AddPanel("Login", gu.forms.loginForm, true, false)
	gu.panels.AddPanel("NewNote", gu.forms.newNoteForm, true, false)
	gu.panels.AddPanel("EditNote", gu.forms.editNoteForm, true, false)
	gu.panels.AddPanel("NewCard", gu.forms.newCardForm, true, false)
	gu.panels.AddPanel("EditCard", gu.forms.editCardForm, true, false)
	gu.panels.AddPanel("NewCred", gu.forms.newCredForm, true, false)
	gu.panels.AddPanel("EditCred", gu.forms.editCredForm, true, false)
	gu.panels.AddPanel("NewFile", gu.forms.newFileForm, true, false)
	gu.panels.AddPanel("EditFile", gu.forms.editFileForm, true, false)
	gu.panels.AddPanel("Collection", gu.layouts.collectionPage, true, false)
	gu.panels.AddPanel("Notes", gu.layouts.notesPage, true, false)
	gu.panels.AddPanel("Cards", gu.layouts.cardsPage, true, false)
	gu.panels.AddPanel("Credentials", gu.layouts.credsPage, true, false)
	gu.panels.AddPanel("Files", gu.layouts.filesPage, true, false)
	gu.panels.AddPanel("Note", gu.layouts.elementPage, true, false)
	gu.panels.AddPanel("File", gu.layouts.elementPage, true, false)
	gu.panels.AddPanel("Card", gu.layouts.elementPage, true, false)
	gu.panels.AddPanel("Cred", gu.layouts.elementPage, true, false)
	gu.panels.AddPanel("Mistake", gu.constrains.constrain, false, false)
	gu.panels.AddPanel("FileHandler", gu.constrains.fileHandler, false, false)
	gu.panels.AddPanel("GetFile", gu.forms.getFileForm, true, false)
}

// Run starts the GUI
func (gu *GUI) Run() error {
	gu.app.SetRoot(gu.panels, true)

	if err := gu.app.Run(); err != nil {
		return err
	}

	return nil
}
