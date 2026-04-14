package handlers

import (
	"strconv"

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

func HandleImplantWhoami(engine *core.Engine, session *uint32, args []string) {
	engine.ImplantWhoami(*session)
}

func HandleImplantPs(engine *core.Engine, session *uint32, args []string) {
	engine.ImplantPs(*session)
}
