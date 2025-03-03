# Moving Circles Life Simulation

A Conway's Game of Life inspired simulation where circles move around and interact based on proximity.

## Prerequisites

1. Go 1.21 or later
2. SFML 2.6 library

### Installing SFML on macOS:
```bash
brew install sfml
```

## Running the Simulation

1. Install dependencies:
```bash
go mod tidy
```

2. Run the simulation:
```bash
go run main.go
```

## Controls
- Close window to exit
- Green circles represent living cells
- Circles move continuously and interact based on proximity
- Rules are inspired by Conway's Game of Life but with movement 