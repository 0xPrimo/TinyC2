package core

import (
	"fmt"
	"hash/crc32"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/0xPrimo/TinyC2/sdk"
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

type Listener struct {
	ID        uint32
	Name      string
	Interface sdk.IListener
	Config    string
}

func (e *Engine) ListenerStart(plugin string, name string, config string) error {
	_, exists := e.Listeners[name]
	if exists {
		return fmt.Errorf("listener %s already exist", name)
	}

	listener, err := e.PluginNewListener(plugin, name, config)
	if err != nil {
		return err
	}

	err = listener.Start()
	if err != nil {
		return err
	}

	e.Listeners[name] = Listener{
		ID:        crc32.ChecksumIEEE([]byte(name)),
		Name:      name,
		Config:    config,
		Interface: listener,
	}

	return nil
}

func (e *Engine) ListenerStop(name string) error {
	listener, exists := e.Listeners[name]
	if !exists {
		return fmt.Errorf("listener %s does not exist", name)
	}

	err := listener.Interface.Stop()
	if err != nil {
		return err
	}

	delete(e.Listeners, name)

	return nil
}

func (e *Engine) ListenerList() error {
	table := pterm.TableData{
		{"ID", "Name", "Config"},
	}

	for _, listener := range e.Listeners {
		table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", listener.ID)), listener.Name, listener.Config})
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

func (e *Engine) ListenerGenerate(name string, dest string) error {
	listener, exists := e.Listeners[name]
	if !exists {
		return fmt.Errorf("listener not found")
	}

	// generate pic with id 0. id 0 means default implant channel
	//
	pic, err := listener.Interface.MakePic(0)
	if err != nil {
		return err
	}

	// err = os.WriteFile("../debug.pic", pic, 0644)
	// if err != nil {
	// 	return err
	// }

	// build cmake project
	//
	src, _ := filepath.Abs("./implant")
	binary, err := buildCmakeProject(src, pic)
	if err != nil {
		return err
	}

	// write binary
	//
	err = os.WriteFile(dest, binary, 0644)
	if err != nil {
		return err
	}

	return nil
}

func buildCmakeProject(src string, config []byte) ([]byte, error) {

	os.MkdirAll(src+"/build", 0755)

	// build project
	//
	args := []string{
		"-S", src,
		"-B", src + "/build",
		"-DDEFAULT_CHANNEL=" + toCArray(config),
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

	// read binary
	//
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
