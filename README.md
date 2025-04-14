# DikuGo

This project aims to produce a faithful re-implementation of DikuMud Gamma in Go, taking advantage of modern language features and machine capacity. The aim is to be fully compatible with DikuMud Gamma zone files and provide an identical gameplay experience, while being more maintainable and extensible.

## Goals

- Create a modern, maintainable Go implementation of DikuMUD
- Maintain 100% compatibility with existing world/player data files
- Use modern Go features (goroutines, channels, mutexes) appropriately
- Make the codebase easily extensible for new features
- Improve stability and crash resistance
- Abstract storage layer to allow future modernization
- Keep initial scope focused on core functionality matching original
- Ensure proper concurrent access and data safety
- Provide clear extension points for future enhancements

## Project Structure

```
├── cmd/
│   └── dikugo/         # Main executable
├── pkg/
│   ├── config/         # Configuration handling
│   ├── game/           # Core game logic
│   ├── network/        # Network handling
│   ├── storage/        # Data storage abstraction
│   ├── utils/          # Utility functions
│   └── world/          # Game world and entities
├── old/                # Original DikuMUD code and data
│   ├── lib/            # Original game data files
│   └── memory-bank/    # Project documentation
├── config.yaml         # Configuration file
└── go.mod              # Go module definition
```

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Building

```bash
go build -o dikugo ./cmd/dikugo
```

### Running

```bash
./dikugo
```

Or with a custom config file:

```bash
./dikugo -config=custom_config.yaml
```

## Development Status

This project is in the initial implementation phase. The core architecture and package structure have been defined, but implementation of game functionality is ongoing.

## License

This project is licensed under the same terms as the original DikuMUD code.
