package main

import (
	"os"

	"github.com/ayntgl/discordo/ui"
	"github.com/zalando/go-keyring"
)

const serviceName = "discordo"

func main() {
	app := ui.NewApp()
	app.EnableMouse(app.Config.General.Mouse)

	token := os.Getenv("DISCORDO_TOKEN")
	if token == "" {
		token, _ = keyring.Get(serviceName, "token")
	}

	if token != "" {
		err := app.Connect(token)
		if err != nil {
			panic(err)
		}

		app.
			SetRoot(ui.NewMainFlex(app), true).
			SetFocus(app.GuildsList)
	} else {
		loginForm := ui.NewLoginForm().
			AddButton("Login", nil)
		app.SetRoot(loginForm, true)
	}

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
