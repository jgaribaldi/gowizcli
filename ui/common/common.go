package common

type state int

const (
	Ready state = iota
	Running
	Done
)

type CmdStatus struct {
	State state
}

func NewCmdStatus() *CmdStatus {
	return &CmdStatus{
		State: Ready,
	}
}

func (c CmdStatus) Start() CmdStatus {
	c.State = Running
	return c
}

func (c CmdStatus) Finish() CmdStatus {
	c.State = Done
	return c
}
