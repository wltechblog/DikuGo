package command

import (
	"github.com/wltechblog/DikuGo/pkg/combat"
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// InitRegistry initializes the command registry with all available commands
func InitRegistry(w *world.World, combatManager *combat.Manager) *Registry {
	registry := NewRegistry()

	// Register commands
	registry.Register(&LookCommand{})
	registry.Register(&SayCommand{})
	registry.Register(&WhoCommand{Characters: w.GetCharacters()})
	registry.Register(&QuitCommand{})
	registry.Register(&KillCommand{CombatManager: combatManager})
	registry.Register(&GetCommand{})
	registry.Register(&DropCommand{})
	registry.Register(&InventoryCommand{})
	registry.Register(&WearCommand{})
	registry.Register(&RemoveCommand{})
	registry.Register(&EquipmentCommand{})
	registry.Register(&ListCommand{})
	registry.Register(&BuyCommand{})
	registry.Register(&SellCommand{})

	// Register admin commands
	registry.Register(&GotoCommand{})
	registry.Register(&RstatCommand{})
	registry.Register(&ValidateRoomsCommand{})
	registry.Register(&TestExitsCommand{})
	registry.Register(&ResetZoneCommand{})

	// Register movement commands
	registry.Register(&MovementCommand{direction: types.DIR_NORTH})
	registry.Register(&MovementCommand{direction: types.DIR_EAST})
	registry.Register(&MovementCommand{direction: types.DIR_SOUTH})
	registry.Register(&MovementCommand{direction: types.DIR_WEST})
	registry.Register(&MovementCommand{direction: types.DIR_UP})
	registry.Register(&MovementCommand{direction: types.DIR_DOWN})

	// Register help command (needs registry)
	helpCmd := &HelpCommand{Registry: registry}
	registry.Register(helpCmd)

	return registry
}
