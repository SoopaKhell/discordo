package main

import (
	"github.com/ayntgl/discordgo"
	"github.com/ayntgl/discordo/ui"
	"github.com/ayntgl/discordo/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zalando/go-keyring"
)

const service = "discordo"

var (
	app               *tview.Application
	loginForm         *tview.Form
	channelsTreeView  *tview.TreeView
	messagesTextView  *tview.TextView
	messageInputField *tview.InputField
	mainFlex          *tview.Flex

	conf    *util.Config
	session *discordgo.Session

	selectedChannel *discordgo.Channel
	selectedMessage int = -1
)

func main() {
	conf = util.LoadConfig()

	tview.Borders.Horizontal = conf.Borders.Horizontal
	tview.Borders.Vertical = conf.Borders.Vertical
	tview.Borders.TopLeft = conf.Borders.TopLeft
	tview.Borders.TopRight = conf.Borders.TopRight
	tview.Borders.BottomLeft = conf.Borders.BottomLeft
	tview.Borders.BottomRight = conf.Borders.BottomRight
	tview.Borders.HorizontalFocus = conf.Borders.HorizontalFocus
	tview.Borders.VerticalFocus = conf.Borders.VerticalFocus
	tview.Borders.TopLeftFocus = conf.Borders.TopLeftFocus
	tview.Borders.TopRightFocus = conf.Borders.TopRightFocus
	tview.Borders.BottomLeftFocus = conf.Borders.BottomLeftFocus
	tview.Borders.BottomRightFocus = conf.Borders.BottomRightFocus

	tview.Styles.PrimitiveBackgroundColor = tcell.GetColor(conf.Theme.Background)
	tview.Styles.ContrastBackgroundColor = tcell.GetColor(conf.Theme.Background)
	tview.Styles.MoreContrastBackgroundColor = tcell.GetColor(conf.Theme.Background)
	tview.Styles.BorderColor = tcell.GetColor(conf.Theme.Border)
	tview.Styles.TitleColor = tcell.GetColor(conf.Theme.Title)
	tview.Styles.GraphicsColor = tcell.GetColor(conf.Theme.Graphics)
	tview.Styles.PrimaryTextColor = tcell.GetColor(conf.Theme.Text)
	tview.Styles.SecondaryTextColor = tcell.GetColor(conf.Theme.Text)
	tview.Styles.TertiaryTextColor = tcell.GetColor(conf.Theme.Text)
	tview.Styles.InverseTextColor = tcell.GetColor(conf.Theme.Text)
	tview.Styles.ContrastSecondaryTextColor = tcell.GetColor(conf.Theme.Text)

	app = tview.NewApplication()
	app.
		EnableMouse(conf.Mouse).
		SetInputCapture(onAppInputCapture)

	channelsTreeView = ui.NewChannelsTreeView()
	channelsTreeView.SetSelectedFunc(onChannelsTreeSelected)

	messagesTextView = ui.NewMessagesTextView()
	messagesTextView.
		SetChangedFunc(func() { app.Draw() }).
		SetInputCapture(onMessagesViewInputCapture)

	messageInputField = ui.NewMessageInputField()
	messageInputField.SetInputCapture(onMessageInputFieldInputCapture)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(messagesTextView, 0, 1, false).
		AddItem(messageInputField, 3, 1, false)
	mainFlex = tview.NewFlex().
		AddItem(channelsTreeView, 0, 1, false).
		AddItem(rightFlex, 0, 4, false)

	token, err := keyring.Get(service, "token")
	if err != nil {
		token = conf.Token
	}

	if token != "" {
		app.
			SetRoot(mainFlex, true).
			SetFocus(channelsTreeView)

		session = newSession()
		session.Token = token
		session.Identify.Token = token
		if err := session.Open(); err != nil {
			panic(err)
		}
	} else {
		loginForm = ui.NewLoginForm(onLoginFormLoginButtonSelected, false)
		app.SetRoot(loginForm, true)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func onLoginFormLoginButtonSelected() {
	email := loginForm.GetFormItem(0).(*tview.InputField).GetText()
	password := loginForm.GetFormItem(1).(*tview.InputField).GetText()
	if email == "" || password == "" {
		return
	}

	session = newSession()
	// Login using the email and password
	lr, err := login(email, password)
	if err != nil {
		panic(err)
	}

	if lr.Token != "" && !lr.MFA {
		app.
			SetRoot(mainFlex, true).
			SetFocus(channelsTreeView)

		session.Token = lr.Token
		session.Identify.Token = lr.Token
		if err = session.Open(); err != nil {
			panic(err)
		}

		go keyring.Set(service, "token", lr.Token)
	} else if lr.MFA {
		// The account has MFA enabled, reattempt login with code and ticket.
		loginForm = ui.NewLoginForm(func() {
			code := loginForm.GetFormItem(0).(*tview.InputField).GetText()
			if code == "" {
				return
			}

			lr, err = totp(code, lr.Ticket)
			if err != nil {
				panic(err)
			}

			app.
				SetRoot(mainFlex, true).
				SetFocus(channelsTreeView)

			session.Token = lr.Token
			session.Identify.Token = lr.Token
			if err = session.Open(); err != nil {
				panic(err)
			}

			go keyring.Set(service, "token", lr.Token)
		}, true)

		app.SetRoot(loginForm, true)
	}
}
