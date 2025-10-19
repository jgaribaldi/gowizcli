package ui

import (
	"fmt"
	"gowizcli/client"
	"gowizcli/wiz"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DiscoverModel struct {
	client           *client.Client
	discovering      bool
	broadcastAddress string
	inputs           []textinput.Model
	focused          int
	lights           []wiz.WizLight
	err              error
}

func NewDiscoverModel(client *client.Client) DiscoverModel {
	var inputs []textinput.Model
	inputs = make([]textinput.Model, 0)

	input1 := textinput.New()
	input1.Placeholder = "192"
	input1.CharLimit = 3
	input1.Width = 3
	input1.Prompt = ""
	input1.Validate = octetValidator
	input1.Focus()
	inputs = append(inputs, input1)

	input2 := textinput.New()
	input2.Placeholder = "168"
	input2.CharLimit = 3
	input2.Width = 3
	input2.Prompt = ""
	input2.Validate = octetValidator
	inputs = append(inputs, input2)

	input3 := textinput.New()
	input3.Placeholder = "1"
	input3.CharLimit = 3
	input3.Width = 3
	input3.Prompt = ""
	input3.Validate = octetValidator
	inputs = append(inputs, input3)

	input4 := textinput.New()
	input4.Placeholder = "255"
	input4.CharLimit = 3
	input4.Width = 3
	input4.Prompt = ""
	input4.Validate = octetValidator
	inputs = append(inputs, input4)

	var lights []wiz.WizLight = make([]wiz.WizLight, 0)
	return DiscoverModel{
		client:           client,
		discovering:      false,
		broadcastAddress: "",
		inputs:           inputs,
		focused:          0,
		lights:           lights,
		err:              nil,
	}
}

func (m DiscoverModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DiscoverModel) Update(msg tea.Msg) (DiscoverModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m = nextInput(m)
		case tea.KeyShiftTab:
			m = previousInput(m)
		case tea.KeyEnter:
			m.broadcastAddress = fmt.Sprintf(
				"%s.%s.%s.%s",
				m.inputs[0].View(),
				m.inputs[1].View(),
				m.inputs[2].View(),
				m.inputs[3].View(),
			)
			discoverLightsCmd(m.client, m.broadcastAddress)
			m.discovering = true
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	case discoverOkMsg:
		m.discovering = false
		for _, l := range msg.lights {
			m.lights = append(m.lights, l)
		}
	case discoverErrorMsg:
		m.discovering = false
		m.err = msg.err
	}

	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m DiscoverModel) View() string {
	if m.broadcastAddress == "" {
		return fmt.Sprintf(
			"Broadcast address: %s.%s.%s.%s",
			m.inputs[0].View(),
			m.inputs[1].View(),
			m.inputs[2].View(),
			m.inputs[3].View(),
		)
	}
	if m.discovering {
		return fmt.Sprintf("Executing discovery on broadcast %s...", m.broadcastAddress)
	}

	if m.err != nil {
		return fmt.Sprintf("Error executing discover: %s", m.err)
	}
	return fmt.Sprintf("Finished discover and found %d lights - Esc to go back to main menu", len(m.lights))
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

func previousInput(m DiscoverModel) DiscoverModel {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
	return m
}

func nextInput(m DiscoverModel) DiscoverModel {
	m.focused = (m.focused + 1) % len(m.inputs)
	return m
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
