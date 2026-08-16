package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
)

func (e *Engine) ImplantGenerate(listenerName string, dest string) error {

	extension, err := e.ListenerExtension(listenerName)
	if err != nil {
		return err
	}

	src, _ := filepath.Abs("../implant")
	binary, err := buildCmakeProject(src, extension)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, binary, 0o644)
	if err != nil {
		return err
	}

	return nil
}

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
