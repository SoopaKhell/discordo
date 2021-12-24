package ui

import (
	"github.com/rivo/tview"
)

func NewMainFlex(app *App) *tview.Flex {
	app.GuildsList.
		ShowSecondaryText(false).
		AddItem("Direct Messages", "", 0, nil).
		SetSelectedFunc(func(idx int, _, _ string, _ rune) {
			onGuildsListSelected(app, idx)
		}).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

	app.ChannelsTreeView.
		SetTopLevel(1).
		SetRoot(tview.NewTreeNode("")).
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 0)

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

	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.GuildsList, 10, 1, false).
		AddItem(app.ChannelsTreeView, 0, 1, false)
	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.MessagesTextView, 0, 1, false).
		AddItem(app.MessageInputField, 3, 1, false)
	mainFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 1, false).
		AddItem(rightFlex, 0, 4, false)

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
