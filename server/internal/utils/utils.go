package utils

import "github.com/pterm/pterm"

func PrintTable(headers []string, rows [][]string) {
	table := pterm.TableData{
		headers,
	}

	table = append(table, rows...)

	pterm.Println()

	pterm.DefaultTable.
		WithHasHeader(true).
		WithBoxed(true).
		WithHeaderStyle(pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold)).
		WithData(table).
		Render()

	pterm.Println()

}
