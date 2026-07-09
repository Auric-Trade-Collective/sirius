package monitor

type ProcessCommand interface {
	Action() func(ProcessEvent)
}

type ProcessEvent interface {}

type ProcessExit struct {
	Code int
}

type ProcessCrash struct {
	Code int
}
