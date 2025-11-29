package ui

import (
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"

	tea "github.com/charmbracelet/bubbletea"
)

type CmdRunner struct {
	client        client.Functions
	lastCmdStatus common.CmdStatus
	lastCmdErr    error
}

func NewCmdRunner(client client.Functions) CmdRunner {
	return CmdRunner{
		client:        client,
		lastCmdStatus: *common.NewCmdStatus(),
		lastCmdErr:    nil,
	}
}

func (c CmdRunner) Run(cmd Command) (CmdRunner, tea.Cmd) {
	if c.lastCmdStatus.State == common.Ready || c.lastCmdStatus.State == common.Done {
		c.lastCmdStatus = c.lastCmdStatus.Start()
		return c, func() tea.Msg {
			result, err := cmd.Run()
			return CmdDone{
				lights: result,
				err:    err,
				cmd:    cmd,
			}
		}
	}
	return c, nil
}

func (c CmdRunner) Finalize(msg CmdDone) CmdRunner {
	c.lastCmdStatus = c.lastCmdStatus.Finish()
	c.lastCmdErr = msg.err
	return c
}

type CmdDone struct {
	lights []wiz.Light
	err    error
	cmd    Command
}

type Command interface {
	Run() ([]wiz.Light, error)
}

type CmdDiscover struct {
	client client.Functions
}

func NewCmdDiscover(client client.Functions) CmdDiscover {
	return CmdDiscover{
		client: client,
	}
}

func (c CmdDiscover) Run() ([]wiz.Light, error) {
	return c.client.Discover()
}

type CmdSwitch struct {
	client client.Functions
	light  wiz.Light
}

func NewCmdSwitch(client client.Functions, light wiz.Light) CmdSwitch {
	return CmdSwitch{
		client: client,
		light:  light,
	}
}

func (c CmdSwitch) Run() ([]wiz.Light, error) {
	if c.light.IsOn != nil && *c.light.IsOn {
		result, err := c.client.TurnOff(c.light.Id)
		if err != nil {
			return nil, err
		}
		return []wiz.Light{*result}, nil
	}

	result, err := c.client.TurnOn(c.light.Id)
	if err != nil {
		return nil, err
	}
	return []wiz.Light{*result}, nil
}

type CmdEraseAll struct {
	client client.Functions
}

func NewCmdEraseAll(client client.Functions) CmdEraseAll {
	return CmdEraseAll{
		client: client,
	}
}

func (c CmdEraseAll) Run() ([]wiz.Light, error) {
	c.client.EraseAll()
	return nil, nil
}

type CmdRefresh struct {
	client client.Functions
}

func NewCmdRefresh(client client.Functions) CmdRefresh {
	return CmdRefresh{
		client: client,
	}
}

func (c CmdRefresh) Run() ([]wiz.Light, error) {
	return c.client.ShowAll()
}
