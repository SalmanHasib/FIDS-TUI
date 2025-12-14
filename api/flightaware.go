package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"fids-tui/models"
)

const (
	flightAwareBaseURL = "https://aeroapi.flightaware.com/aeroapi"
)

// FlightAwareClient handles API interactions with FlightAware
type FlightAwareClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewFlightAwareClient creates a new FlightAware API client
func NewFlightAwareClient(apiKey string) *FlightAwareClient {
	return &FlightAwareClient{
		APIKey:  apiKey,
		BaseURL: flightAwareBaseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AeroAPIDeparture represents a departure from FlightAware API
type AeroAPIDeparture struct {
	Ident        string     `json:"ident"`
	FaFlightID   string     `json:"fa_flight_id"`
	Operator     string     `json:"operator"`
	OperatorIata string     `json:"operator_iata"`
	FlightNumber string     `json:"flight_number"`
	Origin       *Airport   `json:"origin"`
	Destination  *Airport   `json:"destination"`
	Departure    *TimeInfo  `json:"departure"`
	ScheduledOut *time.Time `json:"scheduled_out"`
	EstimatedOut *time.Time `json:"estimated_out"`
	ActualOut    *time.Time `json:"actual_out"`
	Status       string     `json:"status"`
	Gate         string     `json:"gate_origin"`
	BaggageClaim string     `json:"baggage_claim"`
	Remarks      string     `json:"remarks"`
}

// Airport represents airport information
type Airport struct {
	Code     string `json:"code"`
	CodeIata string `json:"code_iata"`
	CodeIcao string `json:"code_icao"`
	City     string `json:"city"`
}

// TimeInfo represents time information
type TimeInfo struct {
	Scheduled time.Time `json:"scheduled"`
	Estimated time.Time `json:"estimated"`
	Actual    time.Time `json:"actual"`
}

// AeroAPIResponse represents the response from FlightAware API
type AeroAPIResponse struct {
	ScheduledDepartures []AeroAPIDeparture `json:"scheduled_departures"`
}

// GetDepartures fetches scheduled departures for an airport within the specified hours
// Uses the scheduled_departures endpoint which defaults to 2 hours before current time
// and excludes flights that have already departed (en route)
func (c *FlightAwareClient) GetDepartures(airportCode string, hours int, maxPages int) ([]models.Flight, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/airports/%s/flights/scheduled_departures", c.BaseURL, airportCode)
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Build query parameters
	params := url.Values{}

	// Only add end time parameter if hours is specified (greater than 0)
	if hours > 0 {
		endTime := time.Now().Add(time.Duration(hours) * time.Hour)
		endTimeISO8601 := endTime.Format(time.RFC3339)
		params.Add("end", endTimeISO8601)
	}

	// Only add max_pages parameter if it's greater than 1 (default is 1)
	if maxPages > 1 {
		params.Add("max_pages", fmt.Sprintf("%d", maxPages))
	}

	// Set query parameters if any were added
	if len(params) > 0 {
		reqURL.RawQuery = params.Encode()
	}

	fullURL := reqURL.String()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-apikey", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("API authentication failed: check your FLIGHTAWARE_API_KEY")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("airport not found: %s", airportCode)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if len(body) == 0 {
		return []models.Flight{}, nil // Empty response is valid, just no flights
	}

	var apiResp AeroAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Filter and convert to our Flight model
	flights := make([]models.Flight, 0)
	now := time.Now()
	// scheduled_departures endpoint defaults to 2 hours before current time
	// We only need to filter by the future cutoff time if hours is specified
	var cutoffTime *time.Time
	if hours > 0 {
		ct := now.Add(time.Duration(hours) * time.Hour)
		cutoffTime = &ct
	}
	maxFlights := 50

	for _, dep := range apiResp.ScheduledDepartures {
		if len(flights) >= maxFlights {
			break
		}

		// Try to get scheduled time from various possible fields
		var scheduled time.Time
		var hasScheduled bool

		// First try the nested departure.scheduled field
		if dep.Departure != nil && !dep.Departure.Scheduled.IsZero() {
			scheduled = dep.Departure.Scheduled
			hasScheduled = true
		} else if dep.ScheduledOut != nil && !dep.ScheduledOut.IsZero() {
			// Try scheduled_out field
			scheduled = *dep.ScheduledOut
			hasScheduled = true
		} else if dep.EstimatedOut != nil && !dep.EstimatedOut.IsZero() {
			// Fall back to estimated_out if scheduled is not available
			scheduled = *dep.EstimatedOut
			hasScheduled = true
		}

		if !hasScheduled {
			continue
		}

		if scheduled.IsZero() {
			continue
		}

		// Filter flights departing within the specified future window (if cutoff is set)
		// scheduled_departures endpoint already excludes en route flights and includes past 2 hours
		if cutoffTime != nil && scheduled.After(*cutoffTime) {
			continue
		}

		flight := c.convertToFlight(dep, scheduled)
		flights = append(flights, flight)
	}

	return flights, nil
}

// convertToFlight converts an AeroAPI departure to our Flight model
func (c *FlightAwareClient) convertToFlight(dep AeroAPIDeparture, scheduled time.Time) models.Flight {
	// Use operator_iata (2-letter) for airline code
	airlineCode := dep.OperatorIata
	if airlineCode == "" {
		// Fall back to operator (3-letter) if IATA not available
		airlineCode = dep.Operator
	}

	// Use operator (3-letter code) as airline name, or fall back to IATA code
	airlineName := dep.Operator
	if airlineName == "" {
		airlineName = dep.OperatorIata
	}
	if airlineName == "" {
		airlineName = "UNK" // Unknown
	}

	// Get flight number (just the number part)
	flightNumber := dep.FlightNumber
	if flightNumber == "" {
		// Try to extract number from ident (e.g., "BAW114" -> "114")
		// This is a fallback, but flight_number should be available
		flightNumber = dep.Ident
	}

	// Prepend airline code to flight number (e.g., "BA" + "114" = "BA114")
	fullFlightNumber := airlineCode + " " + flightNumber

	flight := models.Flight{
		AirlineCode:        airlineCode,
		AirlineName:        airlineName,
		FlightNumber:       fullFlightNumber,
		ScheduledDeparture: scheduled,
	}

	if dep.Destination != nil {
		flight.DestinationCode = dep.Destination.CodeIata
		if flight.DestinationCode == "" {
			flight.DestinationCode = dep.Destination.Code
		}
		flight.DestinationCity = dep.Destination.City
	}

	flight.Gate = dep.Gate

	// Determine status and remarks based on API status
	status := dep.Status
	remarks := dep.Remarks

	// Map status to our enum
	switch {
	case status == "Cancelled" || remarks == "Cancelled":
		flight.Status = models.StatusCancelled
		flight.Remarks = models.RemarksCancelled
	case status == "Taxiing / Delayed" || remarks == "Taxiing / Delayed":
		flight.Status = models.StatusTaxiingDelayed
		flight.Remarks = models.RemarksTaxiingDelayed
	case status == "Taxiing / Left Gate" || remarks == "Taxiing / Left Gate":
		flight.Status = models.StatusTaxiingLeftGate
		flight.Remarks = models.RemarksTaxiingLeftGate
	case status == "Scheduled / Delayed" || status == "Delayed" || remarks == "Delayed":
		flight.Status = models.StatusDelayed
		// Check if there's an estimated departure time
		if dep.EstimatedOut != nil && !dep.EstimatedOut.IsZero() {
			flight.EstimatedDeparture = dep.EstimatedOut
		} else if dep.Departure != nil && !dep.Departure.Estimated.IsZero() {
			flight.EstimatedDeparture = &dep.Departure.Estimated
		}
		// Remarks will be set in UpdateFlights after timezone conversion
		flight.Remarks = models.RemarksDelayed
	default:
		flight.Status = models.StatusOnTime
		flight.Remarks = models.RemarksOnTime
	}

	return flight
}
