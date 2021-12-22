package ui

import (
	"github.com/ayntgl/discordgo"
	"github.com/ayntgl/discordo/config"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application

	ChannelsTreeView  *tview.TreeView
	MessagesTextView  *tview.TextView
	MessageInputField *tview.InputField

	Session *discordgo.Session
	Config  config.Config
}

func NewApp() *App {
	s, _ := discordgo.New()
	return &App{
		Application: tview.NewApplication(),

		ChannelsTreeView:  tview.NewTreeView(),
		MessagesTextView:  tview.NewTextView(),
		MessageInputField: tview.NewInputField(),

		Session: s,
		Config:  config.LoadConfig(),
	}
}

func (app *App) Connect(token string) error {
	app.Session.UserAgent = "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0"
	app.Session.Identify = discordgo.Identify{
		Token: token,
		Properties: discordgo.IdentifyProperties{
			OS:      "Linux",
			Browser: "Firefox",
			Device:  "",
		},
		Compress:           false,
		Intents:            0,
		LargeThreshold:     0,
		GuildSubscriptions: false,
	}

	return app.Session.Open()
}
