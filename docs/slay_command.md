# Slay Command Implementation

## Overview

The `slay` command is an admin-only command that instantly kills a target using the normal combat mechanics. Unlike bypassing the combat system entirely, it works by temporarily boosting the attacker's damage to ensure a lethal hit while using all existing combat functions.

## Usage

```
slay <target>
```

**Examples:**
- `slay cityguard` - Slays the cityguard in the current room
- `slay orc` - Slays an orc (if present)

## Requirements

- **Level**: 20+ (Admin level)
- **Position**: Standing
- **Target**: Must be in the same room
- **Target State**: Target must be alive

## How It Works

### 1. Normal Combat Integration

The slay command uses the existing combat system rather than bypassing it:

1. **Damage Calculation**: Temporarily boosts the attacker's `DamRoll` to `target.HP + 50`
2. **Hit Chance**: Temporarily boosts the attacker's `HitRoll` to `50` to ensure the attack hits
3. **Combat System**: Calls `CombatManager.StartCombat()` to initiate normal combat
4. **Restoration**: Immediately restores original `DamRoll` and `HitRoll` values

### 2. Experience and Death Handling

- **Experience**: Awarded through the normal combat system using `calculateExperience()`
- **Death**: Handled through the existing `HandleCharacterDeath()` function
- **Corpse Creation**: Uses the standard corpse creation system
- **Respawn**: NPCs are scheduled for respawn normally

### 3. Combat Flow

```go
// Store original values
originalDamRoll := character.DamRoll
originalHitRoll := character.HitRoll

// Boost damage to ensure kill
character.DamRoll = targetHP + 50  // Lethal damage
character.HitRoll = 50             // Ensure hit

// Use normal combat system
CombatManager.StartCombat(character, target)

// Restore original values
character.DamRoll = originalDamRoll
character.HitRoll = originalHitRoll
```

## Error Handling

The command provides appropriate error messages for various scenarios:

- **No target**: `"slay whom?"`
- **Target not found**: `"they aren't here"`
- **Self-targeting**: `"you can't slay yourself"`
- **Already dead**: `"<target> is already dead"`
- **Insufficient level**: Handled by command system level check

## Messages

### Success Messages

**To Attacker:**
```
You prepare to slay <target> with divine power!
You slay <target> with divine power!
You gain <X> experience points.
```

**To Target:**
```
<attacker> prepares to slay you with divine power!
<attacker> slays you with divine power!
```

**To Room:**
```
<attacker> prepares to slay <target> with divine power!
<attacker> slays <target> with divine power!
```

## Implementation Details

### Files Modified

1. **pkg/command/slay.go**: Main slay command implementation
2. **pkg/command/registry.go**: Command registration
3. **pkg/command/slay_test.go**: Comprehensive unit tests

### Command Properties

- **Name**: `"slay"`
- **Aliases**: None
- **Min Position**: `POS_STANDING`
- **Min Level**: `20` (Admin level)
- **Logged**: `true` (For admin oversight)

### Integration Points

- **Combat Manager**: Uses `CombatManagerInterface` for combat initiation
- **Death Handler**: Integrates with `HandleCharacterDeath()` function
- **Experience System**: Uses existing experience calculation and awarding
- **Message System**: Uses standard character messaging

## Benefits

### 1. Consistency
- Uses all existing combat mechanics
- Maintains consistency with normal combat flow
- Preserves experience and death handling logic

### 2. Safety
- Admin-level restriction prevents abuse
- Proper error handling for edge cases
- Logging for administrative oversight

### 3. Integration
- Works with existing combat managers
- Compatible with all death-related systems
- Maintains corpse creation and respawn logic

## Testing

Comprehensive unit tests cover:

- **Successful slay**: Verifies damage, death, and experience
- **Error cases**: No target, target not found, self-targeting, already dead
- **Command properties**: Name, aliases, level requirements, logging
- **Integration**: Combat manager interaction, death handling

### Test Coverage

```bash
go test ./pkg/command -v -run TestSlayCommand
```

Expected results:
- ✅ `TestSlayCommand_Execute`
- ✅ `TestSlayCommand_Properties` 
- ✅ `TestSlayCommand_AlreadyDead`

## Usage Examples

### Admin Slaying a Mob

```
> slay cityguard
You prepare to slay a cityguard with divine power!
You slay a cityguard with divine power!
You gain 300 experience points.
```

### Error Cases

```
> slay
slay whom?

> slay nonexistent
They aren't here.

> slay self
You can't slay yourself.
```

### Level Restriction

```
> slay cityguard
You don't have sufficient level to use this command.
```

## Original DikuMUD Compatibility

The slay command maintains compatibility with original DikuMUD principles:

- **Admin Commands**: Follows the pattern of other admin commands
- **Combat Integration**: Uses existing combat mechanics rather than shortcuts
- **Experience System**: Maintains the original experience calculation
- **Death Handling**: Uses the standard death and corpse creation system

## Security Considerations

- **Level Restriction**: Only level 20+ characters can use the command
- **Logging**: All slay commands are logged for administrative review
- **No Player Targeting**: Could be extended to prevent slaying other players
- **Room Restriction**: Only targets in the same room can be slayed

## Future Enhancements

Potential improvements could include:

1. **Player Protection**: Prevent slaying other players entirely
2. **Confirmation**: Require confirmation for high-level targets
3. **Cooldown**: Add a cooldown period between slay uses
4. **Audit Trail**: Enhanced logging with timestamps and reasons
