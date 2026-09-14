package implant

type IImplantManager interface {
	ImplantProcess(listener string, data []byte) ([]byte, error)
	ImplantExecute(id string, cmd string, args ...string) error
	ImplantGenerate(listenerName string) ([]byte, error)
	ImplantList() []Implant
	ImplantCommandList() []Command
	ImplantChannelList(id string) ([]Channel, bool)
	ImplantExists(id string) bool
	ImplantDBSync() error
}

type Channel struct {
	ID       uint32
	Name     string
	Fallback bool
	InUse    bool
}

type Packet struct {
	ID         string
	TaskResult []TaskResult
}

type Command struct {
	Name              string
	Description       string
	NumberOfArguments int
	execute           func(manager *Manager, id string, args ...string) (Task, error)
	process           func(manager *Manager, id string, result TaskResult)
}

type Task struct {
	Cmd       string
	Args      []any
	Artifacts []Artifact
}

type TaskResult struct {
	Cmd      string
	Status   string
	Output   string
	Artifact []byte
}

type Artifact struct {
	Name string
	Path string
	Data []byte
}
type rawTask struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Output   string `json:"output"`
	Artifact string `json:"artifact"`
}

type rawPacket struct {
	ID    *uint64   `json:"id"`
	Tasks []rawTask `json:"tasks"`
}
