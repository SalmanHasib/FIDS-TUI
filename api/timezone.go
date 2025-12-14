package api

import (
	"time"
)

// GetAirportTimezone returns the IANA timezone for a given airport code
// This is a simplified mapping - in production, you might want a more comprehensive database
func GetAirportTimezone(airportCode string) *time.Location {
	// Map of airport codes to IANA timezone identifiers
	timezoneMap := map[string]string{
		// US East Coast
		"JFK": "America/New_York",
		"LGA": "America/New_York",
		"EWR": "America/New_York",
		"BOS": "America/New_York",
		"MIA": "America/New_York",
		"ATL": "America/New_York",
		"CLT": "America/New_York",
		"DCA": "America/New_York",
		"IAD": "America/New_York",
		"PHL": "America/New_York",
		"BWI": "America/New_York",

		// US Central
		"ORD": "America/Chicago",
		"MDW": "America/Chicago",
		"DFW": "America/Chicago",
		"IAH": "America/Chicago",
		"MSP": "America/Chicago",
		"STL": "America/Chicago",
		"DTW": "America/Detroit",
		"CLE": "America/New_York",

		// US Mountain
		"DEN": "America/Denver",
		"PHX": "America/Phoenix",
		"SLC": "America/Denver",

		// US West Coast
		"LAX": "America/Los_Angeles",
		"SFO": "America/Los_Angeles",
		"SAN": "America/Los_Angeles",
		"SEA": "America/Los_Angeles",
		"PDX": "America/Los_Angeles",
		"LAS": "America/Los_Angeles",

		// Europe
		"LHR": "Europe/London",
		"LGW": "Europe/London",
		"CDG": "Europe/Paris",
		"FRA": "Europe/Berlin",
		"AMS": "Europe/Amsterdam",
		"MAD": "Europe/Madrid",
		"FCO": "Europe/Rome",
		"ZUR": "Europe/Zurich",
		"VIE": "Europe/Vienna",
		"CPH": "Europe/Copenhagen",
		"ARN": "Europe/Stockholm",
		"OSL": "Europe/Oslo",
		"HEL": "Europe/Helsinki",
		"DUB": "Europe/Dublin",
		"BRU": "Europe/Brussels",

		// Asia
		"NRT": "Asia/Tokyo",
		"HND": "Asia/Tokyo",
		"ICN": "Asia/Seoul",
		"PEK": "Asia/Shanghai",
		"PVG": "Asia/Shanghai",
		"HKG": "Asia/Hong_Kong",
		"SIN": "Asia/Singapore",
		"BKK": "Asia/Bangkok",
		"DXB": "Asia/Dubai",
		"AUH": "Asia/Dubai",
		"IST": "Europe/Istanbul",

		// Middle East
		"TLV": "Asia/Jerusalem",
		"CAI": "Africa/Cairo",
		"JED": "Asia/Riyadh",
		"RUH": "Asia/Riyadh",

		// Australia
		"SYD": "Australia/Sydney",
		"MEL": "Australia/Melbourne",
		"BNE": "Australia/Brisbane",
		"PER": "Australia/Perth",

		// Canada
		"YYZ": "America/Toronto",
		"YVR": "America/Vancouver",
		"YUL": "America/Toronto",
		"YYC": "America/Edmonton",

		// South America
		"GRU": "America/Sao_Paulo",
		"GIG": "America/Sao_Paulo",
		"EZE": "America/Argentina/Buenos_Aires",
		"SCL": "America/Santiago",
		"LIM": "America/Lima",
		"BOG": "America/Bogota",
		"MEX": "America/Mexico_City",
	}

	// Look up timezone
	if tzName, ok := timezoneMap[airportCode]; ok {
		if loc, err := time.LoadLocation(tzName); err == nil {
			return loc
		}
	}

	// Default to UTC if not found
	return time.UTC
}
