# Moving Circles Life Simulation

A Conway's Game of Life inspired simulation where circles move around and interact based on proximity. This project demonstrates emergent behavior through simple rules applied to moving entities.

## Features
- Dynamic movement of entities represented as circles
- Proximity-based interactions between entities
- Continuous simulation with real-time visualization
- Inspired by Conway's Game of Life rules but with added movement mechanics
- SFML-powered graphics for smooth rendering

## Prerequisites

1. Go 1.21 or later
2. SFML 2.6 library

### Installing Prerequisites

#### Go Installation
1. Download Go from [official website](https://golang.org/dl/)
2. Follow the installation instructions for your OS

#### SFML Installation
On macOS:
```bash
brew install sfml
```

On Linux:
```bash
sudo apt-get install libsfml-dev  # Debian/Ubuntu
sudo pacman -S sfml               # Arch Linux
```

On Windows:
1. Download SFML 2.6 from [SFML website](https://www.sfml-dev.org/download.php)
2. Follow the Windows installation instructions from SFML documentation

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd predator_prey_prototype
```

2. Install Go dependencies:
```bash
go mod tidy
```

## Running the Simulation

Start the simulation:
```bash
go run main.go
```

## Controls and Functionality

### Basic Controls
- Close window to exit the simulation
- The simulation runs automatically once started

### Entities
- Green circles represent living cells
- Each cell moves continuously in the environment
- Cells interact with nearby cells based on proximity rules

### Simulation Rules
- Cells move autonomously in the environment
- Interactions occur when cells come within proximity of each other
- Movement patterns and interactions create emergent behavior
- Rules are inspired by Conway's Game of Life but adapted for continuous movement

## Development

The project is structured using Go modules and SFML for graphics rendering. The main simulation logic is separated from the rendering code for better maintainability 