package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"

	tea "github.com/charmbracelet/bubbletea"
)

// Models & model methods

type model struct {
	keymap    keymap
	help      help.Model
	stopwatch stopwatch.Model

	game_state int

	last_jump         int64
	gate_lower_height int
	gate_upper_height int
	player_position   int
}

type keymap struct {
	action key.Binding
	quit   key.Binding
}

const (
	starting = iota
	playing
	// paused
	dead
	exiting
)

func get_starting_state() model {
	u, l := get_first_gate()
	m := model{
		keymap: keymap{
			action: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "<<ERROR>>"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "<<ERROR>>"),
			),
		},
		help:              help.NewModel(),
		stopwatch:         stopwatch.NewWithInterval(time.Millisecond),
		game_state:        starting,
		last_jump:         0,
		gate_upper_height: u,
		gate_lower_height: l,
		player_position:   get_starting_player_position(),
	}

	return m
}

// View methods

func (m model) View() string {
	switch m.game_state {
	case starting:
		return get_starting_view(m)
	case playing:
		return get_game_view(m)
	case dead:
		return get_dead_view(m)
	case exiting:
		return get_aborted_view(m)
	default:
		return "<< ERROR >>"
	}
}

func get_starting_view(m model) string {
	return "<< READY >>"
}

func get_game_view(m model) string {
	m.keymap.action.SetHelp("space", "jump")
	m.keymap.quit.SetHelp("q", "abort")

	return fmt.Sprintf(
		"CUR ALTITUDE:\t\t%dm\n"+
			"NEXT GATE CEILING:\t%dm\n"+
			"NEXT GATE FLOOR:\t%dm\n"+
			"NEXT GATE IN:\t\t%dm\n"+
			"MISSION CLOCK:\t\t%dms\n"+
			"%s",
		m.player_position,
		m.gate_upper_height,
		m.gate_lower_height,
		get_distance_to_next_gate(m.stopwatch.Elapsed().Milliseconds()),
		m.stopwatch.Elapsed().Milliseconds(),
		m.help.ShortHelpView([]key.Binding{
			m.keymap.action,
			m.keymap.quit,
		}),
	)
}

func get_dead_view(m model) string {
	return "<< DEAD >>"
}

func get_aborted_view(m model) string {
	return "<< MISSION ABORTED >>"
}

// Controller methods

func (m model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.game_state {
	case playing:
		return get_playing_update(m, msg)
	default:
		return get_halted_update(m, msg)
	}
}

func get_playing_update(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	// Tick the game state, pre-input
	m.player_position = get_next_player_position(m.player_position, m.last_jump)
	m.gate_upper_height, m.gate_lower_height = get_next_gate(m.gate_upper_height, m.gate_lower_height, m.stopwatch.Elapsed().Milliseconds())

	if !is_player_alive(m.gate_upper_height, m.gate_lower_height, m.player_position, m.stopwatch.Elapsed().Milliseconds()) {
		m.game_state = dead
	} else {
		m.last_jump += 1

		// React as necessary to key input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keymap.action):
				m.last_jump = 0
			case key.Matches(msg, m.keymap.quit):
				m.game_state = exiting
			}
		}
	}

	// Update the stopwatch
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func get_halted_update(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	// React as necessary to key input
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.action):
			m.game_state = playing
			return m, m.stopwatch.Reset()
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		}
	}

	// Update the stopwatch
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

// ~

func main() {
	if err := tea.NewProgram(get_starting_state()).Start(); err != nil {
		fmt.Printf("Houston, we have a problem (%v)", err)
		os.Exit(1)
	}
}
