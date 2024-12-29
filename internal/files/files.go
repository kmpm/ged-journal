package files

type Module struct {
	logPath string
}

func New(logPath string) *Module {
	return &Module{
		logPath: logPath,
	}
}
