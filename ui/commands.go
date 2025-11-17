package ui

import (
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"

	tea "github.com/charmbracelet/bubbletea"
)

type CmdRunner struct {
	client        *client.Client
	lastCmdStatus common.CmdStatus
	lastCmdErr    error
}

func NewCmdRunner(client *client.Client) CmdRunner {
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
	client *client.Client
}

func NewCmdDiscover(client *client.Client) CmdDiscover {
	return CmdDiscover{
		client: client,
	}
}

func (c CmdDiscover) Run() ([]wiz.Light, error) {
	cmd := client.Command{
		CommandType: client.Discover,
		Parameters:  []string{},
	}
	return c.client.Execute(cmd)
}

type CmdSwitch struct {
	client *client.Client
	light  wiz.Light
}

func NewCmdSwitch(client *client.Client, light wiz.Light) CmdSwitch {
	return CmdSwitch{
		client: client,
		light:  light,
	}
}

func (c CmdSwitch) Run() ([]wiz.Light, error) {
	if c.light.IsOn != nil && *c.light.IsOn {
		cmd := client.Command{
			CommandType: client.TurnOff,
			Parameters: []string{
				c.light.Id,
			},
		}
		return c.client.Execute(cmd)
	}

	cmd := client.Command{
		CommandType: client.TurnOn,
		Parameters: []string{
			c.light.Id,
		},
	}
	return c.client.Execute(cmd)
}

type CmdEraseAll struct {
	client *client.Client
}

func NewCmdEraseAll(client *client.Client) CmdEraseAll {
	return CmdEraseAll{
		client: client,
	}
}

func (c CmdEraseAll) Run() ([]wiz.Light, error) {
	cmd := client.Command{
		CommandType: client.Reset,
		Parameters:  []string{},
	}

	return c.client.Execute(cmd)
}

type CmdRefresh struct {
	client *client.Client
}

func NewCmdRefresh(client *client.Client) CmdRefresh {
	return CmdRefresh{
		client: client,
	}
}

func (c CmdRefresh) Run() ([]wiz.Light, error) {
	cmd := client.Command{
		CommandType: client.Show,
		Parameters:  []string{},
	}
	return c.client.Execute(cmd)
}
