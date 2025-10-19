package ui

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/wiz"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ipAddressInput struct {
	inputs  []textinput.Model
	focused int
}

func newIpAddressInput() ipAddressInput {
	var inputs []textinput.Model = make([]textinput.Model, 4)

	inputs[0] = newOctetInput("192")
	inputs[1] = newOctetInput("168")
	inputs[2] = newOctetInput("1")
	inputs[3] = newOctetInput("255")
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

func (i ipAddressInput) GetValue() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		strings.TrimSpace(i.inputs[0].View()),
		strings.TrimSpace(i.inputs[1].View()),
		strings.TrimSpace(i.inputs[2].View()),
		strings.TrimSpace(i.inputs[3].View()),
	)
}

func (i ipAddressInput) Update(msg tea.Msg) (ipAddressInput, []tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(i.inputs))

	i.inputs[0], cmds[0] = i.inputs[0].Update(msg)
	i.inputs[1], cmds[1] = i.inputs[1].Update(msg)
	i.inputs[2], cmds[2] = i.inputs[2].Update(msg)
	i.inputs[3], cmds[3] = i.inputs[3].Update(msg)

	return i, cmds
}

type discoverData struct {
	lights []wiz.WizLight
	err    error
}

func newDiscoverData() discoverData {
	return discoverData{
		lights: make([]wiz.WizLight, 0),
		err:    nil,
	}
}

func (d discoverData) Result(lights []wiz.WizLight) discoverData {
	for _, l := range lights {
		d.lights = append(d.lights, l)
	}
	d.err = nil
	return d
}

func (d discoverData) Error(err error) discoverData {
	d.err = err
	d.lights = make([]wiz.WizLight, 0)
	return d
}

type DiscoverModel struct {
	client    *client.Client
	input     ipAddressInput
	data      discoverData
	cmdStatus commandStatus
}

func NewDiscoverModel(client *client.Client) DiscoverModel {
	return DiscoverModel{
		client:    client,
		input:     newIpAddressInput(),
		data:      newDiscoverData(),
		cmdStatus: newCommandStatus(),
	}
}

func newOctetInput(ph string) textinput.Model {
	input := textinput.New()
	input.Placeholder = ph
	input.CharLimit = 3
	input.Width = 3
	input.Prompt = ""
	input.Validate = octetValidator

	return input
}

func (m DiscoverModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DiscoverModel) Update(msg tea.Msg) (DiscoverModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.input = m.input.nextOctet()
		case tea.KeyShiftTab:
			m.input = m.input.previousOctet()
		case tea.KeyEnter:
			if !m.cmdStatus.isStarted() {
				broadcastAddress := m.input.GetValue()
				m.cmdStatus = m.cmdStatus.start()
				return m, discoverLightsCmd(m.client, broadcastAddress)
			}
			if m.cmdStatus.isFinished() {
				m.cmdStatus = m.cmdStatus.reset()
				return m, nil
			}
		}

	case discoverOkMsg:
		m.data = m.data.Result(msg.lights)
		m.cmdStatus = m.cmdStatus.finish()
	case discoverErrorMsg:
		m.data = m.data.Error(msg.err)
		m.cmdStatus = m.cmdStatus.finish()
	}

	var cmds []tea.Cmd = make([]tea.Cmd, 0)
	m.input, cmds = m.input.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m DiscoverModel) View() string {
	if m.cmdStatus.isRunning() {
		return "Executing discovery..."
	}

	if m.cmdStatus.isFinished() {
		if m.data.err != nil {
			return fmt.Sprintf("Error executing discover: %s", m.data.err)
		} else {
			return fmt.Sprintf("Finished discover and found %d lights - Esc to go back to main menu", len(m.data.lights))
		}
	}

	return m.input.GetValue()
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

type commandStatus struct {
	running  bool
	finished bool
}

func newCommandStatus() commandStatus {
	return commandStatus{
		running:  false,
		finished: false,
	}
}

func (s commandStatus) start() commandStatus {
	s.running = true
	s.finished = false
	return s
}

func (s commandStatus) finish() commandStatus {
	s.running = false
	s.finished = true
	return s
}

func (s commandStatus) reset() commandStatus {
	s.running = false
	s.finished = false
	return s
}

func (s commandStatus) isFinished() bool {
	return !s.running && s.finished
}

func (s commandStatus) isStarted() bool {
	return s.running && !s.finished
}

func (s commandStatus) isRunning() bool {
	return s.running
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
	lights []wiz.WizLight
}

type discoverErrorMsg struct {
	err error
}
