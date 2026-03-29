package cli

import (
	"fmt"
	"os"

	"github.com/0xPrimo/TinyC2/server/internal/core"

	"github.com/c-bata/go-prompt"
)

type Cli struct {
	SessionID uint32
	Engine    *core.Engine
}

func NewCli(path string) *Cli {
	return &Cli{
		Engine: core.NewEngine(path),
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

	defer os.Exit(0)
	p.Run()
}

func (c *Cli) LivePrefix() (string, bool) {
	if c.SessionID == 0 {
		return "", false
	}
	return fmt.Sprintf("implant(%X)> ", c.SessionID), true
}
