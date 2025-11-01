package ui

import "gowizcli/wiz"

func resetData() fetchDoneMsg {
	return fetchDoneMsg{
		lights: []wiz.Light{},
		err:    nil,
	}
}

type fetchDoneMsg struct {
	lights []wiz.Light
	err    error
}

type switchDoneMsg struct {
	light wiz.Light
	err   error
}

type discoverDoneMsg struct {
	lights []wiz.Light
	err    error
}

func resetDiscoverData() discoverDoneMsg {
	return discoverDoneMsg{
		lights: []wiz.Light{},
		err:    nil,
	}
}
