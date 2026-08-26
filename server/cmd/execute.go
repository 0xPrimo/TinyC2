package cmd

import (
	"sort"
	"strconv"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/pterm/pterm"
	lua "github.com/yuin/gopher-lua"
)

type Command struct {
	Name              string
	Description       string
	Example           string
	NumberOfArguments int
	Execute           func(core *Cli, args ...string) error
}

func (c *Cli) Executor(in string) {
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	args := strings.Split(in, " ")
	var cmds [][]string

	// help menu
	if args[0] == "help" {
		// server
		for _, cmd := range ServerCommandList {
			cmds = append(cmds, []string{cmd.Name, cmd.Description})
		}

		// implant
		if c.id != "" {
			for _, cmd := range c.Engine.ImplantCommandList() {
				cmds = append(cmds, []string{cmd.Name, cmd.Description})
			}

			for _, cmd := range ImplantCommandList {
				cmds = append(cmds, []string{cmd.Name, cmd.Description})
			}

			// script commands
			for _, script := range c.ScriptedCommands {
				for _, cmd := range script {
					cmds = append(cmds, []string{cmd.Name, cmd.Description})
				}
			}
		}

		// sort commands
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i][0] < cmds[j][0]
		})

		printHelpMenu(cmds)
	}

	if c.id != "" {
		// implant commands
		for _, cmd := range c.Engine.ImplantCommandList() {
			if cmd.Name == args[0] {
				err := c.Engine.ImplantExecute(c.id, args[0], args...)
				if err != nil {
					logger.Error("%v", err)
					return
				}
				return
			}
		}

		// server side implant commands
		for _, command := range ImplantCommandList {
			if args[0] == command.Name {
				if len(args)-1 < command.NumberOfArguments {
					logger.Error("usage: \n\t%s\n", command.Example)
					return
				}

				err := command.Execute(c, args...)
				if err != nil {
					logger.Error("%v", err)
					return
				}
				return
			}
		}

		// script registered commands
		for _, script := range c.ScriptedCommands {
			for _, cmd := range script {
				if cmd.Name == args[0] {
					c.L.Push(cmd.Callback)

					num, err := strconv.ParseUint(c.id, 16, 32)
					if err != nil {
						logger.Error("Error: %v", err)
						return
					}

					c.L.Push(lua.LNumber(num))
					for _, arg := range args[1:] {
						c.L.Push(lua.LString(arg))
					}

					if err := c.L.PCall(len(args[1:])+1, 0, nil); err != nil {
						logger.Error("Failed to run script: %v", err)
					}
					return
				}
			}
		}
	}

	// server commands
	for _, command := range ServerCommandList {
		if args[0] == command.Name {
			if len(args)-1 < command.NumberOfArguments {
				logger.Error("usage: \n\t%s\n", command.Example)
				return
			}

			err := command.Execute(c, args...)
			if err != nil {
				logger.Error("%v", err)
			}
			return
		}
	}

}

func printHelpMenu(commands [][]string) {
	var commandData pterm.TableData

	for _, cmd := range commands {
		if len(cmd) >= 2 {
			coloredName := pterm.LightCyan("  " + cmd[0])
			commandData = append(commandData, []string{coloredName, cmd[1]})
		}
	}
	pterm.DefaultTable.WithData(commandData).Render()
}

//
//	command := args[0]
//	cmdargs := args[1:]
//
//	if c.SessionID != 0 {
//		switch command {
//		case "channel":
//			if len(args) < 1 {
//				logger.Info("channel [register|switch|list|remove]")
//				return
//			}
//
//			err := c.Engine.ImplantExecute(c.id, "channel", args...)
//			if err != nil {
//				logger.Error("failed to execute command: %v", err)
//				return
//			}
//		case "ps":
//			handler.HandleImplantPs(c.Engine, &c.SessionID, cmdargs)
//		case "cd":
//			handler.HandleImplantCd(c.Engine, &c.SessionID, cmdargs)
//		case "cp":
//			handler.HandleImplantCp(c.Engine, &c.SessionID, cmdargs)
//		case "shell":
//			handler.HandleImplantShell(c.Engine, &c.SessionID, cmdargs)
//		case "download":
//			handler.HandleImplantDownload(c.Engine, &c.SessionID, cmdargs)
//		case "upload":
//			handler.HandleImplantUpload(c.Engine, &c.SessionID, cmdargs)
//		case "run":
//			handler.HandleImplantRun(c.Engine, &c.SessionID, cmdargs)
//		case "execute-assembly":
//			handler.HandleImplantExecuteAssembly(c.Engine, &c.SessionID, cmdargs)
//		case "job":
//			handler.HandleImplantJob(c.Engine, &c.SessionID, cmdargs)
//		case "inline-execute":
//			handler.HandleImplantInlineExecute(c.Engine, &c.SessionID, cmdargs)
//		case "inject-shellcode":
//			handler.HandleImplantInjectShellcode(c.Engine, &c.SessionID, cmdargs)
//		case "token_info":
//			handler.HandleImplantTokenInfo(c.Engine, &c.SessionID, cmdargs)
//		case "token_rev2self":
//			handler.HandleImplantTokenRev2Self(c.Engine, &c.SessionID, cmdargs)
//		case "token_make":
//			handler.HandleImplantTokenMake(c.Engine, &c.SessionID, cmdargs)
//		case "token_steal":
//			handler.HandleImplantTokenSteal(c.Engine, &c.SessionID, cmdargs)
//		case "connect":
//			handler.HandleImplantPivotConnect(c.Engine, &c.SessionID, cmdargs)
//		case "back":
//			c.SessionID = 0
//		case "help":
//			logger.Info(
//				`Usage:
//	    channel                                       - Manage channels
//	       register [name]                            - Register channel
//	       switch   [name]                            - Switch channel
//	       remove   [name]                            - Remove channel
//	       list                                       - List channels
//	                -----------------------------
//	    job
//	       stop   [id]                                - Stop job
//	       list                                       - List jobs
//	                -----------------------------
//	    ps                                            - List process
//	    cd                                            - Change process working directory
//	    cp                                            - Copy file to target directory
//	    shell                                         - Task a shell command via cmd.exe
//	    download                                      - Download file from target machine
//	    upload                                        - Upload file to target machine
//	    run                                           - Task executable that exits on target machine
//	    execute-assembly                              - Task a .NET application
//	    inline-execute                                - Task a Beacon Object File
//        shellcode-inject                              - Inject shellcode into running process
//        token_info                                    - Show current process token informations
//		token_rev2self                                - Release any access token that have been created or stolen
//		token_make                                    - Impersonate a user
//		token_steal                                   - Steal a process access token
//        back                                          - Exit interactive mode
//	    help                                          - Print help menu
//	`,
//			)
//		case "exit":
//			handler.HandleExit(c.Engine, cmdargs)
//		default:
//
//			for _, cmds := range c.UserCommands {
//				for _, cmd := range cmds {
//					if cmd.Name == command {
//						c.L.Push(cmd.Callback)
//						c.L.Push(lua.LNumber(c.SessionID))
//						for _, arg := range cmdargs {
//							c.L.Push(lua.LString(arg))
//						}
//
//						if err := c.L.PCall(len(cmdargs)+1, 0, nil); err != nil {
//							fmt.Printf("Failed to run script: %v\n", err)
//						}
//						return
//					}
//				}
//			}
//
//			fmt.Printf("Unknown command: %s\n", command)
//			return
//		}
//		return
//	}
//
//	switch command {
//	case "plugin":
//		handler.HandlePlugin(c.Engine, cmdargs)
//	case "listener":
//		handler.HandleListener(c.Engine, cmdargs)
//	case "implant":
//		handler.HandleImplant(c.Engine, &c.SessionID, cmdargs)
//	case "script_load":
//		if len(cmdargs) < 1 {
//			logger.Error("need arguments")
//			return
//		}
//		c.HandleScriptLoad(cmdargs[0])
//	case "script_unload":
//		if len(cmdargs) < 1 {
//			logger.Error("need arguments")
//			return
//		}
//		c.HandleScriptUnload(cmdargs[0])
//
//	case "help":
//		logger.Info(
//			`Usage:
//      implant                                       - Manage implants
//         interact [id]                              - Interact with implant
//         kill     [id]                              - Kill implant
//         list                                       - List implants
//
//      plugin                                        - Manage plugins
//         register [name] [plugin.so]                - Register listener plugin
//         remove   [name]                            - Remove listener plugin
//         list                                       - List plugins
//
//      listener                                      - Manage listeners
//         generate [name] [dest]						- Generate implant (EXE)
//	     start    [plugin] [name] [config.yaml]     - Start listener
//         stop     [name]                            - Stop listener
//         list                                       - List listeners
//
//	  script_load   [path]                          - Load lua script
//	  script_unload [path]                          - Unload lua script
//      help                                          - Print help menu
//`,
//		)
//	case "exit":
//		handler.HandleExit(c.Engine, cmdargs)
//	default:
//		for _, cmds := range c.UserCommands {
//			for _, cmd := range cmds {
//				if cmd.Name == command {
//					c.L.Push(cmd.Callback)
//					c.L.Push(lua.LNumber(c.SessionID))
//					for _, arg := range cmdargs {
//						c.L.Push(lua.LString(arg))
//					}
//
//					if err := c.L.PCall(len(cmdargs)+1, 0, nil); err != nil {
//						fmt.Printf("Failed to run script: %v\n", err)
//					}
//					return
//				}
//			}
//		}
//
//		fmt.Printf("Unknown command: %s\n", command)
//		return
//	}
