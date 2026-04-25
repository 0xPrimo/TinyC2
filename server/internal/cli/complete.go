package cli

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

func (c *Cli) Completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	args := strings.Split(text, " ")
	currword := d.GetWordBeforeCursor()

	if c.SessionID != 0 {
		return c.completeInteractive(args, currword)
	}

	if currword != "" && len(args) <= 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "plugin", Description: "Manage plugins"},
			{Text: "listener", Description: "Manage listeners"},
			{Text: "implant", Description: "Manage implants"},
			{Text: "help", Description: "Help menu"},
			{Text: "exit", Description: "Exit"},
		}, currword, true)
	}

	cmd := args[0]
	cmdargs := args[1:]

	switch cmd {
	case "plugin":
		return c.completePlugin(cmdargs, currword)
	case "listener":
		return c.completeListener(cmdargs, currword)
	case "implant":
		return c.completeImplant(cmdargs, currword)
	}

	return []prompt.Suggest{}
}

func (c *Cli) completeInteractive(args []string, currword string) []prompt.Suggest {
	if currword != "" && len(args) <= 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "channel", Description: "Manage implant channels"},
			{Text: "ps", Description: "List process"},
			{Text: "cd", Description: "Change directory"},
			{Text: "cp", Description: "Copy file"},
			{Text: "shell", Description: "Run shell command via cmd.exe"},
			{Text: "download", Description: "Download file"},
			{Text: "upload", Description: "Upload file"},
			{Text: "run", Description: "Run executable"},
			{Text: "execute-assembly", Description: "Execute .NET application"},
			{Text: "job", Description: "Manage implant jobs"},
			{Text: "help", Description: "Help menu"},
			{Text: "back", Description: "Exit interactive mode"},
		}, currword, true)
	}

	cmd := args[0]
	cmdargs := args[1:]

	switch cmd {
	case "channel":
		return c.completeInteractiveChannel(cmdargs, currword)
	case "job":
		if len(args) == 1 {
			return prompt.FilterHasPrefix([]prompt.Suggest{
				{Text: "stop"},
				{Text: "list"},
			}, currword, true)
		}
	}

	return []prompt.Suggest{}
}

func (c *Cli) completeInteractiveChannel(args []string, currword string) []prompt.Suggest {
	if len(args) == 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "register"},
			{Text: "switch"},
			{Text: "list"},
			{Text: "remove"},
		}, currword, true)
	}

	if len(args) == 2 {
		subcmd := args[0]
		var suggest []prompt.Suggest

		switch subcmd {
		case "register":
			for name := range c.Engine.Listeners {
				suggest = append(suggest, prompt.Suggest{Text: name})
			}
			return prompt.FilterHasPrefix(suggest, currword, true)

		case "remove", "switch":
			if implant, ok := c.Engine.Implants[c.SessionID]; ok {
				for name := range implant.Channels {
					suggest = append(suggest, prompt.Suggest{Text: name})
				}
			}
			return prompt.FilterHasPrefix(suggest, currword, true)
		}
	}

	return []prompt.Suggest{}
}

func (c *Cli) completePlugin(args []string, currword string) []prompt.Suggest {
	if len(args) == 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "register", Description: "Register plugin"},
			{Text: "remove", Description: "Remove plugin"},
			{Text: "list", Description: "List plugins"},
		}, currword, true)
	}

	return []prompt.Suggest{}
}

func (c *Cli) completeImplant(args []string, currword string) []prompt.Suggest {
	if len(args) == 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "interact", Description: "Interact with implant"},
			{Text: "kill", Description: "Kill implant"},
			{Text: "list", Description: "List implants"},
		}, currword, true)
	}

	if len(args) == 2 {
		subcmd := args[0]
		var suggest []prompt.Suggest

		switch subcmd {
		case "kill", "interact":
			for id := range c.Engine.Implants {
				suggest = append(suggest, prompt.Suggest{Text: fmt.Sprintf("%X", id)})
			}
			return prompt.FilterHasPrefix(suggest, currword, true)
		}
	}
	return []prompt.Suggest{}
}

func (c *Cli) completeListener(args []string, currword string) []prompt.Suggest {
	if len(args) == 1 {
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "start", Description: "Start listener"},
			{Text: "stop", Description: "Stop listener"},
			{Text: "generate", Description: "Generate implant EXE"},
			{Text: "list", Description: "List listeners"},
		}, currword, true)
	}

	if len(args) == 2 {
		subcmd := args[0]
		var suggest []prompt.Suggest

		switch subcmd {
		case "start":
			for name := range c.Engine.Plugins {
				suggest = append(suggest, prompt.Suggest{Text: name})
			}
			return prompt.FilterHasPrefix(suggest, currword, true)

		case "stop", "generate":
			for name := range c.Engine.Listeners {
				suggest = append(suggest, prompt.Suggest{Text: name})
			}
			return prompt.FilterHasPrefix(suggest, currword, true)
		}
	}

	return []prompt.Suggest{}
}
