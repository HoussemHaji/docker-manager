package ui

import (
	"docker-cli-tool/internal/docker"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateTable displays Docker containers in a table and shows an action menu when a container is selected
func CreateTable(app *tview.Application) *tview.Table {
	table := tview.NewTable().
		SetBorders(false)

	// Set table headers
	headers := []string{"ID", "STATUS", "CONTAINER NAME"}
	for col, header := range headers {
		table.SetCell(0, col,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignLeft).
				SetExpansion(1))
	}

	// Retrieve Docker container information
	containers, err := docker.ListContainers(true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Populate the table with container data
	for i, container := range containers {
		// Format container ID and name
		containerID := container.ID[:12] // Show first 12 characters of ID
		containerName := container.Names[0][1:]

		// Set the status color based on container state
		statusText := strings.ToUpper(container.State)
		statusColor := tcell.ColorGreen
		if container.State != "running" {
			statusColor = tcell.ColorRed
		}

		// Add rows to the table
		table.SetCell(i*2+1, 0, tview.NewTableCell(containerID).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1))

		table.SetCell(i*2+1, 1, tview.NewTableCell(statusText).
			SetTextColor(statusColor).
			SetExpansion(1))

		table.SetCell(i*2+1, 2, tview.NewTableCell(containerName).
			SetTextColor(tcell.ColorWhite).
			SetExpansion(1))

		// Add separator row
		if i < len(containers)-1 {
			separatorRow := i*2 + 2
			for col := 0; col < 3; col++ {
				table.SetCell(separatorRow, col,
					tview.NewTableCell("-----------------------").
						SetTextColor(tcell.ColorGray).
						SetExpansion(1))
			}
		}
	}

	// Enable row selection
	table.SetSelectable(true, false)

	// Set selected style
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkBlue))

	// Set a function to trigger on row selection
	table.SetSelectedFunc(func(row, column int) {
		if row%2 == 0 || row == 0 {
			return // Skip header and separator rows
		}

		// Get selected container details
		containerIndex := row / 2
		containerID := containers[containerIndex].ID
		containerName := containers[containerIndex].Names[0][1:]

		// Show the container action menu
		showActionMenu(app, containerID, containerName)
	})

	return table
}

// showActionMenu displays a menu of actions for a selected container in the TextView without a modal
func showActionMenu(app *tview.Application, containerID, containerName string) {
	// Create a text view to display the menu
	menu := tview.NewTextView().
		SetDynamicColors(true)

	userInput := ""

	updateMenuText := func() {
		menuText := fmt.Sprintf(`You have selected: [yellow]%s [green]%s[white]

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

Please select your action: %s`, containerName, containerID, userInput)

		menu.SetText(menuText)
	}

	updateMenuText()

	// Handle user input for the menu
	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			userInput += string(event.Rune())
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(userInput) > 0 {
				userInput = userInput[:len(userInput)-1]
			}
		case tcell.KeyEnter:
			if len(userInput) > 0 {
				action := strings.ToLower(userInput[:1])
				userInput = ""
				switch action {
				case "s": // Start container
					go performContainerAction(app, menu, "Starting container..", func() error {
						return docker.StartContainer(containerID)
					})
				case "x": // Stop container
					go performContainerAction(app, menu, "Stopping container..", func() error {
						return docker.StopContainer(containerID)
					})
				case "r": // Restart container
					go performContainerAction(app, menu, "Restarting container..", func() error {
						return docker.RestartContainer(containerID)
					})
				case "p": // Pause container
					go performContainerAction(app, menu, "Pausing container..", func() error {
						return docker.PauseContainer(containerID)
					})
				case "u": // Unpause container
					go performContainerAction(app, menu, "Unpausing container..", func() error {
						return docker.UnpauseContainer(containerID)
					})
				case "l": // Show container logs
					go performContainerAction(app, menu, "Fetching container logs..", func() error {
						return docker.GetContainerLogs(containerID)
					})
				case "d": // Delete container
					go performContainerAction(app, menu, "Deleting container..", func() error {
						return docker.DeleteContainer(containerID)
					})
				case "b": // Go back to the table
					app.SetRoot(CreateTable(app), true)
				default:
					userInput = action
				}
			}
		}

		updateMenuText()
		return nil
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
