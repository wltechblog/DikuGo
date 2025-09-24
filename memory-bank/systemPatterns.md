# System Patterns

## Architecture Overview
DikuGo follows a modular architecture with clear separation of concerns:

```
DikuGo
├── Network Layer (telnet server, client connections)
├── Game Layer (main game loop, command processing)
├── World Layer (rooms, zones, objects, characters)
├── Combat System (combat mechanics, damage calculation)
├── Command System (command parsing and execution)
└── Storage Layer (loading/saving game data)
```

## Key Design Patterns
1. **Command Pattern**: Used for implementing player commands
2. **Factory Pattern**: Used for creating game objects from prototypes
3. **Observer Pattern**: Used for event handling in the game world
4. **Repository Pattern**: Used for data access and persistence

## Component Relationships
- The Network Layer handles client connections and passes input to the Game Layer
- The Game Layer processes commands through the Command System
- The World Layer maintains the state of the game world
- The Combat System handles combat interactions between characters
- The Storage Layer loads and saves game data

## Technical Decisions
1. **Go Language**: Chosen for its simplicity, performance, and concurrency support
2. **Package Structure**: Organized by domain (world, combat, command, etc.)
3. **Data Format**: Compatible with original DikuMUD data files
4. **Concurrency Model**: Uses Go's goroutines and channels for handling multiple players

## Combat System
The combat system follows the original DikuMUD mechanics:
- Turn-based combat with automatic rounds
- Hit and damage calculations based on character stats
- Weapon and armor effects on combat outcomes
- Special attacks and defenses for certain creatures

## Character System
Characters (both players and NPCs) share a common structure:
- Core attributes (strength, intelligence, etc.)
- Combat stats (hit points, armor class, etc.)
- Inventory and equipment
- Position and location in the world