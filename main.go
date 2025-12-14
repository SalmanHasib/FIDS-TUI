package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"fids-tui/api"
	"fids-tui/config"
	"fids-tui/models"
	"fids-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	board        *ui.Board
	apiClient    *api.FlightAwareClient
	cfg          *config.Config
	airportCode  string
	loading      bool
	err          error
	inputMode    bool
	airportInput string
}

type errMsg struct {
	err error
}

type flightsMsg struct {
	flights []models.Flight
	err     error
}

// Initialization
func initialModel(airportCode string, cfg *config.Config) model {
	apiClient := api.NewFlightAwareClient(cfg.APIKey)
	airportTZ := api.GetAirportTimezone(airportCode)
	board := ui.NewBoard(airportCode, airportTZ, cfg.FlightsPerPage)

	return model{
		board:        board,
		apiClient:    apiClient,
		cfg:          cfg,
		airportCode:  airportCode,
		loading:      true,
		inputMode:    false,
		airportInput: "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchFlights(m.apiClient, m.airportCode, m.cfg.LookaheadHours, m.cfg.MaxPages),
		tickAPI(m.cfg.UpdateInterval),
		tickPageRotation(m.cfg.PageRotationInterval),
		tickAnimation(m.cfg.CharAnimationSpeed),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			// Handle input mode
			switch msg.String() {
			case "enter":
				// Validate and change airport
				newCode := strings.ToUpper(strings.TrimSpace(m.airportInput))
				if len(newCode) == 3 {
					// Validate airport code
					valid := true
					for _, r := range newCode {
						if r < 'A' || r > 'Z' {
							valid = false
							break
						}
					}
					if valid {
						m.airportCode = newCode
						m.airportInput = ""
						m.inputMode = false
						m.loading = true
						m.board.Error = ""
						airportTZ := api.GetAirportTimezone(m.airportCode)
						m.board.SetAirport(m.airportCode, airportTZ)
						m.board.SetFlightsPerPage(m.cfg.FlightsPerPage)
						return m, fetchFlights(m.apiClient, m.airportCode, m.cfg.LookaheadHours, m.cfg.MaxPages)
					}
				}
				// Invalid code, exit input mode
				m.airportInput = ""
				m.inputMode = false
				return m, nil
			case "esc":
				// Cancel input mode
				m.airportInput = ""
				m.inputMode = false
				return m, nil
			case "backspace":
				if len(m.airportInput) > 0 {
					m.airportInput = m.airportInput[:len(m.airportInput)-1]
				}
				return m, nil
			default:
				// Add character if it's a letter and we have space
				if len(m.airportInput) < 3 {
					keyStr := msg.String()
					if len(keyStr) == 1 {
						r := rune(keyStr[0])
						if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
							m.airportInput += strings.ToUpper(keyStr)
						}
					}
				}
				return m, nil
			}
		} else {
			// Normal mode
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "a":
				// Enter airport input mode
				m.inputMode = true
				m.airportInput = ""
				return m, nil
			}
		}

	case errMsg:
		m.err = msg.err
		m.board.Error = msg.err.Error()
		m.loading = false
		return m, nil

	case flightsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.board.Error = msg.err.Error()
		} else {
			m.board.Error = ""
			m.board.UpdateFlights(msg.flights)
		}
		return m, nil

	case tickAPIMsg:
		// Fetch flights on API tick
		return m, tea.Batch(
			fetchFlights(m.apiClient, m.airportCode, m.cfg.LookaheadHours, m.cfg.MaxPages),
			tickAPI(m.cfg.UpdateInterval),
		)

	case tickPageRotationMsg:
		// Rotate to next page
		m.board.NextPage()
		return m, tickPageRotation(m.cfg.PageRotationInterval)

	case tickAnimationMsg:
		// Update character animations
		m.board.Tick()
		return m, tickAnimation(m.cfg.CharAnimationSpeed)
	}

	return m, nil
}

func (m model) View() string {
	if m.inputMode {
		// Show input prompt
		prompt := fmt.Sprintf("Enter airport code (3 letters): %s_", m.airportInput)
		return fmt.Sprintf("%s\n\n%s", prompt, m.board.Render())
	}
	if m.loading && len(m.board.Flights) == 0 {
		return "Loading flights...\n"
	}
	view := m.board.Render()
	if !m.inputMode {
		// Add help text at the bottom
		help := "\nPress 'a' to change airport | 'q' to quit"
		view += help
	}
	return view
}

// Commands
type tickAPIMsg time.Time
type tickPageRotationMsg time.Time
type tickAnimationMsg time.Time

func tickAPI(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickAPIMsg(t)
	})
}

func tickPageRotation(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickPageRotationMsg(t)
	})
}

func tickAnimation(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickAnimationMsg(t)
	})
}

func fetchFlights(client *api.FlightAwareClient, airportCode string, hours int, maxPages int) tea.Cmd {
	return func() tea.Msg {
		flights, err := client.GetDepartures(airportCode, hours, maxPages)
		return flightsMsg{flights: flights, err: err}
	}
}

func main() {
	// Parse command line arguments
	var airportCode string
	flag.StringVar(&airportCode, "airport", "", "Airport code (e.g., JFK, LAX)")
	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	// Get airport code from command line, env var, or config
	if airportCode == "" {
		airportCode = cfg.AirportCode
	}
	if airportCode == "" {
		fmt.Fprintf(os.Stderr, "Error: Airport code required. Use -airport flag or set AIRPORT_CODE environment variable.\n")
		os.Exit(1)
	}

	// Validate and normalize airport code (uppercase, 3 letters)
	airportCode = strings.ToUpper(strings.TrimSpace(airportCode))
	if len(airportCode) != 3 {
		fmt.Fprintf(os.Stderr, "Error: Airport code must be 3 letters (e.g., JFK, LAX). Got: %s\n", airportCode)
		os.Exit(1)
	}

	// Validate airport code contains only letters
	for _, r := range airportCode {
		if r < 'A' || r > 'Z' {
			fmt.Fprintf(os.Stderr, "Error: Airport code must contain only letters (e.g., JFK, LAX). Got: %s\n", airportCode)
			os.Exit(1)
		}
	}

	// Validate API key
	if cfg.APIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: FLIGHTAWARE_API_KEY environment variable is required.\n")
		os.Exit(1)
	}

	// Initialize and run the program
	p := tea.NewProgram(initialModel(airportCode, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
