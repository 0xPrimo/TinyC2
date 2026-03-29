package core

import (
	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"

	"github.com/pterm/pterm"
)

func (e *Engine) LogInfo(plugin string, format string, args ...any) {
	logger.Info(pterm.Sprint(pterm.Cyan(plugin))+": "+format, args...)
}

func (e *Engine) LogError(plugin string, format string, args ...any) {
	logger.Error(pterm.Sprint(pterm.Cyan(plugin))+": "+format, args...)
}

func (e *Engine) LogWarning(plugin string, format string, args ...any) {
	logger.Warning(pterm.Sprint(pterm.Cyan(plugin))+": "+format, args...)
}

func (e *Engine) LogSuccess(plugin string, format string, args ...any) {
	logger.Success(pterm.Sprint(pterm.Cyan(plugin))+": "+format, args...)
}
