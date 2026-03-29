package handlers

import (
	"os"
	"os/exec"
	"github.com/0xPrimo/TinyC2/server/internal/core"
)

func HandleExit(engine *core.Engine, args []string) {
	cmd := exec.Command("stty", "-raw", "echo")
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
	os.Exit(0)
}
