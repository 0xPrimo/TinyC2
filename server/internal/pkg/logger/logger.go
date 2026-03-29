package logger

import (
	"fmt"

	"github.com/pterm/pterm"
)

func Info(format string, a ...any) {
	fmt.Printf("%s\n", fmt.Sprintf("[%s] %s", pterm.Bold.Sprint(pterm.Blue("*")), fmt.Sprintf(format, a...)))
}

func Success(format string, a ...any) {
	fmt.Printf("%s\n", fmt.Sprintf("[%s] %s", pterm.Bold.Sprint(pterm.Green("+")), fmt.Sprintf(format, a...)))
}

func Warning(format string, a ...any) {
	fmt.Printf("%s\n", fmt.Sprintf("[%s] %s", pterm.Bold.Sprint(pterm.Yellow("!")), fmt.Sprintf(format, a...)))
}

func Error(format string, a ...any) {
	fmt.Printf("%s\n", fmt.Sprintf("[%s] %s", pterm.Bold.Sprint(pterm.Red("-")), fmt.Sprintf(format, a...)))
}
