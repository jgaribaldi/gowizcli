package common

type CommandStatus struct {
	running  bool
	finished bool
}

func NewCommandStatus() CommandStatus {
	return CommandStatus{
		running:  false,
		finished: false,
	}
}

func (s CommandStatus) Start() CommandStatus {
	s.running = true
	s.finished = false
	return s
}

func (s CommandStatus) Finish() CommandStatus {
	s.running = false
	s.finished = true
	return s
}

func (s CommandStatus) Reset() CommandStatus {
	s.running = false
	s.finished = false
	return s
}

func (s CommandStatus) IsFinished() bool {
	return !s.running && s.finished
}

func (s CommandStatus) IsStarted() bool {
	return s.running && !s.finished
}

func (s CommandStatus) IsRunning() bool {
	return s.running
}
