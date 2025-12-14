package ui

import (
	"fids-tui/models"
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Board manages the flight board display
type Board struct {
	Flights        []*FlightRow
	CurrentPage    int
	TotalPages     int
	AirportCode    string
	AirportTZ      *time.Location
	FlightsPerPage int
	Error          string
	Styles         *SplitFlapStyles
}

// NewBoard creates a new flight board
func NewBoard(airportCode string, airportTZ *time.Location, flightsPerPage int) *Board {
	return &Board{
		Flights:        make([]*FlightRow, 0),
		CurrentPage:    0,
		TotalPages:     1,
		AirportCode:    airportCode,
		AirportTZ:      airportTZ,
		FlightsPerPage: flightsPerPage,
		Styles:         NewSplitFlapStyles(),
	}
}

// UpdateFlights updates the flight list and creates/updates flight rows
func (b *Board) UpdateFlights(flights []models.Flight) {
	// Convert departure times to airport local time
	for i := range flights {
		if b.AirportTZ != nil {
			flights[i].ScheduledDeparture = flights[i].ScheduledDeparture.In(b.AirportTZ)
			// Convert estimated departure time to local time if present
			if flights[i].EstimatedDeparture != nil {
				localEst := flights[i].EstimatedDeparture.In(b.AirportTZ)
				flights[i].EstimatedDeparture = &localEst
				// Update remarks for delayed flights with estimated time
				if flights[i].Status == models.StatusDelayed {
					estTimeStr := localEst.Format("15:04")
					flights[i].Remarks = models.Remarks(fmt.Sprintf("Delayed EST: %s", estTimeStr))
				}
			}
		}
	}

	// Sort flights by departure time ascending
	sort.Slice(flights, func(i, j int) bool {
		return flights[i].ScheduledDeparture.Before(flights[j].ScheduledDeparture)
	})

	// Create a map of existing flights by flight number
	existingMap := make(map[string]*FlightRow)
	for _, row := range b.Flights {
		if row.Flight != nil {
			key := row.Flight.FlightNumber
			existingMap[key] = row
		}
	}

	// Update or create flight rows
	newRows := make([]*FlightRow, 0, len(flights))
	for i := range flights {
		flight := &flights[i]
		key := flight.FlightNumber

		if existingRow, exists := existingMap[key]; exists {
			// Update existing row (triggers animation)
			existingRow.Update(flight)
			newRows = append(newRows, existingRow)
		} else {
			// Create new row
			row := NewFlightRow(flight)
			newRows = append(newRows, row)
		}
	}

	b.Flights = newRows
	b.updatePagination()
}

// updatePagination updates pagination info
func (b *Board) updatePagination() {
	flightsPerPage := b.FlightsPerPage
	if flightsPerPage <= 0 {
		flightsPerPage = 10 // Default fallback
	}
	totalFlights := len(b.Flights)

	if totalFlights == 0 {
		b.TotalPages = 1
		b.CurrentPage = 0
		return
	}

	b.TotalPages = (totalFlights + flightsPerPage - 1) / flightsPerPage
	if b.CurrentPage >= b.TotalPages {
		b.CurrentPage = b.TotalPages - 1
	}
	if b.CurrentPage < 0 {
		b.CurrentPage = 0
	}
}

// NextPage moves to the next page
func (b *Board) NextPage() {
	b.updatePagination()
	b.CurrentPage = (b.CurrentPage + 1) % b.TotalPages
}

// GetCurrentPageFlights returns flights for the current page
// Always returns exactly flightsPerPage rows, filling with empty rows if needed
func (b *Board) GetCurrentPageFlights() []*FlightRow {
	flightsPerPage := b.FlightsPerPage
	if flightsPerPage <= 0 {
		flightsPerPage = 10 // Default fallback
	}
	start := b.CurrentPage * flightsPerPage
	end := start + flightsPerPage

	result := make([]*FlightRow, flightsPerPage)

	// Copy actual flights
	actualEnd := end
	if start >= len(b.Flights) {
		actualEnd = start
	} else if end > len(b.Flights) {
		actualEnd = len(b.Flights)
	}

	copyCount := 0
	if start < len(b.Flights) {
		copyCount = actualEnd - start
		copy(result, b.Flights[start:actualEnd])
	}

	// Fill remaining slots with empty rows
	for i := copyCount; i < flightsPerPage; i++ {
		result[i] = NewFlightRow(nil)
	}

	return result
}

// Tick updates all flight row animations
func (b *Board) Tick() {
	// Update all flights
	for _, row := range b.Flights {
		row.Tick()
	}
}

// Render renders the entire board
func (b *Board) Render() string {
	var sections []string

	// Airport header
	airportHeader := b.renderAirportHeader()
	sections = append(sections, airportHeader)

	// Error message if any
	if b.Error != "" {
		errorMsg := b.Styles.Error.Render("ERROR: " + b.Error)
		sections = append(sections, errorMsg)
	}

	// Table header
	header := b.renderHeader()
	sections = append(sections, header)

	// Flight rows for current page (always shows flightsPerPage rows)
	pageFlights := b.GetCurrentPageFlights()
	for _, row := range pageFlights {
		if row != nil {
			rowStr := row.Render(b.Styles)
			sections = append(sections, rowStr)
		}
	}

	// Page info
	pageInfo := b.renderPageInfo()
	sections = append(sections, pageInfo)

	// Combine all sections
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return b.Styles.Background.Render(content)
}

// renderAirportHeader renders the airport code header
func (b *Board) renderAirportHeader() string {
	label := fmt.Sprintf("DEPARTURES - %s", b.AirportCode)
	return b.Styles.AirportLabel.Render(label)
}

// SetAirport updates the airport code and timezone
func (b *Board) SetAirport(airportCode string, airportTZ *time.Location) {
	b.AirportCode = airportCode
	b.AirportTZ = airportTZ
	// Reset to first page when airport changes
	b.CurrentPage = 0
}

// SetFlightsPerPage updates the flights per page setting
func (b *Board) SetFlightsPerPage(flightsPerPage int) {
	b.FlightsPerPage = flightsPerPage
}

// renderHeader renders the table header
func (b *Board) renderHeader() string {
	status := b.Styles.Header.Render("S")
	flightNum := b.Styles.Header.Render(fmt.Sprintf("%-8s", "FLIGHT"))
	time := b.Styles.Header.Render(fmt.Sprintf("%-8s", "TIME"))
	destination := b.Styles.Header.Render(fmt.Sprintf("%-20s", "DESTINATION"))
	gate := b.Styles.Header.Render(fmt.Sprintf("%-6s", "GATE"))
	remarks := b.Styles.Header.Render(fmt.Sprintf("%-20s", "REMARKS"))

	return fmt.Sprintf("%s %s %s %s %s %s", status, flightNum, time, destination, gate, remarks)
}

// renderPageInfo renders pagination information
func (b *Board) renderPageInfo() string {
	if b.TotalPages <= 1 {
		return ""
	}

	totalFlights := len(b.Flights)
	flightsPerPage := b.FlightsPerPage
	if flightsPerPage <= 0 {
		flightsPerPage = 10 // Default fallback
	}
	start := b.CurrentPage*flightsPerPage + 1
	end := (b.CurrentPage + 1) * flightsPerPage
	if end > totalFlights {
		end = totalFlights
	}
	info := fmt.Sprintf("Page %d/%d (%d-%d of %d)",
		b.CurrentPage+1, b.TotalPages, start, end, totalFlights)
	return b.Styles.PageInfo.Render(info)
}
