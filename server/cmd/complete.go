package cmd

import (
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
)

func (c *Cli) Completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	args := strings.Split(text, " ")
	word := d.GetWordBeforeCursor()

	if c.id != "" {
		return c.completeImplantCommand(args, word)
	}

	return c.completeServerCommand(args, word)
}

func (c *Cli) completeImplantCommand(args []string, word string) []prompt.Suggest {
	var suggestions []prompt.Suggest

	if word != "" && len(args) <= 1 {
		for _, cmd := range c.Engine.ImplantCommandList() {
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}

		for _, cmd := range ImplantCommandList {
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}

		suggestions = append(suggestions, c.completeServerCommand(args, word)...)

		sort.Slice(suggestions, func(i, j int) bool {
			return suggestions[i].Text < suggestions[j].Text
		})

		return prompt.FilterHasPrefix(suggestions, word, false)
	}

	cmd := args[0]
	switch cmd {
	case "channel.register":
		for _, listener := range c.Engine.ListenerList() {
			suggestions = append(suggestions, prompt.Suggest{Text: listener.Name})
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Text < suggestions[j].Text
	})

	return prompt.FilterHasPrefix(suggestions, word, false)
}

func (c *Cli) completeServerCommand(args []string, word string) []prompt.Suggest {
	var suggestions []prompt.Suggest

	if word != "" && len(args) <= 1 {
		for _, cmd := range ServerCommandList {
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
		sort.Slice(suggestions, func(i, j int) bool {
			return suggestions[i].Text < suggestions[j].Text
		})
		return prompt.FilterHasPrefix(suggestions, word, false)
	} else if len(args) == 2 {
		switch args[0] {
		case "listener_start":
			for _, plugin := range c.Engine.PluginList() {
				suggestions = append(suggestions, prompt.Suggest{Text: plugin.Name})
			}
		case "listener_stop":
			for _, listener := range c.Engine.ListenerList() {
				suggestions = append(suggestions, prompt.Suggest{Text: listener.Name})
			}
		case "implant_interact":
			for _, implant := range c.Engine.ImplantList() {
				suggestions = append(suggestions, prompt.Suggest{Text: implant.ID})
			}
		case "implant_generate":
			for _, plugin := range c.Engine.PluginList() {
				suggestions = append(suggestions, prompt.Suggest{Text: plugin.Name})
			}
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Text < suggestions[j].Text
	})
	return prompt.FilterHasPrefix(suggestions, word, false)
}
