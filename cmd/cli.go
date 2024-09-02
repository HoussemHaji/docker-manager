package cmd

import (
	"docker-cli-tool/internal/ui"

	"github.com/rivo/tview"
)

// initialize cli and start app
func StartCli() {
	app := tview.NewApplication()
	table := ui.CreateTable(app)

	flex := tview.NewFlex().AddItem(table, 0, 1, true)

	if err := app.SetRoot(flex, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
