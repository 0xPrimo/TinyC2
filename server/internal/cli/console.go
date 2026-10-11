package cli

import (
	"fmt"
	"os"

	"github.com/0xPrimo/TinyC2/server/internal/core"
	lua "github.com/yuin/gopher-lua"

	"github.com/c-bata/go-prompt"
)

type UserCommand struct {
	Name        string
	Description string
	Callback    *lua.LFunction
}

type Cli struct {
	SessionID    uint32
	L            *lua.LState
	UserCommands map[string][]UserCommand
	Commands     []prompt.Suggest
	Engine       *core.Engine
}

func NewCli(path string) *Cli {
	return &Cli{
		L:            lua.NewState(),
		UserCommands: make(map[string][]UserCommand),
		Engine:       core.NewEngine(path),
		Commands: []prompt.Suggest{
			{Text: "exit", Description: "Exit the console"},
		},
	}
}

func (c *Cli) Start() {
	p := prompt.New(
		c.Executor,
		c.Completer,
		prompt.OptionPrefix("tinyc2> "),
		prompt.OptionLivePrefix(c.LivePrefix),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
	)

	c.L.SetGlobal("register_command", c.L.NewFunction(c.luaRegisterCommand))

	apiTable := c.L.NewTable()
	c.L.SetField(apiTable, "inline_execute", c.L.NewFunction(c.luaInlineExecute))
	c.L.SetField(apiTable, "execute_assembly", c.L.NewFunction(c.luaExecuteAssembly))
	c.L.SetField(apiTable, "bof_pack", c.L.NewFunction(c.luaBofPack))
	c.L.SetGlobal("api", apiTable)
	defer os.Exit(0)
	defer c.L.Close()
	p.Run()
}

func (c *Cli) LivePrefix() (string, bool) {
	if c.SessionID == 0 {
		return "", false
	}
	return fmt.Sprintf("implant(%X)> ", c.SessionID), true
}
