package ui

import (
	"github.com/rivo/tview"
)

func NewMainFlex(app *App) *tview.Flex {
	mainFlex := tview.NewFlex()

	app.ChannelsTreeView.
		SetTopLevel(1).
		SetRoot(tview.NewTreeNode("")).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)
	mainFlex.AddItem(app.ChannelsTreeView, 0, 1, false)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	app.MessagesTextView.
		SetRegions(true).
		SetDynamicColors(true).
		SetWordWrap(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	app.MessageInputField.
		SetPlaceholder("Message...").
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	rightFlex.
		AddItem(app.MessagesTextView, 0, 1, false).
		AddItem(app.MessageInputField, 3, 1, false)
	mainFlex.AddItem(rightFlex, 0, 4, false)

	return mainFlex
}

func NewLoginForm() *tview.Form {
	loginForm := tview.NewForm()
	loginForm.
		SetButtonsAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	loginForm.
		AddInputField("Email", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, 0, nil).
		AddPasswordField("MFA Code (optional)", "", 0, 0, nil)

	return loginForm
}
