package ui

import (
	"docker-cli-tool/internal/docker"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// containers table

func CreateTable(app *tview.Application) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	containers, err := docker.ListContainers(true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for i, container := range containers {
		table.SetCell(i+1, 0, tview.NewTableCell(container.ID).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 1, tview.NewTableCell(container.Names[0]).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 2, tview.NewTableCell(container.State).SetAlign(tview.AlignCenter))
	}

	table.SetSelectable(true, false).SetSelectedFunc(func(row, column int) {
		containerID := table.GetCell(row, 0).Text
		fmt.Printf("Selected container ID: %s\n", containerID)
	})

	return table
}
