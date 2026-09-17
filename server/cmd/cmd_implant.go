package cmd

import (
	"fmt"

	"github.com/0xPrimo/TinyC2/server/internal/utils"
	"github.com/pterm/pterm"
)

var ImplantCommandList = []Command{
	{
		Name:              "back",
		Description:       "run session in background",
		NumberOfArguments: 0,
		Execute: func(core *Cli, args ...string) error {
			core.id = ""
			return nil
		},
	},
	{
		Name:              "channel.list",
		Description:       "list implant communication channels",
		NumberOfArguments: 0,
		Execute: func(core *Cli, args ...string) error {
			channels, ok := core.Engine.ImplantChannelList(core.id)
			if !ok {
				return fmt.Errorf("implant not found")
			}

			var (
				rows [][]string
			)

			for _, channel := range channels {
				var row []string

				row = append(row, channel.Name)

				if channel.InUse {
					row = append(row, pterm.Green("true"))
				} else {
					row = append(row, pterm.Red("false"))
				}

				//if channel.Fallback {
				//	row = append(row, pterm.Green("true"))
				//} else {
				//	row = append(row, pterm.Red("false"))
				//}

				rows = append(rows, row)
			}

			utils.PrintTable([]string{"Name", "InUse"}, rows)
			return nil
		},
	},
}
