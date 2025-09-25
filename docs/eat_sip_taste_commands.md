# Eat, Sip, and Taste Commands Implementation

## Overview

I have implemented the `eat`, `sip`, and `taste` commands following the original DikuMUD mechanics exactly. These commands handle food consumption, liquid consumption, and sampling mechanics with proper hunger/thirst/drunk condition management and poison effects.

## Commands Implemented

### 1. **Eat Command**
Allows characters to consume food items from their inventory.

**Usage**: `eat <food_item>`

**Original DikuMUD Behavior**:
- Only works with `ITEM_FOOD` type objects (admins level 22+ can eat anything)
- Prevents eating when character is too full (fullness > 20)
- Increases fullness by the food's value
- Handles poison effects from contaminated food
- Completely consumes the food item

### 2. **Sip Command**
Allows characters to take small sips from drink containers.

**Usage**: `sip <drink_container>`

**Original DikuMUD Behavior**:
- Only works with `ITEM_DRINKCON` type objects
- Prevents sipping when too drunk (drunk > 10)
- Consumes exactly 1 unit of liquid
- Applies drink effects based on liquid type (thirst, fullness, drunkenness)
- Handles poison effects from contaminated liquids

### 3. **Taste Command**
Allows characters to sample food or drink containers.

**Usage**: `taste <food_or_drink>`

**Original DikuMUD Behavior**:
- For drink containers: redirects to `sip` command
- For food: consumes 1 unit of food value, adds 1 fullness
- Handles poison effects with shorter duration (2 rounds)
- Destroys food if no value remains after tasting

## Technical Implementation

### Food Mechanics

**Food Values (from `obj.Prototype.Value` array)**:
- `Value[0]`: Food value (how much fullness it provides)
- `Value[1]`: Unused for food
- `Value[2]`: Unused for food  
- `Value[3]`: Poison flag (1 = poisoned, 0 = safe)

**Fullness System**:
- Characters have a `Conditions[COND_FULL]` value (0-24)
- Eating is prevented when fullness > 20
- Food value is added directly to fullness
- Fullness is capped at 24

**Poison Effects**:
- Poisoned food applies `AFF_POISON` bitvector
- Poison duration = `food_value * 2` rounds for eat
- Poison duration = 2 rounds for taste
- Only affects characters below level 21

### Drink Mechanics

**Drink Container Values**:
- `Value[0]`: Maximum capacity
- `Value[1]`: Current liquid amount
- `Value[2]`: Liquid type (0=water, 1=beer, 2=wine, etc.)
- `Value[3]`: Poison flag (1 = poisoned, 0 = safe)

**Liquid Effects** (from original DikuMUD drink table):
```go
// DRUNK, FULL, THIRST effects per unit
{0, 1, 10}, // water
{3, 2, 5},  // beer  
{5, 2, 5},  // wine
{2, 2, 5},  // ale
// ... etc
```

**Condition Updates**:
- Effects are applied as: `(effect * amount) / 4`
- All conditions are capped at 24
- Drunk condition > 10 prevents sipping

### Poison System

**Poison Affects**:
- Type: `SPELL_POISON` (33)
- Bitvector: `AFF_POISON`
- Location: `APPLY_NONE`
- Duration varies by consumption method

**Poison Messages**:
- Eat: "Ooups, it tasted rather strange ?!!?"
- Sip: "Ooups, it tasted rather strange!"
- Taste: "Ooups, it did not taste good at all!"

## Files Created/Modified

### New Command Files
1. **`pkg/command/eat.go`** - Eat command implementation
2. **`pkg/command/sip.go`** - Sip command implementation  
3. **`pkg/command/taste.go`** - Taste command implementation

### Test Files
4. **`pkg/command/eat_test.go`** - Comprehensive eat command tests
5. **`pkg/command/sip_taste_test.go`** - Sip and taste command tests

### Modified Files
6. **`pkg/command/registry.go`** - Registered new commands

## Test Coverage

### Eat Command Tests
- ✅ Basic food consumption
- ✅ Poisoned food handling
- ✅ Too full prevention
- ✅ Non-food rejection
- ✅ Admin privileges (can eat anything)
- ✅ Item not found handling
- ✅ No arguments handling

### Sip Command Tests  
- ✅ Basic liquid consumption
- ✅ Too drunk prevention
- ✅ Empty container handling
- ✅ Drink effects application
- ✅ Liquid amount reduction

### Taste Command Tests
- ✅ Food tasting (partial consumption)
- ✅ Drink container redirection to sip
- ✅ Complete food consumption when depleted
- ✅ Poisoned food handling
- ✅ Non-food/non-drink rejection

## Original DikuMUD Compatibility

### Message Format
All messages match the original DikuMUD exactly:
- "You eat the bread."
- "TestPlayer eats a loaf of bread."
- "You sip the water."
- "You taste the apple."
- "You are too full to eat more!"
- "You simply fail to reach your mouth!" (too drunk)

### Mechanics Fidelity
- Condition calculations use original formulas
- Poison durations match original values
- Admin level privileges (22+ for eat, 21+ for poison immunity)
- Exact same value array usage and interpretation

### Error Handling
- "you can't find it!" (item not found)
- "your stomach refuses to eat that!?!" (non-food)
- "you can't sip from that!" (non-drink container)
- "taste that?!? Your stomach refuses!" (non-food/non-drink)

## Usage Examples

### Eating Food
```
> eat bread
You eat a loaf of bread.
> eat bread
You are too full to eat more!
```

### Sipping Drinks
```
> sip bottle
You sip the water.
> sip bottle
You simply fail to reach your mouth!  (if too drunk)
```

### Tasting Items
```
> taste apple
You taste a red apple.
> taste bottle
You sip the wine.  (redirected to sip)
> taste sword
Taste that?!? Your stomach refuses!
```

### Poison Effects
```
> eat meat
You eat a piece of meat.
Ooups, it tasted rather strange ?!!?
TestPlayer coughs and utters some strange sounds.
```

## Integration

The commands are fully integrated with the existing DikuGo systems:

- **Command Registry**: Registered in `InitRegistry()`
- **Condition System**: Uses existing `Conditions[3]int` array
- **Affect System**: Uses existing `Affect` struct and bitvectors
- **Object System**: Works with `ObjectInstance` and `Object` prototypes
- **Room System**: Sends messages to room occupants
- **Inventory System**: Properly removes consumed items

## Benefits

### ✅ **Complete DikuMUD Compatibility**
- Exact behavior match with original C implementation
- Same message format and error handling
- Identical mechanics and calculations

### ✅ **Robust Implementation**
- Comprehensive error handling
- Proper boundary checking (fullness, drunkenness)
- Safe object manipulation

### ✅ **Thorough Testing**
- 15+ unit tests covering all scenarios
- Edge case handling verified
- Poison mechanics tested

### ✅ **Clean Code Structure**
- Follows existing DikuGo patterns
- Proper separation of concerns
- Reusable helper functions

The eat, sip, and taste commands are now fully functional and provide an authentic DikuMUD experience with proper food/drink mechanics, condition management, and poison effects.
