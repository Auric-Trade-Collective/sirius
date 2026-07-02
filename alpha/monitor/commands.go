package monitor

type ProcessCommand interface {
	Action() func(ProcessEvent)
}
