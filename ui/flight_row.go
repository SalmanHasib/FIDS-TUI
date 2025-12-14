package ui

import (
	"fids-tui/models"
	"fmt"
)

// FlightRow represents an animated flight row
type FlightRow struct {
	Flight          *models.Flight
	StatusAnim      *AnimatedText
	FlightNumAnim   *AnimatedText
	TimeAnim        *AnimatedText
	DestinationAnim *AnimatedText
	GateAnim        *AnimatedText
	RemarksAnim     *AnimatedText
}

// NewFlightRow creates a new flight row with animations
func NewFlightRow(flight *models.Flight) *FlightRow {
	row := &FlightRow{
		Flight:          flight,
		StatusAnim:      NewAnimatedText(1),
		FlightNumAnim:   NewAnimatedText(8), // Full flight number with airline code
		TimeAnim:        NewAnimatedText(8), // HH:MM format
		DestinationAnim: NewAnimatedText(20),
		GateAnim:        NewAnimatedText(6),
		RemarksAnim:     NewAnimatedText(20),
	}

	// Initialize animated text with flight data if available
	if flight != nil {
		statusChar := getStatusChar(flight.Status)
		row.StatusAnim.Update(statusChar)

		flightNum := truncate(flight.FlightNumber, 8)
		row.FlightNumAnim.Update(flightNum)

		timeStr := flight.ScheduledDeparture.Format("15:04")
		row.TimeAnim.Update(timeStr)

		dest := truncate(flight.GetDestination(), 20)
		row.DestinationAnim.Update(dest)

		gate := truncate(flight.Gate, 6)
		if gate == "" {
			gate = "     "
		}
		row.GateAnim.Update(gate)

		remarks := truncate(string(flight.Remarks), 20)
		row.RemarksAnim.Update(remarks)
	}

	return row
}

// Update updates the flight data and triggers animations
func (fr *FlightRow) Update(flight *models.Flight) {
	fr.Flight = flight

	// Update status indicator
	statusChar := getStatusChar(flight.Status)
	fr.StatusAnim.Update(statusChar)

	// Update flight number (already includes airline code prefix, e.g., "BA114")
	flightNum := truncate(flight.FlightNumber, 8)
	fr.FlightNumAnim.Update(flightNum)

	// Update departure time (HH:MM format)
	timeStr := flight.ScheduledDeparture.Format("15:04")
	fr.TimeAnim.Update(timeStr)

	// Update destination
	dest := truncate(flight.GetDestination(), 20)
	fr.DestinationAnim.Update(dest)

	// Update gate
	gate := truncate(flight.Gate, 6)
	if gate == "" {
		gate = "     "
	}
	fr.GateAnim.Update(gate)

	// Update remarks
	remarks := truncate(string(flight.Remarks), 20)
	fr.RemarksAnim.Update(remarks)
}

// Tick updates all animations
func (fr *FlightRow) Tick() {
	fr.StatusAnim.Tick()
	fr.FlightNumAnim.Tick()
	fr.TimeAnim.Tick()
	fr.DestinationAnim.Tick()
	fr.GateAnim.Tick()
	fr.RemarksAnim.Tick()
}

// Render renders the flight row with split-flap styling
func (fr *FlightRow) Render(styles *SplitFlapStyles) string {
	if fr.Flight == nil {
		// Empty row (68 characters to match row width)
		return styles.Text.Render("                                                                    ")
	}

	statusColor := fr.Flight.GetStatusColor()

	// Render status with color
	statusText := fr.StatusAnim.Render()
	statusStyle := styles.StatusLight(statusColor)
	statusRendered := statusStyle.Render(statusText)

	// Render other fields
	flightNum := styles.Text.Render(fr.FlightNumAnim.Render())
	timeStr := styles.Text.Render(fr.TimeAnim.Render())
	destination := styles.Text.Render(fr.DestinationAnim.Render())
	gate := styles.Text.Render(fr.GateAnim.Render())
	remarks := styles.Text.Render(fr.RemarksAnim.Render())

	// Combine with proper spacing
	return fmt.Sprintf("%s %s %s %s %s %s",
		statusRendered,
		flightNum,
		timeStr,
		destination,
		gate,
		remarks,
	)
}

// getStatusChar returns a character icon for the status
// Uses simple ASCII-compatible characters that work in all terminals
// The character will be colored by the StatusLight style
func getStatusChar(status models.FlightStatus) string {
	switch status {
	case models.StatusOnTime:
		return "*" // Asterisk - green (on time)
	case models.StatusDelayed:
		return "!" // Exclamation - orange (delayed)
	case models.StatusTaxiingLeftGate:
		return ">" // Greater than - yellow (taxiing)
	case models.StatusTaxiingDelayed:
		return ">" // Greater than - orange (taxiing delayed)
	case models.StatusCancelled:
		return "X" // X - red (cancelled)
	default:
		return " " // Space for unknown status
	}
}

// truncate truncates a string to the specified length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
