# FIDS-TUI

A beautiful terminal-based Flight Information Display System (FIDS) that shows real-time flight departure information using the FlightAware API.

![FIDS-TUI](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

- ✈️ **Real-time Flight Departures** - View scheduled departures from any airport
- 🎨 **Beautiful TUI** - Terminal user interface with split-flap display aesthetics
- 🔄 **Auto-refresh** - Automatically updates flight information at configurable intervals
- 📄 **Pagination** - Navigate through multiple pages of flights with automatic rotation
- 🎭 **Animations** - Smooth character animations for a retro display feel
- 🌍 **Timezone Support** - Automatically displays times in the airport's local timezone
- ⌨️ **Interactive** - Change airports on the fly with simple keyboard commands
- 🚦 **Status Indicators** - Color-coded status lights (green/yellow/orange/red) for flight status

## Prerequisites

- Go 1.24 or later
- A FlightAware API key ([Get one here](https://www.flightaware.com/commercial/aeroapi/))

## Installation

### From Source

```bash
git clone https://github.com/SalmanHasib/FIDS-TUI.git
cd FIDS-TUI
go build -o fids-tui
```

### Using Go Install

```bash
go install github.com/SalmanHasib/FIDS-TUI@latest
```

## Configuration

### Environment Variables

The application can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `FLIGHTAWARE_API_KEY` | **Required** - Your FlightAware API key | - |
| `AIRPORT_CODE` | Default airport code (3-letter IATA code) | - |
| `UPDATE_INTERVAL` | How often to fetch new flight data | `10m` |
| `PAGE_ROTATION_INTERVAL` | How often to rotate to next page | `15s` |
| `MAX_PAGES` | Maximum number of pages to fetch from API | `3` |

### Command Line Arguments

```bash
fids-tui -airport JFK
```

- `-airport`: Airport code (3-letter IATA code, e.g., JFK, LAX, LHR)

## Usage

1. Set your FlightAware API key:
   ```bash
   export FLIGHTAWARE_API_KEY="your-api-key-here"
   ```

2. Run the application:
   ```bash
   ./fids-tui -airport JFK
   ```

   Or use the default airport from environment:
   ```bash
   export AIRPORT_CODE="JFK"
   ./fids-tui
   ```

3. **Keyboard Controls:**
   - `a` - Change airport (enter a 3-letter airport code)
   - `q` or `Ctrl+C` - Quit the application

## Display Information

The FIDS board displays the following information for each flight:

- **Status** - Color-coded status indicator:
  - 🟢 Green: On Time
  - 🟡 Yellow: Taxiing / Left Gate
  - 🟠 Orange: Delayed or Taxiing / Delayed
  - 🔴 Red: Cancelled
- **Flight Number** - Airline code and flight number
- **Time** - Scheduled departure time (in airport local timezone)
- **Destination** - Destination airport code and city
- **Gate** - Gate assignment
- **Remarks** - Flight status remarks (e.g., "Delayed EST: 14:30")

## Project Structure

```
FIDS-TUI/
├── api/              # FlightAware API integration
│   ├── flightaware.go
│   └── timezone.go
├── config/           # Configuration management
│   └── config.go
├── models/           # Data models
│   └── flight.go
├── ui/               # Terminal UI components
│   ├── animation.go
│   ├── board.go
│   ├── flight_row.go
│   └── styles.go
├── main.go           # Application entry point
├── go.mod
└── go.sum
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library

## API Rate Limits

Please be aware of FlightAware API rate limits. The application is configured with reasonable defaults, but you may need to adjust `UPDATE_INTERVAL` based on your API plan.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [FlightAware](https://www.flightaware.com/) for providing the AeroAPI
- [Charm](https://charm.sh/) for the amazing Bubble Tea TUI framework
- Inspired by classic airport split-flap displays

## Troubleshooting

### "API authentication failed"
- Verify your `FLIGHTAWARE_API_KEY` environment variable is set correctly
- Check that your API key is valid and has not expired

### "Airport not found"
- Ensure you're using a valid 3-letter IATA airport code
- Some smaller airports may not be available in the FlightAware database

### No flights displayed
- The airport may not have any scheduled departures in the configured time window
- Try adjusting the `UPDATE_INTERVAL` or check the airport code


## Screenshots

![FIDS-TUI Screenshot](/screenshots/fids-tui.png)

*Example: FIDS-TUI displaying departures from LAX airport*

---

Made with ❤️, Go, and Cursor  
_An experiment on working with AI tools to quickly prototype a TUI. A more comprehensive (and hand built) version is in the works._
