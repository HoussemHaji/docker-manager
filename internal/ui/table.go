package ui

import (
	"docker-cli-tool/internal/docker"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateTable displays Docker containers in a table and shows an action menu when a container is selected
func CreateTable(app *tview.Application) *tview.Table {
	// Create a new table with borders enabled
	table := tview.NewTable().
		SetBorders(true)

	// Set table headers
	table.SetCell(0, 0, tview.NewTableCell("ID").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDefault))

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
		// Format container ID and name
		containerID := container.ID[:10]
		containerName := container.Names[0][1:]

		// Set the status color based on container state
		statusText := container.State
		statusColor := tcell.ColorGreen
		if container.State != "running" {
			statusText = "STOPPED"
			statusColor = tcell.ColorRed
		}

		// Add rows to the table
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

	// focus bg color for selected row to gray
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGray))

	// Set a function to trigger on row selection
	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return // Skip header row
		}

		// Get selected container details
		containerID := containers[row-1].ID
		containerName := containers[row-1].Names[0][1:]

		// Show the container action menu
		showActionMenu(app, containerID, containerName)
	})

	return table
}

// showActionMenu displays a menu of actions for a selected container in the TextView without a modal
func showActionMenu(app *tview.Application, containerID, containerName string) {
	// Create a text view to display the menu
	menu := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(`You have selected: [yellow]%s [green]%s[white]

Options:
  [ s ] Start
  [ x ] Stop
  [ r ] Restart
  [ p ] Pause
  [ u ] Unpause
  ---------------
  [ l ] Show Logs
  [ d ] Delete
  ---------------
  [ b ] Go back

Please select your action:`, containerName, containerID))

	// Handle user input for the menu
	// the menu should respond by typing the corresponding letter and showing it on the terminal and then pressing enter
	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 's': // Start container
			go performContainerAction(app, menu, "Starting container..", func() error {
				return docker.StartContainer(containerID)
			})
		case 'x': // Stop container
			go performContainerAction(app, menu, "Stopping container..", func() error {
				return docker.StopContainer(containerID)
			})
		case 'r': // Restart container
			go performContainerAction(app, menu, "Restarting container..", func() error {
				return docker.RestartContainer(containerID)
			})
		case 'p': // Pause container
			go performContainerAction(app, menu, "Pausing container..", func() error {
				return docker.PauseContainer(containerID)
			})
		case 'u': // Unpause container
			go performContainerAction(app, menu, "Unpausing container..", func() error {
				return docker.UnpauseContainer(containerID)
			})
		case 'l': // Show container logs
			go performContainerAction(app, menu, "Fetching container logs..", func() error {
				return docker.GetContainerLogs(containerID)
			})
		case 'b': // Go back to the table
			app.SetRoot(CreateTable(app), true)
		}
		return nil // No need for 'Enter' to confirm
	})

	// Show the menu
	app.SetRoot(menu, true).SetFocus(menu)
}

// performContainerAction performs the container action and updates the UI
func performContainerAction(app *tview.Application, menu *tview.TextView, message string, actionFunc func() error) {
	menu.SetText(message)
	err := actionFunc()
	time.Sleep(2 * time.Second) // Simulate some delay for user experience
	if err != nil {
		menu.SetText(fmt.Sprintf("[red]Failed[white] to perform action: %v", err))
	} else {
		menu.SetText("[green]Action completed successfully![white]\n\nPress [ b ] to go back to the table.")
	}
	app.Draw() // Force a redraw to update the UI
}
