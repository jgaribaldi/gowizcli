package discover

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/ui/common"
	"gowizcli/wiz"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	client    *client.Client
	input     ipAddressInput
	data      discoverData
	cmdStatus common.CmdStatus
}

func NewModel(client *client.Client, defaultBcastAddr string) Model {
	return Model{
		client:    client,
		input:     newIpAddressInput(defaultBcastAddr),
		data:      newDiscoverData(),
		cmdStatus: *common.NewCmdStatus(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.input = m.input.nextOctet()
		case tea.KeyShiftTab:
			m.input = m.input.previousOctet()
		case tea.KeyEnter:
			if m.cmdStatus.State == common.Ready {
				broadcastAddress := m.input.getValue()
				m.cmdStatus = m.cmdStatus.Start()
				return m, discoverLightsCmd(m.client, broadcastAddress)
			}
			if m.cmdStatus.State == common.Done {
				m.cmdStatus = *common.NewCmdStatus()
				m.data = newDiscoverData()
				return m, nil
			}
		}

	case discoverOkMsg:
		m.data = m.data.result(msg.lights)
		m.cmdStatus = m.cmdStatus.Finish()
	case discoverErrorMsg:
		m.data = m.data.error(msg.err)
		m.cmdStatus = m.cmdStatus.Finish()
	}

	var cmds []tea.Cmd = make([]tea.Cmd, 0)
	m.input, cmds = m.input.update(msg)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.cmdStatus.State == common.Running {
		return "Executing discovery..."
	}

	if m.cmdStatus.State == common.Done {
		if m.data.err != nil {
			return fmt.Sprintf("Error executing discover: %s", m.data.err)
		} else {
			return fmt.Sprintf("Finished discover and found %d lights - Esc to go back to main menu", len(m.data.lights))
		}
	}

	return m.input.getValue()
}

type ipAddressInput struct {
	inputs  []textinput.Model
	focused int
}

func newIpAddressInput(defaultBcastAddr string) ipAddressInput {
	octets := strings.Split(defaultBcastAddr, ".")

	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[0] = newOctetInput(octets[0], "192")
	inputs[1] = newOctetInput(octets[1], "168")
	inputs[2] = newOctetInput(octets[2], "1")
	inputs[3] = newOctetInput(octets[3], "255")
	inputs[0].Focus()

	return ipAddressInput{
		inputs:  inputs,
		focused: 0,
	}
}

func (i ipAddressInput) previousOctet() ipAddressInput {
	i.focused--
	if i.focused < 0 {
		i.focused = len(i.inputs) - 1
	}
	for ii := range i.inputs {
		i.inputs[ii].Blur()
	}
	i.inputs[i.focused].Focus()
	return i
}

func (i ipAddressInput) nextOctet() ipAddressInput {
	i.focused = (i.focused + 1) % len(i.inputs)
	for ii := range i.inputs {
		i.inputs[ii].Blur()
	}
	i.inputs[i.focused].Focus()
	return i
}

func (i ipAddressInput) getValue() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		strings.TrimSpace(i.inputs[0].View()),
		strings.TrimSpace(i.inputs[1].View()),
		strings.TrimSpace(i.inputs[2].View()),
		strings.TrimSpace(i.inputs[3].View()),
	)
}

func (i ipAddressInput) update(msg tea.Msg) (ipAddressInput, []tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(i.inputs))

	i.inputs[0], cmds[0] = i.inputs[0].Update(msg)
	i.inputs[1], cmds[1] = i.inputs[1].Update(msg)
	i.inputs[2], cmds[2] = i.inputs[2].Update(msg)
	i.inputs[3], cmds[3] = i.inputs[3].Update(msg)

	return i, cmds
}

func newOctetInput(defaultValue string, ph string) textinput.Model {
	input := textinput.New()
	input.SetValue(defaultValue)
	input.Placeholder = ph
	input.CharLimit = 3
	input.Width = 3
	input.Prompt = ""
	input.Validate = octetValidator

	return input
}

func octetValidator(octet string) error {
	number, err := strconv.ParseInt(octet, 10, 64)
	if err != nil {
		return err
	}

	if number < 0 || number > 255 {
		return fmt.Errorf("incorrect octet: %s", octet)
	}

	return nil
}

type discoverData struct {
	lights []wiz.Light
	err    error
}

func newDiscoverData() discoverData {
	return discoverData{
		lights: make([]wiz.Light, 0),
		err:    nil,
	}
}

func (d discoverData) result(lights []wiz.Light) discoverData {
	for _, l := range lights {
		d.lights = append(d.lights, l)
	}
	d.err = nil
	return d
}

func (d discoverData) error(err error) discoverData {
	d.err = err
	d.lights = make([]wiz.Light, 0)
	return d
}

func discoverLightsCmd(c *client.Client, broadcastAddress string) tea.Cmd {
	return func() tea.Msg {
		cmd := client.Command{
			CommandType: client.Discover,
			Parameters: []string{
				broadcastAddress,
			},
		}
		result, err := c.Execute(cmd)
		if err != nil {
			return discoverErrorMsg{err: err}
		}
		return discoverOkMsg{lights: result}
	}
}

type discoverOkMsg struct {
	lights []wiz.Light
}

type discoverErrorMsg struct {
	err error
}
