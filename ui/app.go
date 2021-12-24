package ui

import (
	"sort"

	"github.com/ayntgl/discordo/config"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application

	GuildsList        *tview.List
	ChannelsTreeView  *tview.TreeView
	MessagesTextView  *tview.TextView
	MessageInputField *tview.InputField

	State  *state.State
	Config config.Config
}

func NewApp() *App {
	return &App{
		Application: tview.NewApplication(),

		GuildsList:        tview.NewList(),
		ChannelsTreeView:  tview.NewTreeView(),
		MessagesTextView:  tview.NewTextView(),
		MessageInputField: tview.NewInputField(),

		Config: config.LoadConfig(),
	}
}

// Connect initializes a new state and modifies default properties, adds gateway event handlers, and opens a new websocket connection to the Discord gateway.
func (app *App) Connect(token string) error {
	app.State = state.New(token)

	api.UserAgent = app.Config.General.UserAgent
	gateway.DefaultIdentity = gateway.IdentifyProperties{
		OS:      "Linux",
		Browser: "Firefox",
		Device:  "",
	}
	app.State.AddHandler(app.onSessionReady)

	return app.State.Open(app.State.Context())
}

func (app *App) onSessionReady(r *gateway.ReadyEvent) {
	sort.Slice(r.Guilds, func(a, b int) bool {
		found := false
		for _, guildID := range r.UserSettings.GuildPositions {
			if found && guildID == r.Guilds[b].ID {
				return true
			}
			if !found && guildID == r.Guilds[a].ID {
				found = true
			}
		}

		return false
	})

	for _, g := range r.Guilds {
		app.GuildsList.AddItem(g.Name, "", 0, nil)
	}
}
