package ui

import (
	"docker-cli-tool/internal/docker"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
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
  [ l ] Show Logs
  [ d ] Delete
  [ e ] Execute Command
  [ n ] View Network Info
  [ f ] Filter Containers
  ---------------
  [ b ] Go back

Please select your action: %s`, containerName, containerID, userInput)

		menu.SetText(menuText)
	}

	updateMenuText()

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
				case "s":
					go performContainerAction(app, menu, "Starting container..", func() error {
						return docker.StartContainer(containerID)
					})
				case "x":
					go performContainerAction(app, menu, "Stopping container..", func() error {
						return docker.StopContainer(containerID)
					})
				case "r":
					go performContainerAction(app, menu, "Restarting container..", func() error {
						return docker.RestartContainer(containerID)
					})
				case "p":
					go performContainerAction(app, menu, "Pausing container..", func() error {
						return docker.PauseContainer(containerID)
					})
				case "u":
					go performContainerAction(app, menu, "Unpausing container..", func() error {
						return docker.UnpauseContainer(containerID)
					})
				case "l":
					go performContainerAction(app, menu, "Fetching container logs..", func() error {
						return docker.GetContainerLogs(containerID)
					})
				case "d":
					go performContainerAction(app, menu, "Deleting container..", func() error {
						return docker.DeleteContainer(containerID)
					})
				case "e":
					showExecuteCommandPrompt(app, containerID)
				case "n":
					showNetworkInfo(app, containerID)
				case "f":
					showFilterOptions(app)
				case "b":
					app.SetRoot(CreateTable(app), true)
				default:
					userInput = action
				}
			}
		}

		updateMenuText()
		return nil
	})

	app.SetRoot(menu, true).SetFocus(menu)
}

func showExecuteCommandPrompt(app *tview.Application, containerID string) {
	input := tview.NewInputField().
		SetLabel("Enter command: ").
		SetFieldWidth(0)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := input.GetText()
			go performContainerAction(app, tview.NewTextView(), fmt.Sprintf("Executing command: %s", command), func() error {
				return docker.ExecCommandInContainer(containerID, strings.Fields(command))
			})
		}
	})

	app.SetRoot(input, true).SetFocus(input)
}

func showNetworkInfo(app *tview.Application, containerID string) {
	go performContainerAction(app, tview.NewTextView(), "Fetching network info..", func() error {
		networkInfo, err := docker.GetContainerNetworkInfo(containerID)
		if err != nil {
			return err
		}
		infoText := "Network Information:\n"
		for network, ip := range networkInfo {
			infoText += fmt.Sprintf("%s: %s\n", network, ip)
		}
		textView := tview.NewTextView().SetText(infoText)
		app.SetRoot(textView, true)
		return nil
	})
}

func showFilterOptions(app *tview.Application) {
	form := tview.NewForm().
		AddDropDown("Filter by:", []string{"Name", "Status"}, 0, nil).
		AddInputField("Filter value:", "", 0, nil, nil).
		AddButton("Apply", nil).
		AddButton("Cancel", func() {
			app.SetRoot(CreateTable(app), true)
		})

	form.GetButton(0).SetSelectedFunc(func() {
		filterType, _ := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
		filterValue := form.GetFormItem(1).(*tview.InputField).GetText()

		var filteredContainers []types.Container
		var err error

		if filterType == 0 { // Filter by Name
			filteredContainers, err = docker.FilterContainersByName(filterValue)
		} else { // Filter by Status
			filteredContainers, err = docker.FilterContainersByStatus(filterValue)
		}

		if err != nil {
			showError(app, err)
			return
		}

		filteredTable := createFilteredTable(app, filteredContainers)
		app.SetRoot(filteredTable, true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func createFilteredTable(app *tview.Application, containers []types.Container) *tview.Table {
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

func showError(app *tview.Application, err error) {
	errorView := tview.NewModal().
		SetText(fmt.Sprintf("Error: %v", err)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(CreateTable(app), true)
		})

	app.SetRoot(errorView, true)
}

func performContainerAction(app *tview.Application, view tview.Primitive, message string, actionFunc func() error) {
	textView, ok := view.(*tview.TextView)
	if !ok {
		textView = tview.NewTextView()
	}
	textView.SetText(message)
	app.SetRoot(textView, true)

	err := actionFunc()
	if err != nil {
		textView.SetText(fmt.Sprintf("[red]Failed[white] to perform action: %v", err))
	} else {
		textView.SetText("[green]Action completed successfully![white]\n\nPress any key to go back to the table.")
	}

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		app.SetRoot(CreateTable(app), true)
		return nil
	})

	app.Draw()
}
