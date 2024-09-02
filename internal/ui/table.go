package ui

import (
	"docker-cli-tool/internal/docker"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateTable displays Docker containers in a formatted table without background color
func CreateTable(app *tview.Application) *tview.Table {
	// Create a new table with borders enabled
	table := tview.NewTable().
		SetBorders(true)

	// Set table headers with no background color
	table.SetCell(0, 0, tview.NewTableCell("ID").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDefault)) // Set background color to default

	table.SetCell(0, 1, tview.NewTableCell("Name").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDefault))

	table.SetCell(0, 2, tview.NewTableCell("Status").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDefault))

	// Retrieve Docker container information
	containers, err := docker.ListContainers(true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Populate the table with container data
	for i, container := range containers {
		// Format container ID (trim it to the first 10 characters for readability)
		containerID := container.ID[:10]

		// Format container name (remove the leading '/' from the container name)
		containerName := container.Names[0][1:]

		// Set the status and apply color based on the container's state
		statusText := container.State
		statusColor := tcell.ColorGreen
		if container.State != "running" {
			statusText = "STOPPED"
			statusColor = tcell.ColorRed
		}

		// Add rows to the table, with transparent grey background color
		table.SetCell(i+1, 0, tview.NewTableCell(containerID).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorDefault))

		table.SetCell(i+1, 1, tview.NewTableCell(containerName).
			SetTextColor(tcell.ColorBlue).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorDefault))

		table.SetCell(i+1, 2, tview.NewTableCell(statusText).
			SetTextColor(statusColor).
			SetAlign(tview.AlignCenter).
			SetBackgroundColor(tcell.ColorDefault))
	}

	// Enable row selection but not column selection
	table.SetSelectable(true, false)

	//set the width of the table to full
	table.SetFixed(1, 3)

	// Set the selection color to grey with transparency
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Dim(true))

	// Set a function to trigger on row selection
	table.SetSelectedFunc(func(row, column int) {
		containerID := table.GetCell(row, 0).Text
		fmt.Printf("Selected container ID: %s\n", containerID)
		// You can add further actions here (e.g., stop/start a container by ID)
	})

	return table
}
