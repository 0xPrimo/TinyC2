package cmd

import (
	"fmt"
	"os"

	"github.com/0xPrimo/TinyC2/server/internal/engine"
	lua "github.com/yuin/gopher-lua"

	"github.com/c-bata/go-prompt"
)

type UserCommand struct {
	Name        string
	Description string
	Callback    *lua.LFunction
}

type Cli struct {
	id               string
	L                *lua.LState
	ScriptedCommands map[string][]UserCommand
	Engine           *engine.Engine
}

func NewCli(path string) *Cli {
	engine := engine.NewEngine(path)

	return &Cli{
		L:                lua.NewState(),
		ScriptedCommands: map[string][]UserCommand{},
		Engine:           engine,
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

	c.luaInitialize()

	// Debug
	c.Executor("listener_start tcp tcp-1 ../plugins/tcp/config.yaml")
	c.Executor("implant_generate tcp-1 /home/primo/Share/implant.exe")

	defer os.Exit(0)
	defer c.L.Close()
	p.Run()
}

func (c *Cli) LivePrefix() (string, bool) {
	if c.id == "" {
		return "", false
	}
	return fmt.Sprintf("implant(%s)> ", c.id), true
}
