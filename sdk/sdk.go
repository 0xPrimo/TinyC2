package sdk

type IEngine interface {
	// ImplantProcess process implant response ( expect json ) and returns
	// queued tasks
	ImplantProcess(listener string, data []byte) ([]byte, error)

	// Logger APIs
	LogInfo(plugin string, format string, args ...any)
	LogError(plugin string, format string, args ...any)
	LogWarning(plugin string, format string, args ...any)
	LogSuccess(plugin string, format string, args ...any)
}

type IListener interface {
	// Start start the listener
	Start() error

	// Stop stop the listener
	Stop() error

	// MakePic return PIC
	MakePic(uint32) ([]byte, error)
}
