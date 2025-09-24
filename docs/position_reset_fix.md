# Player Position Reset Fix

## Problem

Players were getting stuck in invalid positions (fighting, dead, sleeping, etc.) after logging out and back in. This caused commands to not work properly because the game thought the player was still in a fighting state or other invalid position.

## Root Cause

The DikuGo server was saving and loading the `Position` field from player files, which meant players could log back in with positions like:
- `POS_FIGHTING` (7) - Stuck in combat state
- `POS_DEAD` (0) - Dead state
- `POS_SLEEPING` (4) - Sleeping state
- Other invalid positions

This differed from the original DikuMUD behavior, which always resets players to standing position when they enter the game.

## Solution

### 1. Player Position Reset Function

Added `resetPlayerCharacter()` function in `pkg/world/world.go` that:
- Clears the `Fighting` field (removes combat target)
- Sets `Position` to `POS_STANDING` (8)
- Ensures minimum HP, Mana, and Move points (at least 1)
- Logs the reset for debugging

```go
func (w *World) resetPlayerCharacter(ch *types.Character) {
    // Clear fighting state
    ch.Fighting = nil
    
    // Reset position to standing (most important fix)
    ch.Position = types.POS_STANDING
    
    // Ensure minimum hit points, mana, and movement
    if ch.HP <= 0 {
        ch.HP = 1
    }
    // ... etc
}
```

### 2. Automatic Reset on Login

Modified `AddCharacter()` in `pkg/world/world.go` to call `resetPlayerCharacter()` for all non-NPC characters when they enter the game:

```go
// Reset character state when entering the game (like original DikuMUD reset_char)
if !character.IsNPC {
    w.resetPlayerCharacter(character)
    w.InitializeCharacterSkills(character)
}
```

### 3. Storage Layer Changes

Modified `SavePlayer()` in `pkg/storage/player_storage.go` to:
- Always save `Position` as `POS_STANDING` (regardless of actual position)
- Exclude transient fields like `Fighting`, `InRoom`, `World`
- Use explicit field copying to avoid mutex copying issues

## Original DikuMUD Behavior

This fix matches the original DikuMUD behavior from `old/db.c`:

```c
void reset_char(struct char_data *ch) {
    // ... other resets ...
    ch->specials.position = POSITION_STANDING;
    ch->specials.default_pos = POSITION_STANDING;
    // ... ensure minimum HP/mana/move ...
}
```

The `reset_char()` function is called every time a player enters the game in the original DikuMUD.

## Benefits

1. **No Stuck Fighting States**: Players can't get stuck unable to use commands after combat
2. **Consistent Login Experience**: All players always start in standing position
3. **Original Compatibility**: Matches original DikuMUD behavior exactly
4. **Robust Error Recovery**: Handles edge cases like negative HP/mana/move
5. **Clean State Management**: Separates persistent data from transient game state

## Testing

Comprehensive tests were added in `pkg/world/player_position_test.go`:

- `TestPlayerPositionResetOnLogin`: Verifies players in various positions are reset to standing
- `TestNPCPositionNotReset`: Ensures NPCs keep their original positions
- `TestResetPlayerCharacterFunction`: Tests the reset function directly
- `TestPlayerStorageDoesNotSavePosition`: Documents storage behavior

## Files Modified

1. **pkg/world/world.go**: Added `resetPlayerCharacter()` and integrated it into `AddCharacter()`
2. **pkg/storage/player_storage.go**: Modified `SavePlayer()` to exclude transient fields
3. **pkg/world/player_position_test.go**: Added comprehensive tests

## Usage

The fix is automatic - no configuration needed. Players will always enter the game in standing position regardless of how they logged out.

## Debugging

The reset function logs its actions:
```
Reset player character PlayerName: Position=8, HP=50, Mana=100, Move=100
```

Position value 8 corresponds to `POS_STANDING`, confirming the reset worked correctly.
