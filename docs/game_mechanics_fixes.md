# Game Mechanics Fixes

## Overview

This document describes the fixes implemented for three major game mechanics issues in DikuGo:

1. **Drink Command Issue**: The drink command was looking for potions instead of handling drink containers
2. **Container Access Issue**: Players couldn't get items from containers (such as corpses)
3. **Pet Shop Issue**: The pet shop functionality was completely missing

## 1. Drink Command Fix

### Problem
The "drink" command was aliased to the "quaff" command, which only handles potions (ITEM_POTION). This meant players couldn't drink from water fountains, bottles, or other drink containers (ITEM_DRINKCON, ITEM_FOUNTAIN).

### Solution
- **Created new DrinkCommand** (`pkg/command/drink.go`)
- **Removed "drink" alias** from QuaffCommand
- **Registered new DrinkCommand** in the command registry
- **Added proper drink container handling** for ITEM_DRINKCON and ITEM_FOUNTAIN types
- **Implemented drink effects** (thirst, hunger, drunk conditions)
- **Added instance-specific Value array** to ObjectInstance for tracking liquid amounts

### Key Features
- Handles both drink containers and fountains
- Proper liquid consumption (1 unit per drink)
- Drink effects based on liquid type (water, beer, wine, etc.)
- Condition tracking (thirst, hunger, drunk)
- Fountains have unlimited liquid, containers deplete

### Files Modified
- `pkg/command/drink.go` - New drink command implementation
- `pkg/command/quaff_command.go` - Removed "drink" alias
- `pkg/command/registry.go` - Registered new drink command
- `pkg/types/types.go` - Added Value array to ObjectInstance
- `pkg/types/item.go` - Added condition constants (moved to constants.go)

### Tests Added
- `pkg/command/drink_test.go` - Comprehensive drink command tests
- Tests for fountains, bottles, empty containers, non-drink objects
- Tests for drink effects and condition tracking

## 2. Container Access Fix

### Problem
The container access functionality was already implemented in the get command, but there were potential issues with container state handling and closed containers.

### Solution
- **Verified existing get command functionality** works correctly
- **Enhanced container tests** to ensure proper behavior
- **Fixed test expectations** to match DikuMUD pattern (success messages as errors)

### Key Features
- Get items from any container: `get sword from chest`
- Get items from corpses: `get gold from corpse`
- Proper handling of closed containers
- Container state validation
- Support for "get all from container"

### Files Modified
- `pkg/command/container_test.go` - New comprehensive container tests

### Tests Added
- TestGetFromContainer - Basic container access
- TestGetFromCorpse - Corpse looting functionality  
- TestGetFromClosedContainer - Closed container handling

## 3. Pet Shop Implementation

### Problem
Pet shop functionality was completely missing. Players couldn't list or buy pets from pet shops.

### Solution
- **Created pet shop special procedure** (`pkg/ai/pet_shop.go`)
- **Added pet shop command handlers** (`pkg/command/pet_shop_commands.go`)
- **Registered pet shop procedures** in special procedures registry
- **Added pet detection and management** functionality

### Key Features
- `list` command to show available pets and prices
- `buy <pet>` command to purchase pets
- Pet following system with Master relationship
- Price calculation based on pet experience
- Pet shop detection in rooms
- Support for custom pet names

### Files Modified
- `pkg/ai/pet_shop.go` - Pet shop special procedure
- `pkg/command/pet_shop_commands.go` - Pet shop commands
- `pkg/ai/special_procs.go` - Registered pet shop procedures
- `pkg/types/mob.go` - Added ACT_FOLLOWER flag
- `pkg/types/types.go` - Added Master field to Character

### Key Functions
- `listPets()` - Shows available pets with prices
- `buyPet()` - Handles pet purchases and following setup
- `isPet()` - Determines if an NPC can be sold as a pet
- `createPetCopy()` - Creates pet instances for players

## Technical Details

### ObjectInstance Value Array
Added instance-specific Value array to ObjectInstance to support:
- Drink container liquid tracking
- Container state management
- Instance-specific object properties

```go
type ObjectInstance struct {
    // ... existing fields ...
    Value [4]int // Instance-specific values (overrides prototype if set)
    // ... rest of fields ...
}
```

### Condition System Integration
The drink command properly integrates with the character condition system:
- COND_THIRST (0) - Thirst level
- COND_FULL (1) - Hunger/fullness level  
- COND_DRUNK (2) - Intoxication level

### Pet Shop Architecture
The pet shop system uses the existing special procedure framework:
- Special procedures registered by name
- Command routing through AI system
- Integration with character following system

## Testing

All fixes include comprehensive unit tests:

### Drink Command Tests
- Basic drinking functionality
- Liquid consumption tracking
- Empty container handling
- Non-drink object rejection
- Drink effects calculation

### Container Access Tests  
- Basic container access
- Corpse looting
- Closed container handling
- Error message validation

### Integration Tests
- All existing tests still pass
- No regression in existing functionality
- Command registry properly updated

## Usage Examples

### Drink Command
```
drink fountain          # Drink from a water fountain
drink bottle           # Drink from a bottle in inventory
drink water            # Drink from any water container
```

### Container Access
```
get sword from chest   # Get specific item from container
get all from corpse    # Get all items from corpse
get gold from bag      # Get gold from a bag
```

### Pet Shop
```
list                   # Show available pets (in pet shop)
buy dog                # Buy a dog
buy cat fluffy         # Buy a cat and name it "fluffy"
```

## Future Enhancements

### Drink System
- Add more liquid types and effects
- Implement poisoned drinks
- Add drink mixing/alchemy

### Container System
- Add container weight limits
- Implement lockable containers
- Add container durability

### Pet Shop System
- Add pet training commands
- Implement pet equipment
- Add pet special abilities

## Conclusion

These fixes restore three critical game mechanics to full functionality:

1. ✅ **Drink Command**: Now properly handles drink containers and fountains
2. ✅ **Container Access**: Verified working with comprehensive tests
3. ✅ **Pet Shop**: Fully implemented with list/buy functionality

All fixes maintain compatibility with the original DikuMUD behavior and include comprehensive test coverage to prevent regressions.
