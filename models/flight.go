package models

import "time"

// FlightStatus represents the status of a flight
type FlightStatus int

const (
	StatusOnTime FlightStatus = iota
	StatusDelayed
	StatusTaxiingLeftGate
	StatusTaxiingDelayed
	StatusCancelled
)

// String returns the string representation of the flight status
func (s FlightStatus) String() string {
	switch s {
	case StatusOnTime:
		return "On Time"
	case StatusDelayed:
		return "Delayed"
	case StatusTaxiingLeftGate:
		return "Taxiing / Left Gate"
	case StatusTaxiingDelayed:
		return "Taxiing / Delayed"
	case StatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// Remarks represents the remarks/status message for a flight
type Remarks string

const (
	RemarksOnTime          Remarks = "On Time"
	RemarksDelayed         Remarks = "Delayed"
	RemarksTaxiingLeftGate Remarks = "Taxiing / Left Gate"
	RemarksTaxiingDelayed  Remarks = "Taxiing / Delayed"
	RemarksCancelled       Remarks = "Cancelled"
)

// Flight represents a flight departure
type Flight struct {
	Status             FlightStatus
	AirlineCode        string // 2-letter IATA code
	AirlineName        string // Full airline name/operator code
	FlightNumber       string // Full flight number with airline code prefix
	DestinationCode    string
	DestinationCity    string
	Gate               string
	Remarks            Remarks
	ScheduledDeparture time.Time
	EstimatedDeparture *time.Time // Estimated departure time (for delayed flights)
}

// GetStatusColor returns the color code for the status light
func (f *Flight) GetStatusColor() string {
	switch f.Status {
	case StatusOnTime:
		return "green"
	case StatusDelayed:
		return "orange"
	case StatusTaxiingLeftGate:
		return "yellow" // Yellow for taxiing
	case StatusTaxiingDelayed:
		return "orange" // Orange for delayed taxiing
	case StatusCancelled:
		return "red"
	default:
		return "white"
	}
}

// GetDestination returns formatted destination string (code + city)
func (f *Flight) GetDestination() string {
	if f.DestinationCity != "" {
		return f.DestinationCode + " " + f.DestinationCity
	}
	return f.DestinationCode
}
