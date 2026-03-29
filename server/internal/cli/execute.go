package cli

import (
	"strings"
	"github.com/0xPrimo/TinyC2/server/internal/cli/handlers"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
)

func (c *Cli) Executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}

	args := strings.Split(in, " ")
	command := args[0]
	cmdargs := args[1:]

	if c.SessionID != 0 {
		switch command {
		case "channel":
			handlers.HandleImplantChannel(c.Engine, &c.SessionID, cmdargs)
		case "whoami":
			handlers.HandleImplantWhoami(c.Engine, &c.SessionID, cmdargs)
		case "back":
			c.SessionID = 0
		case "help":
			logger.Info(
				`Usage:
      channel                                       - Manage channels
         register [name]                            - Register channel
         switch   [name]                            - Switch channel
         list                                       - List channels

      whoami                                        - Get implant context
      back                                          - Exit interactive mode
      help                                          - Print help menu
`)
		case "exit":
			handlers.HandleExit(c.Engine, cmdargs)
		default:
			logger.Error("unknown command: %s\n", command)
		}
		return
	}

	switch command {
	case "plugin":
		handlers.HandlePlugin(c.Engine, cmdargs)
	case "listener":
		handlers.HandleListener(c.Engine, cmdargs)
	case "implant":
		handlers.HandleImplant(c.Engine, &c.SessionID, cmdargs)
	case "help":
		logger.Info(
			`Usage:
      implant                                       - Manage implants
         interact [id]                              - Interact with implant
         kill     [id]                              - Kill implant
         list                                       - List implants

      plugin                                        - Manage plugins
         register [name] [plugin.so]                - Register listener plugin
         remove   [name]                            - Remove listener plugin
         list                                       - List plugins

      listener                                      - Manage listeners
         generate [name] [dest]						- Generate implant (EXE)
	     start    [plugin] [name] [config.yaml]     - Start listener
         stop     [name]                            - Stop listener
         list                                       - List listeners

      help                                          - Print help menu
`)
	case "exit":
		handlers.HandleExit(c.Engine, cmdargs)
	default:
		logger.Error("unknown command: %s\n", command)
	}
}
