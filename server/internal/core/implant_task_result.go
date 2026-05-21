package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/pterm/pterm"
)

type TaskResultHandler func(TaskResult) error

var TaskResultHandlers = map[string]TaskResultHandler{
	"channel.switch":   handleOutput,
	"channel.regsiter": handleOutput,
	"cd":               handleOutput,
	"cp":               handleOutput,
	"shell":            handleOutput,
	"upload":           handleOutput,
	"run":              handleOutput,
	"execute-assembly": handleOutput,
	"inline-execute":   handleOutput,
	"inject-shellcode": handleOutput,
	"ps":               handlePs,
	"download":         handleDownload,
	"job.list":         handleJobList,
}

func (e *Engine) ImplantTaskResultDispatch(id uint32, tasks []TaskResult) {
	// process tasks
	for _, task := range tasks {
		if handler, ok := TaskResultHandlers[task.Name]; ok {
			if err := handler(task); err != nil {
				logger.Error("failed to process task result %s result: %v", task.Name, err)
			}
		}
	}
}

type FileInfo struct {
	Name string `json:"name"`
	Size uint32 `json:"size"`
	Data string `json:"data"`
}

func handleDownload(task TaskResult) error {
	var file FileInfo

	err := json.Unmarshal([]byte(task.Artifact), &file)
	if err != nil {
		logger.Error("error occurred during unmarshaling: %v", err)
		return fmt.Errorf("error occurred during unmarshaling: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(file.Data)
	if err != nil {
		logger.Error("failed to decode base64: %v", err)
		return fmt.Errorf("failed to decode base64: %v", err)
	}

	err = os.MkdirAll("uploads", 0o755)
	if err != nil {
		logger.Error("failed to create directory: %v", err)
		return fmt.Errorf("failed to create directory: %v", err)
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

	err = os.WriteFile(dest, data, 0o644)
	if err != nil {
		logger.Error("failed to write to file: %v", err)
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

type ProcessInfo struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	Pid     uint32 `json:"pid"`
	PPid    uint32 `json:"ppid"`
}

func handlePs(task TaskResult) error {
	var pslist []ProcessInfo
	err := json.Unmarshal([]byte(task.Artifact), &pslist)
	if err != nil {
		logger.Error("error occurred during unmarshaling: %v", err)
		return fmt.Errorf("error occurred during unmarshaling: %v", err)
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

type JobInfo struct {
	ID uint32 `json:"id"`
}

func handleJobList(task TaskResult) error {
	var joblist []JobInfo
	err := json.Unmarshal([]byte(task.Artifact), &joblist)
	if err != nil {
		logger.Error("error occurred during unmarshaling: %v", err)
		return fmt.Errorf("error occurred during unmarshaling: %v", err)
	}

	table := pterm.TableData{
		{"ID"},
	}

	for _, job := range joblist {
		table = append(table, []string{
			pterm.Cyan(fmt.Sprintf("%X", job.ID)),
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

func handleOutput(task TaskResult) error {
	logger.Success("recieved output:\n%s\n", task.Output)
	return nil
}
