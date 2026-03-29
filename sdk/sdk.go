package sdk

type IEngine interface {
	ImplantProcess(listener string, data []byte) ([]byte, error)

	LogInfo(plugin string, format string, args ...any)
	LogError(plugin string, format string, args ...any)
	LogWarning(plugin string, format string, args ...any)
	LogSuccess(plugin string, format string, args ...any)
}

type IListener interface {
	Start() error
	Stop() error
	MakePic(uint32) ([]byte, error)
}
