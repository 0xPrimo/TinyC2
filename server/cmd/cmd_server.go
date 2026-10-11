package cmd

import (
	"fmt"
	"os"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/utils"
	"github.com/pterm/pterm"
)

var ServerCommandList = []Command{
	{
		Name:              "plugin_register",
		Description:       "Register plugin",
		Example:           "plugin_register [/path/to/go/plugin]",
		NumberOfArguments: 1,
		Execute: func(core *Cli, args ...string) error {
			meta, err := core.Engine.PluginRegister(core.Engine, args[1])
			if err != nil {
				return err
			}
			logger.Success("plugin registered: %s [%s]", meta.Name, meta.Type)
			return nil
		},
	},
	{
		Name:              "plugin_unregister",
		Description:       "Unregister plugin",
		NumberOfArguments: 1,
		Example:           "plugin_unregister tcp",
		Execute: func(core *Cli, args ...string) error {
			return core.Engine.PluginUnregister(args[1])
		},
	},
	{
		Name:              "plugin_list",
		Description:       "List plugins",
		Example:           "",
		NumberOfArguments: 0,
		Execute: func(core *Cli, args ...string) error {
			var rows [][]string
			for _, plugin := range core.Engine.PluginList() {
				rows = append(rows, []string{plugin.Name, plugin.Type, plugin.Path})
			}

			utils.PrintTable([]string{"Name", "Type", "Path"}, rows)
			return nil
		},
	},
	{
		Name:              "listener_start",
		Description:       "Start listener",
		NumberOfArguments: 3,
		Example:           "listener_start [plugin-name] [listener-name] [path/to/config.yaml]",
		Execute: func(core *Cli, args ...string) error {
			err := core.Engine.ListenerStart(args[1], args[2], args[3])
			if err != nil {
				return err
			}

			logger.Success("listener %s started", args[2])
			return nil
		},
	},
	{
		Name:              "listener_stop",
		Description:       "Stop listener",
		NumberOfArguments: 1,
		Example:           "listener_stop [listener-name]",
		Execute: func(core *Cli, args ...string) error {
			err := core.Engine.ListenerStop(args[1])
			if err != nil {
				return err
			}
			logger.Success("listener stopped: %s", args[1])
			return nil
		},
	},
	{
		Name:              "listener_list",
		Description:       "List listeners",
		NumberOfArguments: 0,
		Example:           "",
		Execute: func(core *Cli, args ...string) error {
			var rows [][]string
			for _, listener := range core.Engine.ListenerList() {
				rows = append(rows, []string{listener.Name, listener.Protocol})
			}
			utils.PrintTable([]string{"Name", "Protocol"}, rows)
			return nil
		},
	},
	{
		Name:              "implant_generate",
		Description:       "Generate implant payload",
		NumberOfArguments: 2,
		Example:           "implant_generate [listener-name] [/out/payload.exe]",
		Execute: func(core *Cli, args ...string) error {
			payload, err := core.Engine.ImplantGenerate(args[1])
			if err != nil {
				return err
			}

			err = os.WriteFile(args[2], payload, 0600)
			if err != nil {
				return err
			}

			logger.Success("implant generated: %s", args[2])
			return nil
		},
	},
	{
		Name:              "implant_list",
		Description:       "List implants",
		Example:           "",
		NumberOfArguments: 0,
		Execute: func(core *Cli, args ...string) error {
			var rows [][]string
			for _, implant := range core.Engine.ImplantList() {
				var status string

				if implant.Alive() {
					status = pterm.Green("alive")
				} else {
					status = pterm.Red("dead")
				}

				rows = append(rows, []string{
					implant.ID,
					implant.ChannelCurrent(),
					implant.Meta["user"].(string),
					implant.Meta["host"].(string),
					implant.Meta["domain"].(string),
					implant.Meta["pid"].(string),
					implant.Meta["os"].(string),
					status,
				})
			}

			utils.PrintTable([]string{"ID", "Channel", "User", "Computer", "Domain", "Pid", "Os", "Status"}, rows)
			return nil
		},
	},
	{
		Name:              "implant_kill",
		Description:       "Kill implant",
		Example:           "implant_kill [id]",
		NumberOfArguments: 1,
		Execute: func(core *Cli, args ...string) error {
			return core.Engine.ImplantExecute(args[1], "kill")
		},
	},
	{
		Name:              "implant_interact",
		Description:       "Interact with implant",
		Example:           "implant_interact [id]",
		NumberOfArguments: 1,
		Execute: func(core *Cli, args ...string) error {
			for _, implant := range core.Engine.ImplantList() {
				if args[1] == implant.ID {
					core.id = args[1]
					return nil
				}
			}
			return fmt.Errorf("unknown implant id: %s", args[1])

		},
	},
	{
		Name:              "script_load",
		Description:       "Load lua script",
		Example:           "script_load [/path/to/script.lua]",
		NumberOfArguments: 1,
		Execute: func(core *Cli, args ...string) error {
			if err := core.L.DoFile(string(args[1])); err != nil {
				return err
			}

			logger.Success("Script loaded")
			return nil
		},
	},
	{
		Name:              "script_unload",
		Description:       "Unload lua script",
		Example:           "script_unload [script-name]",
		NumberOfArguments: 1,
		Execute: func(core *Cli, args ...string) error {
			delete(core.ScriptedCommands, args[1])
			logger.Success("Script unloaded")
			return nil
		},
	},
}
