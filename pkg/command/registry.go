package command

import (
	"github.com/wltechblog/DikuGo/pkg/types"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// InitRegistry initializes the command registry with all available commands
func InitRegistry(w *world.World, combatManager CombatManagerInterface) *Registry {
	registry := NewRegistry()

	// Register commands
	registry.Register(&LookCommand{})
	registry.Register(&SayCommand{})
	registry.Register(&WhoCommand{Characters: w.GetCharacters()})
	registry.Register(&QuitCommand{})
	registry.Register(&KillCommand{CombatManager: combatManager})
	registry.Register(&GetCommand{})
	registry.Register(&DropCommand{})
	registry.Register(&PutCommand{})
	registry.Register(&ExamineCommand{})
	registry.Register(&InventoryCommand{})
	registry.Register(&WearCommand{})
	registry.Register(&RemoveCommand{})
	registry.Register(&EquipmentCommand{})
	registry.Register(&ListCommand{})
	registry.Register(&BuyCommand{})
	registry.Register(&SellCommand{})
	registry.Register(&ScoreCommand{})
	registry.Register(&PromptCommand{})
	registry.Register(&TimeCommand{})

	// Register magic commands
	registry.Register(&CastCommand{World: w})
	registry.Register(&QuaffCommand{World: w})
	registry.Register(&ReciteCommand{World: w})
	registry.Register(&UseCommand{World: w})

	// Register admin commands
	registry.Register(&GotoCommand{})
	registry.Register(&RstatCommand{})
	registry.Register(&MobstatCommand{})
	registry.Register(&ShopstatCommand{})
	registry.Register(&ResetShopsCommand{})
	registry.Register(&CheckShopsCommand{})
	registry.Register(&WhereMobCommand{})
	registry.Register(&CheckMobsCommand{})
	registry.Register(&ExamineMobCommand{})
	registry.Register(&TestMobParserCommand{})

	// Register combat skill commands
	registry.Register(&BashCommand{CombatManager: combatManager})
	registry.Register(&RescueCommand{CombatManager: combatManager})
	registry.Register(&BackstabCommand{CombatManager: combatManager})
	registry.Register(&StealCommand{CombatManager: combatManager})
	registry.Register(&CreateBakerCommand{})
	registry.Register(&AddBakerCommand{})
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
